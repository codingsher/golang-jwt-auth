package sqlite

import (
	"database/sql"

	"github.com/codingsher/user-jwt-auth/internal/config"
	"github.com/codingsher/user-jwt-auth/internal/types"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
		);`)
	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil
}

func (s *Sqlite) NewUserRegister(userName, userEmail, userPasswordHash string) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO users (username, email, password) VALUES(?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(userName, userEmail, userPasswordHash)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil
}

func (s *Sqlite) GetUserByEmail(userEmail string) (string, error) {
	stmt, err := s.Db.Prepare("SELECT password FROM users WHERE email=? LIMIT 1")
	if err != nil {
		return "err", err
	}
	defer stmt.Close()

	var user types.UserLogin

	err = stmt.QueryRow(userEmail).Scan(&user.Password)
	if err != nil {
		return "err", err
	}

	return user.Password, nil
}
