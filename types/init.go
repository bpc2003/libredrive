package types

import (
	"context"
	"database/sql"
	_ "embed"
	"log"
	"os"
	"path"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"libredrive/crypto"
	"libredrive/models"
)

var db *sql.DB
var Queries *models.Queries
var CTX = context.Background()

//go:embed schema.sql
var ddl string

func init() {
	var err error
	db, err = sql.Open("sqlite3", "users.db")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := db.ExecContext(CTX, ddl); err != nil {
		log.Fatal(err)
	}
	Queries = models.New(db)

	if err = godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	if u, _ := Queries.GetUsers(CTX); len(u) == 0 {
		password, salt := crypto.GeneratePassword(os.Getenv("ADMIN_PASSWORD"), 144)
		_, err = Queries.CreateUser(CTX, models.CreateUserParams{Username: "admin", Password: string(password), Salt: salt, Isadmin: true})
		if err != nil || os.MkdirAll(path.Join("users", "1"), 0750) != nil {
			log.Fatal(err)
		}
	}
}
