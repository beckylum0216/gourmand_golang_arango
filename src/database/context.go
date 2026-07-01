package database

import (
	"context"
	"fmt"

	"github.com/arangodb/go-driver/v2/arangodb"
	"github.com/arangodb/go-driver/v2/connection"
)

func NewArangoDB(ctx context.Context, dsn, username, password, dbName string) (*arangodb.Database, error) {
	conn := connection.NewHttp2Connection(
		connection.DefaultHTTP2ConfigurationWrapper(
			connection.NewRoundRobinEndpoints([]string{dsn}), true,
		),
	)

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
	} else {
		db, err = client.CreateDatabase(ctx, dbName, nil)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get/create database: %w", err)
	}

	return &db, nil
}