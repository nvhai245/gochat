package user

import (
	"time"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/nvhai245/gochat/services/auth/model"
)

func GetUser(username string, db *sqlx.DB) (success bool, authorizedUser model.AuthorizedUser) {
	sqlStatement := `SELECT FROM users WHERE username = $1`
	err := db.Get(&authorizedUser, sqlStatement, username)
	if err != nil {
		log.Println(err)
		return false, model.AuthorizedUser{}
	}
	return true, authorizedUser
}

func UpdateEmail(username string, email string, db *sqlx.DB) (success bool, updatedEmail string) {
	sqlStatement := `UPDATE users
	                 SET email = $1
					 WHERE username = $2
					 RETURNING email`
	err := db.QueryRow(sqlStatement, email, username).Scan(&updatedEmail)
	if err != nil {
		log.Println(err)
		return false, ""
	}
	return true, updatedEmail
}

func UpdateAvatar(username string, avatar string, db *sqlx.DB) (success bool, updatedAvatar string) {
	sqlStatement := `UPDATE users
	                 SET avatar = $1
					 WHERE username = $2
					 RETURNING avatar`
	err := db.QueryRow(sqlStatement, avatar, username).Scan(&updatedAvatar)
	if err != nil {
		log.Println(err)
		return false, ""
	}
	return true, updatedAvatar
}

func UpdatePhone(username string, phone string, db *sqlx.DB) (success bool, updatedPhone string) {
	sqlStatement := `UPDATE users
	                 SET phone = $1
					 WHERE username = $2
					 RETURNING phone`
	err := db.QueryRow(sqlStatement, phone, username).Scan(&updatedPhone)
	if err != nil {
		log.Println(err)
		return false, ""
	}
	return true, updatedPhone
}

func UpdateBirthday(username string, birthday time.Time, db *sqlx.DB) (success bool, updatedBirthday time.Time) {
	sqlStatement := `UPDATE users
	                 SET birthday = $1
					 WHERE username = $2
					 RETURNING birthday`
	err := db.QueryRow(sqlStatement, avatar, username).Scan(&updatedBirthday)
	if err != nil {
		log.Println(err)
		return false, birthday
	}
	return true, updatedBirthday
}

func UpdateFb(username string, fb string, db *sqlx.DB) (success bool, updatedFb string) {
	sqlStatement := `UPDATE users
	                 SET fb = $1
					 WHERE username = $2
					 RETURNING fb`
	err := db.QueryRow(sqlStatement, fb, username).Scan(&updatedFb)
	if err != nil {
		log.Println(err)
		return false, ""
	}
	return true, updatedFb
}

func UpdateInsta(username string, insta string, db *sqlx.DB) (success bool, updatedInsta string) {
	sqlStatement := `UPDATE users
	                 SET insta = $1
					 WHERE username = $2
					 RETURNING insta`
	err := db.QueryRow(sqlStatement, insta, username).Scan(&updatedInsta)
	if err != nil {
		log.Println(err)
		return false, ""
	}
	return true, updatedInsta
}