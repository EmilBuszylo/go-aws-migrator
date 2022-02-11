package migrator_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang-migrator-example/src/migrator"
	"golang-migrator-example/src/migrator/mock"
)

func TestMigrator_Run(t *testing.T) {
	ctx := context.Background()
	t.Run("error cases", func(t *testing.T) {
		t.Run("query recent version error", func(t *testing.T) {
			ctrl := gomock.NewController(t)

			db := mock.NewMockAPI(ctrl)
			db.EXPECT().Query(ctx, gomock.Any()).Return(nil, errors.New("random error"))

			summary, err := migrator.New(db, "foo").Run(ctx, "test", nil)
			assert.Error(t, err)
			assert.Equal(t, "random error", err.Error())
			assert.Empty(t, summary)
		})

		t.Run("error from specific migration func", func(t *testing.T) {
			ctrl := gomock.NewController(t)

			db := mock.NewMockAPI(ctrl)
			db.EXPECT().Query(ctx, gomock.Any()).Return(&dynamodb.QueryOutput{
				Items: []map[string]types.AttributeValue{
					{"version_number": &types.AttributeValueMemberN{Value: "0"}}},
			}, nil)
			summary, err := migrator.New(db, "foo").Run(ctx, "test", []migrator.Definition{
				{Name: "testfunc", Func: func() error { return errors.New("boom") }}})
			assert.Error(t, err)
			assert.Equal(t, "migration 'testfunc' failure: boom", err.Error())
			assert.Empty(t, summary)
		})

		t.Run("error put migration in DDB", func(t *testing.T) {
			ctrl := gomock.NewController(t)

			db := mock.NewMockAPI(ctrl)
			db.EXPECT().Query(ctx, gomock.Any()).Return(&dynamodb.QueryOutput{
				Items: []map[string]types.AttributeValue{
					{"version_number": &types.AttributeValueMemberN{Value: "0"}}},
			}, nil)
			db.EXPECT().PutItem(ctx, gomock.Any()).Return(nil, errors.New("boom"))

			summary, err := migrator.New(db, "foo").Run(ctx, "test", []migrator.Definition{
				{Name: "testfunc", Func: func() error { return nil }}})
			assert.Error(t, err)
			if assert.Len(t, summary.Executions, 1) {
				assert.Equal(t, "testfunc", summary.Executions[0].Name)
			}
			assert.Equal(t, 0, summary.CurrentVersion)
		})

		t.Run("error, migration hole", func(t *testing.T) {
			ctrl := gomock.NewController(t)

			db := mock.NewMockAPI(ctrl)
			db.EXPECT().Query(ctx, gomock.Any()).Return(&dynamodb.QueryOutput{
				Items: []map[string]types.AttributeValue{
					{"version_number": &types.AttributeValueMemberN{Value: "2"}},
				},
			}, nil)

			summary, err := migrator.New(db, "foo").Run(ctx, "test", []migrator.Definition{
				{Name: "testfunc", Func: func() error { return nil }}})

			assert.Error(t, err)
			assert.Equal(t, migrator.ErrMigrationHole, err)
			assert.Empty(t, summary)
		})
	})

	t.Run("success cases", func(t *testing.T) {

		t.Run("no migrationas", func(t *testing.T) {
			ctrl := gomock.NewController(t)

			db := mock.NewMockAPI(ctrl)
			db.EXPECT().Query(ctx, gomock.Any()).Return(&dynamodb.QueryOutput{
				Items: []map[string]types.AttributeValue{
					{"version_number": &types.AttributeValueMemberN{Value: "0"}},
				},
			}, nil)

			summary, err := migrator.New(db, "foo").Run(ctx, "test", nil)
			assert.Empty(t, err)
			assert.Empty(t, summary)
		})
		t.Run("one migration", func(t *testing.T) {
			ctrl := gomock.NewController(t)

			db := mock.NewMockAPI(ctrl)
			db.EXPECT().Query(ctx, gomock.Any()).Return(&dynamodb.QueryOutput{
				Items: []map[string]types.AttributeValue{
					{"version_number": &types.AttributeValueMemberN{Value: "0"}},
				},
			}, nil)
			db.EXPECT().PutItem(ctx, gomock.Any()).Return(nil, nil)

			summary, err := migrator.New(db, "foo").Run(ctx, "test", []migrator.Definition{
				{Name: "testfunc", Func: func() error { return nil }}})
			assert.Empty(t, err)
			if assert.Len(t, summary.Executions, 1) {
				assert.Equal(t, "testfunc", summary.Executions[0].Name)
			}
			assert.Equal(t, 0, summary.StartingVersion)
			assert.Equal(t, 1, summary.CurrentVersion)
		})
		t.Run("two migrations", func(t *testing.T) {
			ctrl := gomock.NewController(t)

			db := mock.NewMockAPI(ctrl)
			db.EXPECT().Query(ctx, gomock.Any()).Return(&dynamodb.QueryOutput{
				Items: []map[string]types.AttributeValue{
					{"version_number": &types.AttributeValueMemberN{Value: "1"}},
				},
			}, nil)
			db.EXPECT().PutItem(ctx, gomock.Any()).Return(nil, nil)
			db.EXPECT().PutItem(ctx, gomock.Any()).Return(nil, nil)

			summary, err := migrator.New(db, "foo").Run(ctx, "test", []migrator.Definition{
				{Name: "testfunc3", Func: func() error { return nil }},
				{Name: "testfunc2", Func: func() error { return nil }},
				{Name: "testfunc1", Func: func() error { return nil }},
			})
			assert.Empty(t, err)
			if assert.Len(t, summary.Executions, 2) {
				assert.Equal(t, "testfunc2", summary.Executions[0].Name)
				assert.Equal(t, "testfunc3", summary.Executions[1].Name)
			}
			assert.Equal(t, 1, summary.StartingVersion)
			assert.Equal(t, 3, summary.CurrentVersion)
		})
	})
}
