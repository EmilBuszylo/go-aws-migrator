package ddb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"golang-migrator-example/src/examples/ddb/scripts"
	"golang-migrator-example/src/migrator"
)

// defsExample contains definitions of our migrations
// conf - needs for use some AWS services inside our function for example: SQS, SNS, lambda ..
// db - necessary to use DDB table
func defsExample(ctx context.Context, conf aws.Config, db *dynamodb.Client) []migrator.Definition {
	return []migrator.Definition{
		{
			Name: "#1 example migration",
			Func: func() error {
				err := scripts.CreateRecord(ctx, db, "example-records-table")
				return err
			},
		},
	}
}
