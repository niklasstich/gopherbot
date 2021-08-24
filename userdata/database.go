package userdata

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"time"
)

var db *sql.DB

var ErrorNoRowsUpdated = errors.New("updating user yielded no updated rows")

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./sqlite.db")
	if err != nil {
		log.Fatal("Failed to open DB file: ", err.Error())
	}

	//create table for users
	createQuery := `CREATE TABLE IF NOT EXISTS users (
	id text PRIMARY KEY,
	points integer NOT NULL DEFAULT 0,
	LastCurrencyClaimTime timestamp NOT NULL
	);`

	_, err = db.Exec(createQuery)
	if err != nil {
		log.Fatal("Failed to create user table: ", err.Error())
	}
	log.Debug("Initialized DB")
}

func createUserDb(id string) (err error) {
	//open tx
	tx, err := db.Begin()
	if err != nil {
		log.Error("Failed to open DB transaction: ", err.Error())
		return
	}
	//prepare statement
	stmtUserAdd, err := tx.Prepare("INSERT INTO users(id,points, LastCurrencyClaimTime) VALUES(?, ?, ?);")
	if err != nil {
		log.Error("Failed to prepare DB new user statement: ", err.Error())
		tx.Rollback()
		return
	}
	defer stmtUserAdd.Close()
	//execute statement
	_, err = stmtUserAdd.Exec(id, 0, time.Now().AddDate(0, 0, -1))
	if err != nil {
		log.Error("Failed to execute DB new user statement: ", err.Error())
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Error("Failed to commit DB new user statement: ", err.Error())
		return
	}
	return
}

//getUserDb returns the *DBUser for given id found in the database.
//If no entry exists, it will simply return nil.
func getUserDb(id string) (user *DBUser) {
	user = &DBUser{}
	res, err := db.Query("SELECT * FROM users WHERE id = (?);", id)
	if err != nil {
		log.Error("Failed to execute DB GetUser query: ", err.Error())
		return nil
	}
	defer res.Close()

	//return nil if no user was found
	if !res.Next() {
		return nil
	}
	err = res.Scan(&user.ID, &user.Points, &user.LastCurrencyClaimTime)
	if err != nil {
		log.Error("Failed to read row from DB GetUser query: ", err.Error())
	}

	//return valid user
	return
}

func updateUserDb(user DBUser) (err error) {
	tx, err := db.Begin()
	if err != nil {
		log.Error("Failed to open DB transaction: ", err.Error())
		return
	}
	stmtUserUpdate, err := tx.Prepare(`UPDATE users 
					SET points = ?,LastCurrencyClaimTime = ?
					WHERE id = ?;`)
	if err != nil {
		log.Error("Failed to prepare DB UserUpdate statement: ", err.Error())
		tx.Rollback()
		return
	}
	defer stmtUserUpdate.Close()

	res, err := stmtUserUpdate.Exec(user.Points, user.LastCurrencyClaimTime, user.ID)
	if err != nil {
		log.Error("Failed to execute DB UpdateUser statement: ", err.Error())
		tx.Rollback()
		return
	}
	affected, _ := res.RowsAffected()
	if affected < 1 {
		return ErrorNoRowsUpdated
	}

	//commit tx
	err = tx.Commit()
	if err != nil {
		log.Error("Failed to commit DB UpdateUser TX: ", err.Error())
	}
	return nil
}

//EnsureDBClosed makes sure the connection to the database is closed properly before quitting gracefully
func EnsureDBClosed() {
	if db != nil {
		err := db.Close()
		if err != nil {
			log.Error("Failed closing DB connection gracefully: ", err.Error())
		}
	}
}

func UserCount() (count int64) {
	res, err := db.Query("SELECT COUNT (id) FROM main.users;")
	if err != nil {
		log.Error("Failed to query user count: ", err.Error())
		return 0
	}
	defer res.Close()
	res.Next()
	err = res.Scan(&count)
	if err != nil {
		log.Error("Failed to read user count: ", err.Error())
		return 0
	}
	return
}
