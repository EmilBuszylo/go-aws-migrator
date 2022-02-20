package scripts

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	"time"
)

// ExampleRecord represents structure of single record in DDB
type ExampleRecord struct {
	ID        string    `dynamodbav:"id"`
	CreatedAt time.Time `dynamodbav:"created_at"`
}

// CreateRecord creates a single record of ExampleRecord in DDB
func CreateRecord(ctx context.Context, db *dynamodb.Client, tableName string) error {
	r := ExampleRecord{
		ID:        uuid.NewString(),
		CreatedAt: time.Now(),
	}

	attrs, err := attributevalue.MarshalMap(r)
	if err != nil {
		return err
	}

	_, err = db.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      attrs,
		TableName: aws.String(tableName),
	})

	if err != nil {
		return err
	}

	return nil
}
