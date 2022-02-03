package postgres

import (
	"conkeys/config"
	"conkeys/storage"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type PostgresStorage struct{}

var db *sql.DB

func (m PostgresStorage) Init() {
	cfg := config.GetConfig()
	connectionUri := "postgres://conkeys:S0jeje1!@localhost/conkeys?sslmode=disable"
	if cfg.Postgres.ConnectionUri != "" {
		connectionUri = cfg.Postgres.ConnectionUri
	}

	var err error
	db, err = sql.Open("postgres", connectionUri)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	_, cErr := db.Exec(`CREATE TABLE IF NOT EXISTS keys (
		Key VARCHAR PRIMARY KEY NOT NULL,
		Value JSON NOT NULL
	)`)
	if cErr != nil {
		log.Fatal(cErr)
	}

}

func (m PostgresStorage) Get(path string) (storage.Value, error) {
	stmt, err := db.Prepare("SELECT value FROM keys WHERE key = $1")
	if err != nil {
		return storage.Value{}, err
	}
	normalizedPath := strings.TrimPrefix(path, "/")
	rows, qErr := stmt.Query(normalizedPath)
	if qErr != nil {
		return storage.Value{}, qErr
	}
	defer rows.Close()

	if rows.Next() {
		var buf []byte
		sErr := rows.Scan(&buf)
		if sErr != nil {
			return storage.Value{}, qErr
		}
		v := storage.Value{}
		uErr := json.Unmarshal(buf, &v)
		if uErr != nil {
			return storage.Value{}, uErr
		}
		return v, nil
	}
	return storage.Value{}, errors.New("no key found")
}

func (m PostgresStorage) GetKeys(pathSearch string) (map[string]storage.Value, error) {
	result := make(map[string]storage.Value)
	normalizedPathSearch := strings.TrimPrefix(pathSearch, "/")

	stmt, err := db.Prepare("SELECT key, value FROM keys WHERE key LIKE $1")
	if err != nil {
		fmt.Fprintf(gin.DefaultWriter, "%s", err)
	}
	rows, qErr := stmt.Query(fmt.Sprintf("%s%%", normalizedPathSearch))
	if qErr != nil {
		fmt.Fprintf(gin.DefaultWriter, "%s", qErr)
	}

	defer rows.Close()
	for rows.Next() {
		var key string
		var buf []byte
		sErr := rows.Scan(&key, &buf)
		if sErr != nil {
			fmt.Fprintf(gin.DefaultWriter, "%s", sErr)
		}
		v := storage.Value{}
		uErr := json.Unmarshal(buf, &v)
		if uErr != nil {
			fmt.Fprintf(gin.DefaultWriter, "%s", uErr)
		}
		result[key] = v
	}

	return result, nil
}

func (m PostgresStorage) GetAllKeys() map[string]storage.Value {
	result := make(map[string]storage.Value)
	rows, qErr := db.Query("SELECT key, value FROM keys")
	if qErr != nil {
		fmt.Fprintf(gin.DefaultWriter, "%s", qErr)
	}

	defer rows.Close()
	for rows.Next() {
		var key string
		var buf []byte
		sErr := rows.Scan(&key, &buf)
		if sErr != nil {
			fmt.Fprintf(gin.DefaultWriter, "%s", sErr)
		}
		v := storage.Value{}
		uErr := json.Unmarshal(buf, &v)
		if uErr != nil {
			fmt.Fprintf(gin.DefaultWriter, "%s", uErr)
		}
		result[key] = v
	}

	return result
}

func (m PostgresStorage) Put(path string, val storage.Value) {
	stmt, err := db.Prepare(`INSERT INTO keys
	(key, value)
	VALUES
	($1, $2)
	ON CONFLICT (key)
	DO
		UPDATE SET value = $2`)
	if err != nil {
		fmt.Fprintf(gin.DefaultWriter, "%s", err)
	}
	normalizedPath := strings.TrimPrefix(path, "/")
	v, mErr := json.Marshal(val)
	if mErr != nil {
		fmt.Fprintf(gin.DefaultWriter, "%s", err)
	}
	stmt.Exec(normalizedPath, v)

}

func (m PostgresStorage) Delete(path string) {
	stmt, err := db.Prepare("DELETE FROM keys WHERE key = $1")
	if err != nil {
		fmt.Fprintf(gin.DefaultWriter, "%s", err)
	}
	normalizedPath := strings.TrimPrefix(path, "/")
	stmt.Exec(normalizedPath)

}
