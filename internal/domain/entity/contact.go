package entity

type Contact struct {
	Id    string `json:"id" bson:"_id"`
	Login string `json:"login" bson:"user"`
}

type ListContact []Contact
