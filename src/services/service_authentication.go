package services

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/arangodb/go-driver/v2/arangodb"
	"github.com/golang-jwt/jwt/v5"

	"gourmand.golang.arango/src/entities"
	"gourmand.golang.arango/src/enums"
)

type Token struct {
	Email    string     `json:"email"`
	Role     enums.Role `json:"role"`
	jwt.RegisteredClaims
}

const authenticationCollection = "authentications"

type AuthenticationService struct {
	db arangodb.Database
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func NewAuthenticationService(db arangodb.Database) (*AuthenticationService, error) {
	privatePath := os.Getenv("SSL_PRIVATE_KEY_PATH")
	if privatePath == "" {
		return nil, errors.New("SSL_PRIVATE_KEY_PATH environment variable is not set")
	}

	keyBytes, err := os.ReadFile(privatePath)
	if err != nil {
		return nil, fmt.Errorf("reading private key file: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("parsing private key: %w", err)
	}

	publicPath := os.Getenv("SSL_PUBLIC_KEY_PATH")
	if publicPath == "" {
		return nil, errors.New("SSL_PUBLIC_KEY_PATH environment variable is not set")
	}

	keyBytes, err = os.ReadFile(publicPath)
	if err != nil {
		return nil, fmt.Errorf("reading public key file: %w", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("parsing public key: %w", err)
	}

	return &AuthenticationService{db: db, privateKey: privateKey, publicKey: publicKey}, nil
}

func (s *AuthenticationService) CreateAuthentication(
	ctx context.Context, email, password string) (*entities.Authentication, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	auth := entities.Authentication{
		Email:    email,
		Password: string(hashedPassword),
	}

	col, err := s.db.GetCollection(ctx, authenticationCollection, nil)
	if err != nil {
		return nil, err
	}

	meta, err := col.CreateDocument(ctx, auth)
	if err != nil {
		return nil, err
	}

	auth.Id = meta.Key

	return &auth, nil
}

func (s *AuthenticationService) GenerateToken(ctx context.Context, email string) (string, error) {

	claims := Token{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(s.privateKey)
}

func (s *AuthenticationService) AuthenticateUser(
	ctx context.Context, email, password string) (bool, error) {

	query := `FOR u IN @@collection 
		FILTER u.email == @email 
		LIMIT 1 
		RETURN u`

	bindVars := map[string]interface{}{
		"@collection": authenticationCollection,
		"email":       email,
	}

	options := &arangodb.QueryOptions{
		BindVars: bindVars,
	}

	cursor, err := s.db.Query(ctx, query, options)
	if err != nil {
		return false, err
	}

	defer cursor.Close()

	var auth entities.Authentication

	_, err = cursor.ReadDocument(ctx, &auth)
	if err != nil {
		return false, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(password)); err != nil {
		return false, errors.New("invalid password")
	}

	return true, nil
}

func (s *AuthenticationService) AuthenticateToken(ctx context.Context, tokenString string) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Token{}, 
		func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.publicKey, nil
	})

	if err != nil {
		return false, fmt.Errorf("parsing token: %w", err)
	}	

	if _, ok := token.Claims.(*Token); ok && token.Valid {
		return true, nil
	}
	return false, errors.New("invalid token")
}

