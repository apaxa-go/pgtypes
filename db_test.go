package pgtypes

import (
	"database/sql"
	"github.com/apaxa-go/helper/strconvh"
	"github.com/jackc/pgx"
	_ "github.com/lib/pq"
	"os"
	"testing"
)

// Following environment variables used for testing: PG_HOST, PG_PORT, PG_DATABASE, PG_USER & PG_PASSWORD
const (
	pgHost     = "PG_HOST"
	pgPort     = "PG_PORT"
	pgDatabase = "PG_DATABASE"
	pgUser     = "PG_USER"
	pgPassword = "PG_PASSWORD"
)

var pgxConn *pgx.Conn
var pgxConf pgx.ConnConfig
var pqConn *sql.DB

func setupPgx() {
	var err error
	pgxConf.Host = os.Getenv(pgHost)
	pgxConf.Port, _ = strconvh.ParseUint16(os.Getenv(pgPort))
	pgxConf.Database = os.Getenv(pgDatabase)
	pgxConf.User = os.Getenv(pgUser)
	pgxConf.Password = os.Getenv(pgPassword)

	pgxConn, err = pgx.Connect(pgxConf)
	if err != nil {
		panic(err)
	}
}

func cleanPgx() {
	if pgxConn != nil {
		_ = pgxConn.Close()
	}
}

func setupPq() {
	// compute connection string
	addPart := func(argName, envName string) string {
		s := os.Getenv(envName)
		if s == "" {
			return ""
		}
		return " " + argName + "=" + s
	}
	connStr := "sslmode=disable" + addPart("host", pgHost) + addPart("port", pgPort) + addPart("user", pgUser) + addPart("password", pgPassword) + addPart("dbname", pgDatabase)

	//
	var err error
	pqConn, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
}

func cleanPq() {
	if pqConn != nil {
		_ = pqConn.Close()
	}
}

func TestMain(m *testing.M) {
	main := func() int {
		setupPgx()
		defer cleanPgx()
		setupPq()
		defer cleanPq()
		return m.Run()
	}
	r := main()
	os.Exit(r)
}
