package tpl

const (
	MongoRepositoryInterface = `package {{PackageName}}

import (
	"github.com/Drafteame/framework/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	Find(*mongo.Filters, *mongo.FindOptions, interface{}) error
	FindOne(*mongo.Filters, interface{}) error
	Insert([]interface{}) error
	InsertOne(interface{}) (*primitive.ObjectID, error)
	Update(*mongo.Filters, interface{}) error
	UpdateOne(*mongo.Filters, interface{}) error
	Delete(*mongo.Filters) error
	DeleteOne(*mongo.Filters) error
}
`

	MongoRepository = `package {{PackageName}}

import (
	"context"

	"github.com/Drafteame/framework/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const collection = "{{CollectionName}}"

type repository struct {
	db mongo.Driver
}

func NewRepository(db mongo.Driver) Repository {
	return &repository{db}
}

func (r *repository) Find(filters *mongo.Filters, opts *mongo.FindOptions, res interface{}) error {
	return r.db.Find(context.Background(), collection, filters, opts, res)
}

func (r *repository) FindOne(filters *mongo.Filters, res interface{}) error {
	return r.db.FindOne(context.Background(), collection, filters, res)
}

func (r *repository) Insert(data []interface{}) error {
	return r.db.Insert(context.Background(), collection, data)
}

func (r *repository) InsertOne(data interface{}) (*primitive.ObjectID, error) {
	return r.db.InsertOne(context.Background(), collection, data)
}

func (r *repository) Update(filters *mongo.Filters, data interface{}) error {
	return r.db.Update(context.Background(), collection, filters, data)
}

func (r *repository) UpdateOne(filters *mongo.Filters, data interface{}) error {
	return r.db.UpdateOne(context.Background(), collection, filters, data)
}

func (r *repository) Delete(filters *mongo.Filters) error {
	return r.db.Delete(context.Background(), collection, filters)
}

func (r *repository) DeleteOne(filters *mongo.Filters) error {
	return r.db.DeleteOne(context.Background(), collection, filters)
}
`
)
