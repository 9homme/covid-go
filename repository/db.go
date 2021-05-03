package repository

import (
	"example.com/covid-go/model"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

var DB DataSource

//go:generate mockgen --source db.go --destination mock/mock_db.go
type DataSource interface {
	GetUsers() ([]model.User, error)
	GetUserByUsername(username string) (model.User, error)
}

type MongoDataSource struct {
	session *mgo.Session
}

func InitMongoDb() error {
	session, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		return err
	}
	DB = MongoDataSource{
		session: session,
	}
	return nil
}

func (r MongoDataSource) GetUserByUsername(username string) (model.User, error) {
	var user model.User
	session := r.session.Copy()
	defer session.Close()
	c := session.DB("covid19").C("user")
	err := c.Find(bson.M{"username": username}).One(&user)
	return user, err
}

func (r MongoDataSource) GetUsers() ([]model.User, error) {
	var users []model.User
	session := r.session.Copy()
	defer session.Close()
	c := session.DB("covid19").C("user")
	err := c.Find(bson.D{}).All(&users)
	return users, err
}
