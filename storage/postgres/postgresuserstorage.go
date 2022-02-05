package postgres

import (
	"conkeys/config"
	"conkeys/storage"
	"database/sql"
	"errors"
	"log"
)

type PostgresUserStorage struct{}

var dbUsers *sql.DB

func (s PostgresUserStorage) Init() {
	cfg := config.GetConfig()
	connectionUri := "postgres://conkeys:S0jeje1!@localhost/conkeys?sslmode=disable"
	if cfg.Postgres.ConnectionUri != "" {
		connectionUri = cfg.Postgres.ConnectionUri
	}

	var err error
	dbUsers, err = sql.Open("postgres", connectionUri)
	if err != nil {
		log.Fatal(err)
	}
	err = dbUsers.Ping()
	if err != nil {
		log.Fatal(err)
	}

	_, cErr := dbUsers.Exec(`CREATE TABLE IF NOT EXISTS users (
		UserName VARCHAR PRIMARY KEY NOT NULL,
		FirstName VARCHAR NULL,
        LastName VARCHAR NULL,
        Email VARCHAR,
        "Password" VARCHAR
	)`)
	if cErr != nil {
		log.Fatal(cErr)
	}
}

func (s PostgresUserStorage) Get(userName string) (storage.User, error) {
	stmt, err := dbUsers.Prepare("SELECT UserName, FirstName, LastName, Email, Password FROM users WHERE UserName = $1")
	if err != nil {
		return storage.User{}, nil
	}
	
	rows, qErr := stmt.Query(userName)
	if qErr != nil {
		return storage.User{}, qErr
	}
	
	defer rows.Close()

	if rows.Next() {
		var name, lastName, email, password string
		sErr := rows.Scan(&name, &lastName, &email, &password)
		if sErr != nil {
			return storage.User{}, sErr
		}
		return storage.User{
			UserName: userName,
			Name: name,
			LastName: lastName,
			Email: email,
			Password: password,
		}, nil
	}
	return storage.User{}, errors.New("no user found")
}

// TODO: Add Query Logic
func (s PostgresUserStorage) GetUsers(query string) ([]storage.User, error) {
	result := make([]storage.User, 0)

	rows, qErr := dbUsers.Query("SELECT UserName, FirstName, LastName, Email FROM users")
	if qErr != nil {
		return result, qErr
	}
	
	for rows.Next() {
		var userName, name, lastName, email string
		sErr := rows.Scan(&userName, &name, &lastName, &email)
		if sErr != nil {
			return result, sErr
		}

		result = append(result, storage.User{
			UserName: userName,
			Name: name,
			LastName: lastName,
			Email: email,
		})
	}

	return result, nil
}

func (s PostgresUserStorage) Add(usr storage.User) error {
	stmt, err := dbUsers.Prepare(`INSERT INTO users
	(UserName, FirstName, LastName, Email)
	VALUES
	($1, $2, $3, $4)`)
	if err != nil {
		return err
	}
	
	_, iErr := stmt.Exec(usr.UserName, usr.Name, usr.LastName, usr.Email)
	return iErr
}

func (s PostgresUserStorage) Update(usr storage.User) error {
	stmt, err := dbUsers.Prepare(`UPDATE users
	SET FirstName = $2, LastName = $3, Email = $4
	WHERE UserName = $1`)

	if err != nil {
		return err
	}
	
	_, iErr := stmt.Exec(usr.UserName, usr.Name, usr.LastName, usr.Email)
	return iErr
}

func (s PostgresUserStorage) Delete(userName string) error {
	stmt, err := dbUsers.Prepare(`DELETE FROM users WHERE UserName = $1`)
	if err != nil {
		return err
	}
	
	_, iErr := stmt.Exec(userName)
	return iErr
}

func (s PostgresUserStorage) SetPassword(userName string, password string) {
	stmt, _ := dbUsers.Prepare(`UPDATE users
	SET Password = $2
	WHERE UserName = $1`)

	// if err != nil {
	// 	return err
	// }
	
	stmt.Exec(userName, password)
	// return iErr
}
