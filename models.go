package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Animator struct {
	ID          uint64      `db:"id" form:"animator_id"`
	UUID        string      `db:"uuid" form:"uuid"`
	Username    string      `db:"username" form:"username"`
	Password    string      `db:"password" form:"password"`
	DisplayName string      `db:"display_name" form:"display_name"`
	Active      bool        `db:"active"`
	CreatedAt   time.Time   `db:"created_at"`
	LastLoginAt pq.NullTime `db:"last_login_at"`
}

type Animation struct {
	ID          uint64      `db:"id"`
	UUID        string      `db:"uuid" form:"uuid"`
	AnimatorID  uint64      `db:"animator_id" form:"animator_id"`
	Title       string      `db:"title" form:"title"`
	Details     string      `db:"details" form:"details"`
	Published   bool        `db:"published" form:"published"`
	Visible     bool        `db:"visible" form:"visible"`
	Animation   string      `db:"animation" form:"animation"`
	CreatedAt   time.Time   `db:"created_at"`
	UpdatedAt   pq.NullTime `db:"updated_at"`
	PublishedAt pq.NullTime `db:"published_at"`
}

type PublishedAnimation struct {
	UUID        string `db:"uuid"`
	Title       string `db:"title"`
	Details     string `db:"details"`
	DisplayName string `db:"display_name"`
}

var DB *sqlx.DB

func InitDB() {
	db, err := sqlx.Connect("postgres", "user=anigram dbname=anigram sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(2)
	db.SetMaxIdleConns(2)

	DB = db

	fmt.Println("Connected to anigram db")
}

// func GetAnimations() ([]PublishedAnimation, error) {
// 	animations := []PublishedAnimation{}

// 	err := config.DB.Select(&animations, "select a.uuid as uuid, a.title as title, a.details as details, ar.display_name as display_name from animation a join animator ar on a.animator_id = ar.id where a.visible = true and a.published = true order by a.published_at desc")

// 	return animations, err
// }
