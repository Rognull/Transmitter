package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type Message struct {
	Session_id string
	Frequency  float64
	Timestamp  time.Time
}

type Database struct {
	db *gorm.DB
}

func (d *Database) ConnectDB(dsn string) error {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	d.db = db
	return nil
}

func (d *Database) Create(msg Message) error {
	result := d.db.Create(msg)
	return result.Error
}
