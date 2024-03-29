package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"conkeys/storage"
)

type PostgresUserStorage struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) *PostgresUserStorage {
	store := &PostgresUserStorage{
		db: db,
	}

	_, err := store.db.Exec(`CREATE TABLE IF NOT EXISTS users (
		username VARCHAR PRIMARY KEY NOT NULL,
		firstname VARCHAR NULL,
        lastname VARCHAR NULL,
        email VARCHAR,
        "password" VARCHAR,
		isAdmin boolean
	)`)
	if err != nil {
		log.Fatal(err)
	}
	return store
}

func (s *PostgresUserStorage) Get(userName string) (storage.User, error) {
	stmt, err := s.db.Prepare("SELECT FirstName, LastName, Email, IsAdmin FROM users WHERE UserName = $1")
	if err != nil {
		return storage.User{}, nil
	}

	rows, qErr := stmt.Query(userName)
	if qErr != nil {
		return storage.User{}, qErr
	}

	defer rows.Close()

	if rows.Next() {
		var name, lastName, email string
		var isAdmin bool
		sErr := rows.Scan(&name, &lastName, &email, &isAdmin)
		if sErr != nil {
			return storage.User{}, sErr
		}
		return storage.User{
			UserName: userName,
			Name:     name,
			LastName: lastName,
			Email:    email,
			IsAdmin:  isAdmin,
		}, nil
	}
	return storage.User{}, errors.New("no user found")
}

func getAllUsers(db *sql.DB) (*sql.Rows, error) {
	return db.Query("SELECT UserName, FirstName, LastName, Email, IsAdmin FROM Users")
}

func getUsersByQuery(db *sql.DB, query string) (*sql.Rows, error) {
	stmt, err := db.Prepare("SELECT UserName, FirstName, LastName, Email, IsAdmin FROM Users WHERE UserName like $1")

	if err != nil {
		return nil, err
	}

	return stmt.Query(fmt.Sprintf("%%%s%%", query))
}

func (s *PostgresUserStorage) GetUsers(query string) ([]storage.User, error) {
	result := make([]storage.User, 0)

	rows, qErr := (func() (*sql.Rows, error) {
		if len(query) > 0 {
			return getUsersByQuery(s.db, query)
		}
		return getAllUsers(s.db)
	})()

	if qErr != nil {
		return result, qErr
	}

	for rows.Next() {
		var userName, name, lastName, email string
		var isAdmin bool
		sErr := rows.Scan(&userName, &name, &lastName, &email, &isAdmin)
		if sErr != nil {
			return result, sErr
		}

		result = append(result, storage.User{
			UserName: userName,
			Name:     name,
			LastName: lastName,
			Email:    email,
			IsAdmin:  isAdmin,
		})
	}

	return result, nil
}

func (s *PostgresUserStorage) Add(usr storage.User) error {
	stmt, err := s.db.Prepare(`INSERT INTO users
	(UserName, FirstName, LastName, Email, IsAdmin)
	VALUES
	($1, $2, $3, $4, $5)`)
	if err != nil {
		return err
	}

	_, iErr := stmt.Exec(usr.UserName, usr.Name, usr.LastName, usr.Email, usr.IsAdmin)
	return iErr
}

func (s *PostgresUserStorage) Update(usr storage.User) error {
	stmt, err := s.db.Prepare(`UPDATE users
	SET FirstName = $2, LastName = $3, Email = $4, IsAdmin = $5
	WHERE UserName = $1`)

	if err != nil {
		return err
	}

	_, iErr := stmt.Exec(usr.UserName, usr.Name, usr.LastName, usr.Email, usr.IsAdmin)
	return iErr
}

func (s *PostgresUserStorage) Delete(userName string) error {
	stmt, err := s.db.Prepare(`DELETE FROM users WHERE UserName = $1`)
	if err != nil {
		return err
	}

	_, iErr := stmt.Exec(userName)
	return iErr
}

func (s *PostgresUserStorage) SetPassword(userName string, password string) error {
	fmt.Printf("Username: %s\nPassword: %s\n", userName, password)
	stmt, err := s.db.Prepare(`UPDATE users
	SET Password = $2
	WHERE UserName = $1`)

	if err != nil {
		return err
	}

	_, iErr := stmt.Exec(userName, password)
	return iErr
}

func (s *PostgresUserStorage) GetPassword(userName string) (string, storage.User, error) {
	stmt, err := s.db.Prepare("SELECT Password, UserName, FirstName, LastName, Email, IsAdmin FROM Users WHERE UserName = $1")
	if err != nil {
		return "", storage.User{}, err
	}
	rows, qErr := stmt.Query(userName)
	if qErr != nil {
		return "", storage.User{}, qErr
	}

	defer rows.Close()

	if rows.Next() {
		var password, userName, firstName, lastName, email string
		var isAdmin bool
		sErr := rows.Scan(&password, &userName, &firstName, &lastName, &email, &isAdmin)
		if sErr != nil {
			return "", storage.User{}, sErr
		}
		return password, storage.User{
			UserName: userName,
			Name:     firstName,
			LastName: lastName,
			Email:    email,
			IsAdmin:  isAdmin,
		}, nil
	}
	return "", storage.User{}, nil
}
