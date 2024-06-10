package repository

import (
	"context"
	"payment-service/internal/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoPaymentRepository struct {
	client *mongo.Client
}

func NewMongoPaymentRepository(client *mongo.Client) domain.PaymentRepository {
	return &MongoPaymentRepository{
		client: client,
	}
}

func (r *MongoPaymentRepository) Save(ctx context.Context, payment *domain.Payment) error {
	collection := r.client.Database("paymentdb").Collection("payments")
	payment.CreatedAt = time.Now()
	payment.UpdatedAt = time.Now()
	_, err := collection.InsertOne(ctx, payment)
	return err
}

func (r *MongoPaymentRepository) FindByID(ctx context.Context, paymentID string) (*domain.Payment, error) {
	collection := r.client.Database("paymentdb").Collection("payments")
	var payment domain.Payment
	err := collection.FindOne(ctx, bson.M{"paymentid": paymentID}).Decode(&payment)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *MongoPaymentRepository) FindByUserID(ctx context.Context, userID string, page, pageSize int) ([]domain.Payment, int, error) {
	collection := r.client.Database("paymentdb").Collection("payments")
	var payments []domain.Payment
	filter := bson.M{"userid": userID}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var payment domain.Payment
		if err = cursor.Decode(&payment); err != nil {
			return nil, 0, err
		}
		payments = append(payments, payment)
	}

	if err = cursor.Err(); err != nil {
		return nil, 0, err
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return payments, int(count), nil
}

func (r *MongoPaymentRepository) UpdateStatus(ctx context.Context, paymentID, status string) error {
	collection := r.client.Database("paymentdb").Collection("payments")
	_, err := collection.UpdateOne(ctx, bson.M{"paymentid": paymentID}, bson.M{"$set": bson.M{"status": status, "updatedat": time.Now()}})
	return err
}
