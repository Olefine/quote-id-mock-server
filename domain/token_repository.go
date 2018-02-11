package domain

import (
	"database/sql"
	"errors"
	"time"
)

func checkDatabaseErr(err error) {
	if err != nil {
		panic(err)
	}
}

// TokenRepository represents storage layer to create/read/update/delete from source
type TokenRepository struct {
	db *sql.DB
}

// Create creates record in underlying store(default is Sqlite3) by inserting
// token string
// created string(in format 2016-01-13)
func (tr *TokenRepository) Create(token string) error {
	stmt, err := tr.db.Prepare("INSERT INTO tokens(token, created) values(?,?);")
	checkDatabaseErr(err)

	todayDate := time.Now().Local().Format("2006-01-02")
	_, err = stmt.Exec(token, todayDate)

	return err
}

// IsTokenValid checks if token exists in database and fresh
func (tr *TokenRepository) IsTokenValid(token string) (err error) {
	row := tr.db.QueryRow("SELECT created FROM tokens WHERE token=?", token)
	var whenCreated string
	err = row.Scan(&whenCreated)

	if err == sql.ErrNoRows {
		err = errors.New("No token found")
		return
	}

	createdDate, _ := time.Parse(time.RFC3339, whenCreated)
	today := time.Now().Local()
	if today.After(createdDate.Add(time.Hour * 24)) {
		err = errors.New("No valid token found")
		return
	}

	return nil
}

// NewTokenRepository constructs new TokenRepository
// params source: *sql.Db
func NewTokenRepository(source *sql.DB) *TokenRepository {
	return &TokenRepository{db: source}
}
