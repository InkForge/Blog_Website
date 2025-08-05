package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// MongoTransactionManager implements domain.ITransactionManager
type MongoTransactionManager struct {
	client *mongo.Client
}

func NewMongoTransactionManager(client *mongo.Client) *MongoTransactionManager {
	return &MongoTransactionManager{
		client: client,
	}
}

func (m *MongoTransactionManager) WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	session, err := m.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		return nil, fn(sessCtx)
	})

	return err
} 