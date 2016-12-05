package pgtypes

import (
	"github.com/apaxa-go/helper/strconvh"
	"github.com/jackc/pgx"
	"os"
	"testing"
)

var conn *pgx.Conn
var conf pgx.ConnConfig

func setupDB() {
	var err error
	conf.Host = os.Getenv("PG_HOST")
	conf.Port, _ = strconvh.ParseUint16(os.Getenv("PG_PORT"))
	conf.Database = os.Getenv("PG_DATABASE")
	conf.User = os.Getenv("PG_USER")
	conf.Password = os.Getenv("PG_PASSWORD")

	conn, err = pgx.Connect(conf)
	if err != nil {
		panic(err)
	}
}

func cleanDB() {
	if conn != nil {
		conn.Close()
	}
}

func TestMain(m *testing.M) {
	main := func() int {
		setupDB()
		defer cleanDB()
		return m.Run()
	}
	r := main()
	os.Exit(r)
}
