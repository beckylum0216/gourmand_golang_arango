package database

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"bytes"
	"io"
	"time"
	"encoding/json"

	"github.com/arangodb/go-driver/v2/arangodb"
	"github.com/arangodb/go-driver/v2/connection"

	"gourmand.golang.arango/src/entities"
)

func NewArangoDB(ctx context.Context, dsn, username, password, dbName string) (*arangodb.Database, error) {
	config := connection.DefaultHTTP2ConfigurationWrapper(
		connection.NewRoundRobinEndpoints([]string{dsn}), 
		true,
	)
	
	conn := connection.NewHttp2Connection(config)
	auth := connection.NewBasicAuth(username, password)

	if err := conn.SetAuthentication(auth); err != nil {
		return nil, fmt.Errorf("failed to set authentication: %w", err)
	}

	client := arangodb.NewClient(conn)

	exists, err := client.DatabaseExists(ctx, dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to check database existence: %w", err)
	}

	var db arangodb.Database
	if exists {
		db, err = client.GetDatabase(ctx, dbName, nil)
		if err != nil {
			return nil, fmt.Errorf("Failed to create database: %w", err)
		}
	} else {
		collection := map[string]bool{
			"authentications": false,
			"persons": false,
			"users": false,
			"authors": false,
			"recipes": false,
			"ingredients": false,
			"persons_users": true,
			"users_authentications": true,
			"persons_authors": true,
			"authors_recipes": true,
			"recipes_ingredients": true,
		}

		db, err = client.CreateDatabase(ctx, dbName, nil)
		if err != nil {
			return nil, fmt.Errorf("Failed to create database: %w", err)
		}

		for key, value := range collection {
			err := CreateCollection(ctx, db, key, value)
			if err != nil {
				return nil, err
			}
		}

		err = PopulateIngredients(ctx, db)
		if err != nil {
			return nil, err
		}
	}

	return &db, nil
}

func CreateCollection(ctx context.Context, 
	db arangodb.Database, collectionName string, edge bool) error {
	exists, err := db.CollectionExists(ctx, collectionName)
	if err != nil {
		return fmt.Errorf("failed to check collection existence: %w", err)
	}

	if !exists {
		properties := &arangodb.CreateCollectionPropertiesV2{}
		if edge {
			t := arangodb.CollectionTypeEdge
			properties.Type = &t
		}

		_, err = db.CreateCollectionV2(ctx, collectionName, properties)
		if err != nil {
			return fmt.Errorf("failed to create collection: %w", err)
		}
	}

	return nil
}

func PopulateIngredients(ctx context.Context, db arangodb.Database) error {
	collection, err := db.GetCollection(ctx, "ingredients", nil)	
	if err != nil {
		return fmt.Errorf("failed to get ingredients collection: %w", err)
	}
	
	ingredients, err := FetchIngredientsFromFDC()
	if err != nil {
		return fmt.Errorf("failed to fetch ingredients: %w", err)
	}

	for _, ingredient := range ingredients {
		_, err = collection.CreateDocument(ctx, ingredient)
		if err != nil {
			return fmt.Errorf("failed to create ingredient document: %w", err)
		}
	}

	return nil
}

func LoadApiKeyFromFile() (string, error) {
	path := os.Getenv("FDC_API_KEY_PATH")
	if path == "" {
		return "", fmt.Errorf("FDC_API_KEY_PATH environment variable is not set")
	}

	key, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read API key from file: %w", err)
	}

	return string(bytes.TrimSpace(key)), nil
}

func FetchIngredientsFromFDC() ([]entities.FoodDetail, error) {
	apiKey, err := LoadApiKeyFromFile()
	if err != nil {
		return nil, fmt.Errorf("failed to load API key: %w", err)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	url := "https://api.nal.usda.gov/fdc/v1/foods/list?api_key=" + apiKey
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ingredients: %w", err)
	}
	
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var ingredients []entities.FoodDetail
	if err := json.Unmarshal(body, &ingredients); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return ingredients, nil
}
