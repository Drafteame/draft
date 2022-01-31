package tpl

const (
	MongoRepository = `package {{PackageName}}

import (
	"context"

	"github.com/Drafteame/engine/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const collection = "{{CollectionName}}"

type {{CamelCaseName}}Repository struct {
	db mongo.Driver
}

func NewRepository(db mongo.Driver) *{{CamelCaseName}}Repository {
	return &{{CamelCaseName}}Repository{db}
}

func (r *{{CamelCaseName}}Repository) Find(filters mongo.Filters, opts *mongo.FindOptions, res interface{}) error {
	return r.db.Find(context.Background(), collection, &filters, opts, res)
}

func (r *{{CamelCaseName}}Repository) FindOne(filters mongo.Filters, res interface{}) error {
	return r.db.FindOne(context.Background(), collection, &filters, res)
}

func (r *{{CamelCaseName}}Repository) Insert(data []interface{}) error {
	return r.db.Insert(context.Background(), collection, data)
}

func (r *{{CamelCaseName}}Repository) InsertOne(data interface{}) (*primitive.ObjectID, error) {
	return r.db.InsertOne(context.Background(), collection, data)
}

func (r *{{CamelCaseName}}Repository) Update(filters mongo.Filters, data interface{}) error {
	return r.db.Update(context.Background(), collection, &filters, data)
}

func (r *{{CamelCaseName}}Repository) UpdateOne(filters mongo.Filters, data interface{}) error {
	return r.db.UpdateOne(context.Background(), collection, &filters, data)
}

func (r *{{CamelCaseName}}Repository) Delete(filters mongo.Filters) error {
	return r.db.Delete(context.Background(), collection, &filters)
}

func (r *{{CamelCaseName}}Repository) DeleteOne(filters mongo.Filters) error {
	return r.db.DeleteOne(context.Background(), collection, &filters)
}`
)
