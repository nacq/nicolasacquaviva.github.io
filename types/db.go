package types

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Datastore interface {
	AddContent(Content) (*mongo.InsertOneResult, error)
	SaveCommand(string, bool, string, string) (*mongo.InsertOneResult, error)
	GetContentByName(string) (*Content, error)
	GetContentByParentDir(string) ([]string, error)
	GetContentByPath(string) (*Content, error)
	GetContentsParentByChild(string) (*Content, error)
	GetFileContent(string) string
}
