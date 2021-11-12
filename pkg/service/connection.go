package service

import (
	"context"
	"database/sql"
	"github.com/sisukasco/henki/pkg/db"

	_ "github.com/lib/pq" //required for postgres
	"github.com/pkg/errors"
)

type DBConnection struct {
	db *sql.DB
	Q  *db.Queries
}

func NewConnection(dburl string) (*DBConnection, error) {

	if len(dburl) <= 0 {
		return nil, errors.New("DBURL config is not set!")
	}

	dbconn, err := sql.Open("postgres", dburl)
	if err != nil {
		return nil, errors.Wrap(err, "Making new DB Connection")
	}

	q := db.New(dbconn)

	conn := DBConnection{db: dbconn, Q: q}

	return &conn, nil
}

func (this *DBConnection) Ping(ctx context.Context) error {
	return this.db.PingContext(ctx)
}

func (this *DBConnection) Close() {

	this.db.Close()
}

func (this *DBConnection) TruncateAll() {
	this.db.Exec("TRUNCATE users")
}
