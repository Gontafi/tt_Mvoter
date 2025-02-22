package models

import "time"

type User struct {
	ID        uint64    `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type Table struct {
	ID        int64     `bson:"_id" json:"id"`
	Name      string    `bson:"name" json:"name"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

type Row struct {
	ID        int64                  `json:"id" bson:"_id"`
	TableID   int64                  `json:"table_id" bson:"table_id"`
	Data      map[string]interface{} `json:"data" bson:"data"`
	CreatedAt time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time              `json:"updated_at" bson:"updated_at"`
}

type Counter struct {
	ID    string `bson:"_id"`
	Value int64  `bson:"value"`
}
