package postgres

import (
	"conkeys/config"
	"conkeys/crypto"
	"crypto/rsa"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type SecurityPostgresStorage struct{}

var dbSecurity *sql.DB

func (s SecurityPostgresStorage) Init() {
	cfg := config.GetConfig()
	connectionUri := "postgres://conkeys:S0jeje1!@localhost/conkeys?sslmode=disable"
	if cfg.Postgres.ConnectionUri != "" {
		connectionUri = cfg.Postgres.ConnectionUri
	}

	var err error
	dbSecurity, err = sql.Open("postgres", connectionUri)
	if err != nil {
		log.Fatal(err)
	}
	err = dbUsers.Ping()
	if err != nil {
		log.Fatal(err)
	}

	_, cErr := dbSecurity.Exec(`CREATE TABLE IF NOT EXISTS SigningPair (
		public VARCHAR NOT NULL,
		"private" VARCHAR NOT NULL
	);
	CREATE TABLE IF NOT EXISTS CryptoPair (
		public VARCHAR NOT NULL,
		"private" VARCHAR NOT NULL
	);`)
	if cErr != nil {
		log.Fatal(cErr)
	}
}

func loadPair(tableName string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	rows, qErr := dbSecurity.Query(fmt.Sprintf("SELECT public, private FROM %s", tableName))
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

func savePair(tableName string, priv *rsa.PrivateKey, pub *rsa.PublicKey) error {
	private := crypto.PrivateKeyToBytes(priv)
	public, err := crypto.PublicKeyToBytes(pub)
	if err != nil {
		return err
	}

	stmt, sErr := dbSecurity.Prepare(fmt.Sprintf(`INSERT INTO %s
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
func (s SecurityPostgresStorage) LoadCryptingPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	return loadPair("CryptoPair")
}

func (s SecurityPostgresStorage) SaveCryptingPair(priv *rsa.PrivateKey, pub *rsa.PublicKey) error {
	return savePair("CryptoPair", priv, pub)
}

func (s SecurityPostgresStorage) LoadSigninPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	return loadPair("SigningPair")
}

func (s SecurityPostgresStorage) SaveSigninPair(priv *rsa.PrivateKey, pub *rsa.PublicKey) error {
	return savePair("SigningPair", priv, pub)
}
