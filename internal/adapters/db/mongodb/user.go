package storage

import (
	"ChatGo/internal/domain/entity"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	db *mongo.Database
}

func New(db *mongo.Database) *Storage {
	return &Storage{db: db}
}

func (bs *Storage) Create(user *entity.User) error {

	parUser := bson.M{"_id": user.Login}

	coll := bs.db.Collection("Users")
	_, err := coll.InsertOne(context.TODO(), parUser)
	if err != nil {
		return err
	}

	return nil
}

func (bs *Storage) Delete(user *entity.User) error {

	parUser := bson.M{"_id": user.Login}

	coll := bs.db.Collection("Users")
	_, err := coll.DeleteOne(context.TODO(), parUser)
	if err != nil {
		return err
	}

	return nil
}

func (bs *Storage) Login(user *entity.User) (*entity.FindUser, error) {

	parUser := bson.M{"_id": user.Login, "pass": user.GetHash()}
	opts := options.FindOne().SetProjection(bson.D{{"_id", 1}})

	coll := bs.db.Collection("Users")
	res := coll.FindOne(context.TODO(), parUser, opts)
	var result entity.FindUser
	err := res.Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (bs *Storage) Find(user string) (*entity.ListUser, error) {

	parUser := bson.M{"_id": bson.M{"$regex": user, "$options": "im"}}
	opts := options.Find().SetProjection(bson.D{{"_id", 1}})

	coll := bs.db.Collection("Users")
	cursor, err := coll.Find(context.TODO(), parUser, opts)
	if err != nil {
		return nil, err
	}

	var results entity.ListUser
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	if err = cursor.Err(); err != nil {
		return nil, err
	}
	return &results, nil
}

func (bs *Storage) FindOne(user string) (*entity.FindUser, error) {

	parUser := bson.M{"_id": user}
	opts := options.FindOne().SetProjection(bson.D{{"_id", 1}})

	coll := bs.db.Collection("Users")
	res := coll.FindOne(context.TODO(), parUser, opts)
	var result entity.FindUser
	err := res.Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (bs *Storage) AddContact(curuser *entity.FindUser, adduser *entity.FindUser) (string, error) {

	parUser := bson.M{"user": curuser.Login, "contact": adduser.Login}

	coll := bs.db.Collection("ContactList")
	result, err := coll.InsertOne(context.TODO(), parUser)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil

}

func (bs *Storage) DeleteContact(id string) error {

	bid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	parUser := bson.M{"_id": bid}

	coll := bs.db.Collection("ContactList")
	_, err = coll.DeleteOne(context.TODO(), parUser)
	if err != nil {
		return err
	}

	return nil

}
