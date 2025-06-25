package transactionsRepository

import (
	"OtusGo/internal/interfaces"
	"OtusGo/internal/transaction"
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ interfaces.TransactionRepositoryInterface = (*TransactionsRepository)(nil)

var collectionName = "transactions"

type TransactionsRepository struct {
	collection *mongo.Collection
}

func NewMongoTransactionRepository() (*TransactionsRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Получаем URI из переменных окружения
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://admin:password@mongodb:27017"
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	collection := client.Database("exchange_service_db").Collection(collectionName)

	return &TransactionsRepository{
		collection: collection,
	}, nil
}

func (r *TransactionsRepository) Add(tx *transaction.Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Генерируем новый ID
	count, err := r.collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return fmt.Errorf("failed to count documents: %v", err)
	}
	tx.ID = int(count) + 1

	_, err = r.collection.InsertOne(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to insert transaction: %v", err)
	}

	return nil
}

func (r *TransactionsRepository) Get(id int) (*transaction.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var tx transaction.Transaction
	filter := bson.M{"id": id}

	err := r.collection.FindOne(ctx, filter).Decode(&tx)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("transaction with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to find transaction: %v", err)
	}

	return &tx, nil
}

func (r *TransactionsRepository) GetAll() []*transaction.Transaction {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil
	}
	defer cursor.Close(ctx)

	var transactions []*transaction.Transaction
	for cursor.Next(ctx) {
		var tx transaction.Transaction
		if err := cursor.Decode(&tx); err != nil {
			continue
		}
		transactions = append(transactions, &tx)
	}

	return transactions
}

func (r *TransactionsRepository) Update(tx *transaction.Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": tx.ID}
	update := bson.M{"$set": tx}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update transaction: %v", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("transaction with ID %d not found", tx.ID)
	}

	return nil
}
