// internal/infrastructure/repository/mongo_payment_repository.go
package repository

import (
	"context"
	"payment-service/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoPaymentRepository struct {
	client *mongo.Client
}

func NewMongoPaymentRepository(client *mongo.Client) domain.PaymentRepository {
	return &mongoPaymentRepository{
		client: client,
	}
}

func (r *mongoPaymentRepository) Save(ctx context.Context, payment *domain.Payment) error {
	collection := r.client.Database("paymentdb").Collection("payments")
	_, err := collection.InsertOne(ctx, payment)
	return err
}

func (r *mongoPaymentRepository) FindByID(ctx context.Context, paymentID string) (*domain.Payment, error) {
	collection := r.client.Database("paymentdb").Collection("payments")
	filter := bson.M{"paymentid": paymentID}
	var payment domain.Payment
	err := collection.FindOne(ctx, filter).Decode(&payment)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *mongoPaymentRepository) FindByUserID(ctx context.Context, userID string, page, pageSize int) ([]domain.Payment, int, error) {
	collection := r.client.Database("paymentdb").Collection("payments")
	filter := bson.M{"userid": userID}
	options := options.Find()
	options.SetSkip(int64((page - 1) * pageSize))
	options.SetLimit(int64(pageSize))

	cursor, err := collection.Find(ctx, filter, options)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var payments []domain.Payment
	for cursor.Next(ctx) {
		var payment domain.Payment
		err := cursor.Decode(&payment)
		if err != nil {
			return nil, 0, err
		}
		payments = append(payments, payment)
	}

	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return payments, int(total), nil
}

func (r *mongoPaymentRepository) UpdateStatus(ctx context.Context, paymentID, status string) error {
	collection := r.client.Database("paymentdb").Collection("payments")
	filter := bson.M{"paymentid": paymentID}
	update := bson.M{"$set": bson.M{"status": status}}
	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}
