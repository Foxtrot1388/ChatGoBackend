package storage

import (
	"ChatGo/internal/domain/entity"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	db *mongo.Database
}

func New(db *mongo.Database) *Storage {
	return &Storage{db: db}
}

func (bs *Storage) CreateUser(user *entity.User) error {

	parUser := bson.M{"_id": user.Login, "pass": user.GetHash()}

	coll := bs.db.Collection("Users")
	_, err := coll.InsertOne(context.TODO(), parUser)
	if err != nil {
		return err
	}

	return nil
}

func (bs *Storage) DeleteUser(user *entity.User) error {

	parUser := bson.M{"_id": user.Login}

	coll := bs.db.Collection("Users")
	_, err := coll.DeleteOne(context.TODO(), parUser)
	if err != nil {
		return err
	}

	return nil
}

func (bs *Storage) LoginUser(user *entity.User) (*entity.FindUser, error) {

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

func (bs *Storage) FindUser(user string) (*entity.ListUser, error) {

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

func (bs *Storage) FindOneUser(user string) (*entity.FindUser, error) {

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
