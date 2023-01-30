package general

import (
	"database/sql"

	_ "github.com/lib/pq"
)

const MAX_OPEN_CONNS = 250

type Pgsql struct {
	*sql.DB
}

func NewPgsql(connectionString string) *Pgsql {
	db, err := sql.Open("postgres", connectionString)
	db.SetMaxOpenConns(MAX_OPEN_CONNS)
	if err != nil {
		panic(err)
	}
	return &Pgsql{db}
}
