package utils

import (
	"context"

	"github.com/arangodb/go-driver/v2/arangodb"
)

func DeleteEdges(ctx context.Context, db arangodb.Database, query string, bindVars map[string]interface{}) error {
    opts := &arangodb.QueryOptions{
        BindVars: bindVars,
    }

    _, err := db.Query(ctx, query, opts)
    return err
}
