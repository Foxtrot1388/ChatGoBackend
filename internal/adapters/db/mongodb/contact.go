package storage

import (
	"ChatGo/internal/domain/entity"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (bs *Storage) ListContact(login string) (*entity.ListContact, error) {

	parCon := bson.M{"user": login}

	coll := bs.db.Collection("ContactList")
	cursor, err := coll.Find(context.TODO(), parCon)
	if err != nil {
		return nil, err
	}

	var results entity.ListContact
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	if err = cursor.Err(); err != nil {
		return nil, err
	}
	return &results, nil
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
