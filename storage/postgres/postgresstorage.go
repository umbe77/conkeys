package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"conkeys/storage"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewKeyStorage(db *sql.DB) *PostgresStorage {
	store := &PostgresStorage{
		db: db,
	}
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS keys (
		Key VARCHAR PRIMARY KEY NOT NULL,
		Value JSON NOT NULL
	);
	CREATE TABLE IF NOT EXISTS EncryptedKeys (
		Key VARCHAR PRIMARY KEY NOT NULL,
		Value VARCHAR NOT NULL
	);`)

	if err != nil {
		log.Fatal(err)
	}

	return store

}

func (s *PostgresStorage) Get(path string) (storage.Value, error) {
	stmt, err := s.db.Prepare("SELECT value FROM keys WHERE key = $1")
	if err != nil {
		return storage.Value{}, err
	}
	rows, qErr := stmt.Query(path)
	if qErr != nil {
		return storage.Value{}, qErr
	}
	defer rows.Close()

	if rows.Next() {
		var buf []byte
		sErr := rows.Scan(&buf)
		if sErr != nil {
			return storage.Value{}, sErr
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

func (s *PostgresStorage) GetEncrypted(path string) (storage.Value, error) {
	stmt, err := s.db.Prepare("SELECT value FROM EncryptedKeys WHERE key = $1")
	if err != nil {
		return storage.Value{}, err
	}
	rows, qErr := stmt.Query(path)
	if qErr != nil {
		return storage.Value{}, qErr
	}
	defer rows.Close()

	if rows.Next() {
		var encryptedValue string
		sErr := rows.Scan(&encryptedValue)
		if sErr != nil {
			return storage.Value{}, sErr
		}
		v := storage.Value{
			T: storage.Crypted,
			V: encryptedValue,
		}

		return v, nil
	}
	return storage.Value{}, errors.New("no key found")
}

func (s *PostgresStorage) GetKeys(pathSearch string) (map[string]storage.Value, error) {
	result := make(map[string]storage.Value)

	stmt, err := s.db.Prepare("SELECT key, value FROM keys WHERE key LIKE $1")
	if err != nil {
		return nil, err
	}
	rows, qErr := stmt.Query(fmt.Sprintf("%s%%", pathSearch))
	if qErr != nil {
		return nil, qErr
	}

	defer rows.Close()
	for rows.Next() {
		var key string
		var buf []byte
		sErr := rows.Scan(&key, &buf)
		if sErr != nil {
			return nil, sErr
		}
		v := storage.Value{}
		uErr := json.Unmarshal(buf, &v)
		if uErr != nil {
			return nil, sErr
		}
		result[key] = v
	}

	return result, nil
}

func (s *PostgresStorage) GetAllKeys() (map[string]storage.Value, error) {
	result := make(map[string]storage.Value)
	rows, qErr := s.db.Query("SELECT key, value FROM keys")
	if qErr != nil {
		return nil, qErr
	}

	defer rows.Close()
	for rows.Next() {
		var key string
		var buf []byte
		sErr := rows.Scan(&key, &buf)
		if sErr != nil {
			continue
			//TODO: Set logging
		}
		v := storage.Value{}
		uErr := json.Unmarshal(buf, &v)
		if uErr != nil {
			continue
			//TODO: Set logging
		}
		result[key] = v
	}

	return result, nil
}

func putKey(path string, val storage.Value, tx *sql.Tx) error {

	stmt, err := tx.Prepare(`INSERT INTO keys
	(key, value)
	VALUES
	($1, $2)
	ON CONFLICT (key)
	DO
		UPDATE SET value = $2`)
	if err != nil {
		return err
	}
	v, mErr := json.Marshal(val)
	if mErr != nil {
		return mErr
	}
	_, iErr := stmt.Exec(path, v)
	return iErr
}

func (s *PostgresStorage) Put(path string, val storage.Value) error {
	ctx := context.Background()
	tx, tErr := s.db.BeginTx(ctx, nil)
	if tErr != nil {
		tx.Rollback()
		return tErr
	}
	pErr := putKey(path, val, tx)
	if pErr != nil {
		tx.Rollback()
		return pErr
	}
	tx.Commit()
	return nil
}

func (s *PostgresStorage) PutEncrypted(path string, maskedValue storage.Value, encryptedValue string) error {
	ctx := context.Background()
	tx, tErr := s.db.BeginTx(ctx, nil)
	if tErr != nil {
		return tErr
	}

	err := putKey(path, maskedValue, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	stmt, err := tx.Prepare(`INSERT INTO EncryptedKeys
	(key, value)
	VALUES
	($1, $2)
	ON CONFLICT (key)
	DO
		UPDATE SET value = $2`)

	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = stmt.Exec(path, encryptedValue)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
func (s *PostgresStorage) Delete(path string) error {
	ctx := context.Background()
	tx, tErr := s.db.BeginTx(ctx, nil)
	if tErr != nil {
		return tErr
	}

	stmtK, err := tx.Prepare(`DELETE FROM keys WHERE key = $1;`)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, dErr := stmtK.Exec(path)
	if dErr != nil {
		tx.Rollback()
		return dErr
	}

	stmtE, err := tx.Prepare(`DELETE FROM EncryptedKeys WHERE key = $1;`)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, dErr = stmtE.Exec(path)
	if dErr != nil {
		tx.Rollback()
		return dErr
	}

	tx.Commit()
	return nil

}
