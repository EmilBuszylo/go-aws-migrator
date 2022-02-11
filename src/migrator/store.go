package migrator

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// New creates an instance of Migrator.
// Table name should point to a valid migration table, example one defined in testdata/serverless.yml
func New(db API, tableName string) *Migrator {
	return &Migrator{db: db, tableName: tableName}
}

// Run launch migration definitions for specific migrationSet
func (m *Migrator) Run(ctx context.Context, migrationSet string, defs []Definition) (Summary, error) {
	v, err := m.get(ctx, migrationSet)
	if err != nil {
		return Summary{}, err
	}

	if len(defs) == 0 {
		return Summary{}, nil
	}

	currentVersion := v.VersionNo
	if len(defs) < currentVersion {
		return Summary{}, ErrMigrationHole
	}

	var executions []Execution
	for i := len(defs) - currentVersion - 1; i >= 0; i-- {
		def := defs[i]
		firedAt := time.Now()
		var elapsed time.Duration
		err := def.Func()
		if err != nil {
			return Summary{
				StartingVersion: v.VersionNo,
				CurrentVersion:  currentVersion,
				Executions:      executions,
			}, fmt.Errorf("migration '%s' failure: %w", def.Name, err)
		}
		elapsed = time.Since(firedAt)

		executions = append(executions, Execution{
			Name:    def.Name,
			FiredAt: firedAt,
			Elapsed: elapsed,
		})

		err = m.put(ctx, version{
			MigrationSet: migrationSet,
			VersionNo:    currentVersion + 1,
			Name:         def.Name,
			FiredAt:      firedAt,
			Elapsed:      elapsed.Nanoseconds(),
		})
		if err != nil {
			return Summary{
				StartingVersion: v.VersionNo,
				CurrentVersion:  currentVersion,
				Executions:      executions,
			}, err
		}

		currentVersion++
	}

	return Summary{
		StartingVersion: v.VersionNo,
		CurrentVersion:  currentVersion,
		Executions:      executions,
	}, nil

}

// put creates new record in migration DDB table
func (m *Migrator) put(ctx context.Context, v version) error {
	item, err := attributevalue.MarshalMap(v)
	if err != nil {
		return err
	}

	expr, err := expression.NewBuilder().WithCondition(
		expression.And(
			expression.AttributeNotExists(expression.Name("migration_set")),
			expression.AttributeNotExists(expression.Name("version_number")))).Build()
	if err != nil {
		return err
	}

	_, err = m.db.PutItem(ctx, &dynamodb.PutItemInput{
		ConditionExpression:      expr.Condition(),
		ExpressionAttributeNames: expr.Names(),
		Item:                     item,
		TableName:                aws.String(m.tableName),
	})
	if err != nil {
		return err
	}

	return nil
}

// list gets our previous migrations by specific migrationSet and limit
func (m *Migrator) list(ctx context.Context, migrationSet string, limit int) ([]version, error) {
	expr, err := expression.NewBuilder().WithKeyCondition(
		expression.KeyEqual(expression.Key("migration_set"), expression.Value(migrationSet))).Build()
	if err != nil {
		return nil, err
	}

	out, err := m.db.Query(ctx, &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		Limit:                     aws.Int32(int32(limit)),
		// descending order
		ScanIndexForward: aws.Bool(false),
		TableName:        aws.String(m.tableName),
	})
	if err != nil {
		return nil, err
	}

	if len(out.Items) == 0 {
		return nil, nil
	}

	var v []version
	err = attributevalue.UnmarshalListOfMaps(out.Items, &v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// get provide latest migration data
func (m *Migrator) get(ctx context.Context, migrationSet string) (version, error) {
	v, err := m.list(ctx, migrationSet, 1)
	if err != nil {
		return version{}, err
	}

	if len(v) != 1 {
		return version{}, nil
	}

	return v[0], nil
}
