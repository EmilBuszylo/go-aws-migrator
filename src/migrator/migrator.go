package migrator

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// ErrMigrationHole means that number of migrations definitions and current version differ.
// Number of migrations definitions should always be greater or equal to current version.
// It's equal when no migrations should be fired.
var ErrMigrationHole = errors.New("too few migrations, did you remove any by accident?")

// API needed to fulfill the contract, we could use *dynamodb.Client instead, but if you want to generate
// mocks for tests create own interface will be better choice.
type API interface {
	PutItem(context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	Query(context.Context, *dynamodb.QueryInput, ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
}

// DefaultDDBProvider proposition of provider for your migration definitions.
// this provider is a perfect choice if you want to make migrate on DDB tables.
type DefaultDDBProvider interface {
	Provide(ctx context.Context, migrationSet string, sess aws.Config, db *dynamodb.Client) ([]Definition, error)
}

// DefaultDDBProviderOptions provide defaults options of migration provider for DDB migrations
type DefaultDDBProviderOptions struct {
	Ctx          context.Context
	MigrationSet string
	Conf         aws.Config
	DB           *dynamodb.Client
}

// Migrator allows firing migration definitions.
type Migrator struct {
	db        API
	tableName string
}

// Definition keeps shape of single definition
type Definition struct {
	Name string
	Func func() error
}

// Summary keeps details about launched migration
type Summary struct {
	StartingVersion int
	CurrentVersion  int
	Executions      []Execution
}

// Execution keeps single migration execution details.
type Execution struct {
	Name    string
	FiredAt time.Time
	Elapsed time.Duration
}

// version keeps data about single done migration
type version struct {
	MigrationSet string    `dynamodbav:"migration_set"`
	VersionNo    int       `dynamodbav:"version_number"`
	Name         string    `dynamodbav:"name"`
	FiredAt      time.Time `dynamodbav:"firedAt"`
	Elapsed      int64     `dynamodbav:"elapsed"`
}
