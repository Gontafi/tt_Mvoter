package repository

import (
	"context"
	"errors"
	"time"
	"tt/internal/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type DynamicDataRepositoryInterface interface {
	CreateTable(ctx context.Context, name string) (int64, error)
	CreateRow(ctx context.Context, tableID int64, data map[string]interface{}) (int64, error)
	GetRows(ctx context.Context, tableID int64) ([]models.Row, error)
	UpdateRow(ctx context.Context, tableID int64, rowID int64, data map[string]interface{}) error
	DeleteRow(ctx context.Context, tableID int64, rowID int64) error
	GetTable(ctx context.Context, tableID int64) (*models.Table, error)
	DeleteTable(ctx context.Context, tableID int64) error
	GetTables(ctx context.Context) ([]models.Table, error)
}

type DynamicDataRepository struct {
	tableCollection   *mongo.Collection
	rowCollection     *mongo.Collection
	counterCollection *mongo.Collection
}

func NewDynamicDataRepository(db *mongo.Database) *DynamicDataRepository {
	tableCollection := db.Collection("tables")
	rowCollection := db.Collection("rows")
	counterCollection := db.Collection("counters")

	indexModel := mongo.IndexModel{
		Keys: bson.D{{"table_id", 1}},
	}
	_, err := rowCollection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		panic(err)
	}

	return &DynamicDataRepository{
		tableCollection:   tableCollection,
		rowCollection:     rowCollection,
		counterCollection: counterCollection,
	}
}

func (r *DynamicDataRepository) getNextSequence(ctx context.Context, sequenceName string) (int64, error) {
	filter := bson.M{"_id": sequenceName}
	update := bson.M{"$inc": bson.M{"value": 1}}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var counter models.Counter
	err := r.counterCollection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&counter)
	if err != nil {
		return 0, err
	}

	return counter.Value, nil
}

func (r *DynamicDataRepository) CreateTable(ctx context.Context, name string) (int64, error) {
	id, err := r.getNextSequence(ctx, "table_id")
	if err != nil {
		return 0, err
	}

	table := models.Table{
		ID:        id,
		Name:      name,
		CreatedAt: time.Now(),
	}

	_, err = r.tableCollection.InsertOne(ctx, table)
	if err != nil {
		return 0, err
	}

	return table.ID, nil
}

func (r *DynamicDataRepository) CreateRow(ctx context.Context, tableID int64, data map[string]interface{}) (int64, error) {
	var table models.Table
	err := r.tableCollection.FindOne(ctx, bson.M{"_id": tableID}).Decode(&table)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 0, errors.New("table not found")
		}
		return 0, err
	}

	rowID, err := r.getNextSequence(ctx, "row_id")
	if err != nil {
		return 0, err
	}

	now := time.Now()
	row := models.Row{
		ID:        rowID,
		TableID:   tableID,
		Data:      data,
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err = r.rowCollection.InsertOne(ctx, row)
	if err != nil {
		return 0, err
	}

	return row.ID, nil
}

func (r *DynamicDataRepository) GetRows(ctx context.Context, tableID int64) ([]models.Row, error) {
	filter := bson.M{"table_id": tableID}
	opts := options.Find().SetSort(bson.D{{"created_at", 1}})

	cursor, err := r.rowCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var rows []models.Row
	if err := cursor.All(ctx, &rows); err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *DynamicDataRepository) UpdateRow(ctx context.Context, tableID, rowID int64, data map[string]interface{}) error {
	filter := bson.M{"_id": rowID, "table_id": tableID}
	update := bson.M{
		"$set": bson.M{
			"data":       data,
			"updated_at": time.Now(),
		},
	}

	result, err := r.rowCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("row not found")
	}

	return nil
}

func (r *DynamicDataRepository) DeleteRow(ctx context.Context, tableID, rowID int64) error {
	filter := bson.M{"_id": rowID, "table_id": tableID}

	result, err := r.rowCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("row not found")
	}

	return nil
}

func (r *DynamicDataRepository) GetTable(ctx context.Context, tableID int64) (*models.Table, error) {
	var table models.Table
	err := r.tableCollection.FindOne(ctx, bson.M{"_id": tableID}).Decode(&table)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("table not found")
		}
		return nil, err
	}

	return &table, nil
}

func (r *DynamicDataRepository) DeleteTable(ctx context.Context, tableID int64) error {
	session, err := r.tableCollection.Database().Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(ctx2 context.Context) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}

		tableFilter := bson.M{"_id": tableID}
		tableResult, err := r.tableCollection.DeleteOne(ctx2, tableFilter)
		if err != nil {
			return err
		}
		if tableResult.DeletedCount == 0 {
			return errors.New("table not found")
		}

		rowsFilter := bson.M{"table_id": tableID}
		_, err = r.rowCollection.DeleteMany(ctx2, rowsFilter)
		if err != nil {
			return err
		}

		return session.CommitTransaction(ctx2)
	})

	return err
}

func (r *DynamicDataRepository) GetTables(ctx context.Context) ([]models.Table, error) {
	cursor, err := r.tableCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tables []models.Table
	if err := cursor.All(ctx, &tables); err != nil {
		return nil, err
	}

	return tables, nil
}
