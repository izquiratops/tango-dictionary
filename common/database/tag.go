package database

type Tag struct {
	Name        string `bson:"_id"`
	Description string `bson:"description"`
}
