package storage

import (
	"ChatGo/internal/domain/entity"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (bs *Storage) CreateMessage(mes *entity.Message) error {

	parUser := bson.M{"body": mes.Body, "Sender": mes.Sender.Login, "Recipient": mes.Recipient.Login, "Date": mes.Date}

	coll := bs.db.Collection("Messages")
	_, err := coll.InsertOne(context.TODO(), parUser)
	if err != nil {
		return err
	}

	return nil
}

func (bs *Storage) ListMessages(sender *entity.FindUser, recipient *entity.FindUser, afterAt interface{}) (*entity.ListMessage, error) {

	opts := options.Find().SetSort(bson.D{{"Date", 1}})
	opts.SetLimit(10)
	opts.SetProjection(bson.D{{"body", 1}, {"Date", 1}})

	var parMess bson.M
	if afterAt != nil {
		parMess = bson.M{"Sender": sender.Login, "Recipient": recipient.Login, "Date": bson.D{{"$lt", afterAt}}}
	} else {
		parMess = bson.M{"Sender": sender.Login, "Recipient": recipient.Login}
	}

	coll := bs.db.Collection("Messages")
	cursor, err := coll.Find(context.TODO(), parMess, opts)
	if err != nil {
		return nil, err
	}

	var results entity.ListMessage
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	if err = cursor.Err(); err != nil {
		return nil, err
	}
	return &results, nil

}
