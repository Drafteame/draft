package tpl

var (
	AdapterInterface = `package {{PackageName}}

import (
	"{{Namespace}}/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Adapter interface{
	GetAll() ([]*models.{{ModelName}}, error)
	GetByID(*primitive.ObjectID) (*models.{{ModelName}}, error)
	Create(*models.{{ModelName}}) (*primitive.ObjectID, error)
	Update(*models.{{ModelName}}) error
	Delete(*primitive.ObjectID) error
}

`

	AdapterMongoImplementation = `package {{PackageName}}

import (
	"{{Namespace}}/internal/models"
	"{{Namespace}}/internal/ports/repositories/{{RepoName}}"

	"github.com/Drafteame/framework/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
	mgo "go.mongodb.org/mongo-driver/mongo"
)

type adapter struct {
	{{RepoName}}Repository {{RepoName}}.Repository
}


func (a *adapter) GetAll() ([]*models.{{ModelName}}, error) {
	items := []*models.{{ModelName}}{}

	if err := a.{{RepoName}}Repository.Find(nil, nil, &items); err != nil {
		if err == mgo.ErrNoDocuments {
			return items, nil
		}

		return nil, err
	}

	return items, nil
}

func (a *adapter) GetByID(itemID *primitive.ObjectID) (*models.{{ModelName}}, error) {
	filters := &mongo.Filters{
		{Key: "_id", Value: itemID},
	}
	
	item := models.{{ModelName}}{}

	if err := a.{{RepoName}}Repository.FindOne(filters, &item); err != nil {
		return nil, err
	}

	return &item, nil
}

func (a *adapter) Create(item *models.{{ModelName}}) (*primitive.ObjectID, error) {
	if err := mongo.AddTimestamps(item); err != nil {
		return nil, err
	}

	return a.{{RepoName}}Repository.InsertOne(item)
}

func (a *adapter) Update(item *models.{{ModelName}}) error {
	if err := mongo.AddUpdatedAt(item); err != nil {
		return err
	}

	filters := &mongo.Filters{
		{Key: "_id", Value: item.ID},
	}

	return a.{{RepoName}}Repository.UpdateOne(filters, item)
}

func (a *adapter) Delete(itemID *primitive.ObjectID) error {
	return a.{{RepoName}}Repository.DeleteOne(&mongo.Filters{
		{Key: "_id", Value: itemID},
	})
}

var _ Adapter = &adapter{}

func NewAdapter(db mongo.Driver) Adapter {
	return &adapter{
		{{RepoName}}Repository: {{RepoName}}.NewRepository(db),
	}
}

`
)
