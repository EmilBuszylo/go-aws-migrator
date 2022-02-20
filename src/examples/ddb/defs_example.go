package ddb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"golang-migrator-example/src/examples/ddb/scripts"
	"golang-migrator-example/src/migrator"
)

// DefsExample contains definitions of our migrations
// db - necessary to use DDB table
func DefsExample(ctx context.Context, db *dynamodb.Client) []migrator.Definition {
	return []migrator.Definition{
		{
			Name: "#3 example migration",
			Func: func() error {
				// Insert your table name in place of mine
				err := scripts.CreateRecord(ctx, db, "example-records-table")
				return err
			},
		},
		{
			Name: "#2 example migration",
			Func: func() error {
				// Insert your table name in place of mine
				err := scripts.CreateRecord(ctx, db, "example-records-table")
				return err
			},
		},
		// The First migration definition
		{
			Name: "#1 example migration",
			Func: func() error {
				// Insert your table name in place of mine
				err := scripts.CreateRecord(ctx, db, "example-records-table")
				return err
			},
		},
	}
}
