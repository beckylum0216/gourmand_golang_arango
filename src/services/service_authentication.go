package services

import (
	"context"
	"errors"
	
	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/arangodb/go-driver/v2/arangodb"

	"gourmand.golang.arango/src/enums"
	"gourmand.golang.arango/src/entities"
)

type Token struct {
	Email	string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role enums.Role `json:"role"`
	jwt.RegisteredClaims
}

const authenticationCollection = "authentications"

type AuthenticationService struct{
	db arangodb.Database
}

func NewAuthenticationService(db arangodb.Database) *AuthenticationService {
	return &AuthenticationService{db: db}
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


func (s *AuthenticationService) GenerateToken(ctx context.Context, email string, password string) (string, error) {
	claims := Token{
		Email:    email,
		Password: password,
		RegisteredClaims: jwt.RegisteredClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("your_secret_key"))
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

	if err := bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(password)); 
	err != nil {
		return false, errors.New("invalid password")
	}
	
	return true, nil
}

