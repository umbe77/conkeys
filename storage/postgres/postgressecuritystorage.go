package postgres

import (
	"crypto/rsa"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"conkeys/crypto"
)

type PostgresSecurityStorage struct {
	db *sql.DB
}

func NewSecStorage(db *sql.DB) *PostgresSecurityStorage {
	store := &PostgresSecurityStorage{
		db: db,
	}
	_, err := store.db.Exec(`CREATE TABLE IF NOT EXISTS SigningPair (
		public VARCHAR NOT NULL,
		"private" VARCHAR NOT NULL
	);
	CREATE TABLE IF NOT EXISTS CryptoPair (
		public VARCHAR NOT NULL,
		"private" VARCHAR NOT NULL
	);`)

	if err != nil {
		log.Fatal(err)
	}

	return store
}

func (s *PostgresSecurityStorage) loadPair(tableName string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	rows, qErr := s.db.Query(fmt.Sprintf("SELECT public, private FROM %s", tableName))
	if qErr != nil {
		return nil, nil, qErr
	}
	defer rows.Close()

	if rows.Next() {
		var private, public string
		sErr := rows.Scan(&public, &private)
		if sErr != nil {
			return nil, nil, sErr
		}

		priv, err := crypto.BytesToPrivateKey([]byte(private))
		if err != nil {
			return nil, nil, err
		}

		pub, err := crypto.BytesToPublicKey([]byte(public))
		if err != nil {
			return nil, nil, err
		}

		return priv, pub, nil
	}
	return nil, nil, errors.New("No pair found")
}

func (s *PostgresSecurityStorage) savePair(tableName string, priv *rsa.PrivateKey, pub *rsa.PublicKey) error {
	private := crypto.PrivateKeyToBytes(priv)
	public, err := crypto.PublicKeyToBytes(pub)
	if err != nil {
		return err
	}

	stmt, sErr := s.db.Prepare(fmt.Sprintf(`INSERT INTO %s
	(public, private)
	VALUES
	($1, $2)`, tableName))
	if sErr != nil {
		return sErr
	}

	_, iErr := stmt.Exec(string(public), string(private))
	if iErr != nil {
		return iErr
	}

	return nil
}
func (s *PostgresSecurityStorage) LoadCryptingPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	return s.loadPair("CryptoPair")
}

func (s *PostgresSecurityStorage) SaveCryptingPair(priv *rsa.PrivateKey, pub *rsa.PublicKey) error {
	return s.savePair("CryptoPair", priv, pub)
}

func (s *PostgresSecurityStorage) LoadSigninPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	return s.loadPair("SigningPair")
}

func (s *PostgresSecurityStorage) SaveSigninPair(priv *rsa.PrivateKey, pub *rsa.PublicKey) error {
	return s.savePair("SigningPair", priv, pub)
}
