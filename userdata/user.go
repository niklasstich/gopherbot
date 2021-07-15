package userdata

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"math"
	"time"
)

type DBUser struct {
	ID                    string
	Points                int64
	LastCurrencyClaimTime time.Time
}

var ErrorUserAlreadyExists = errors.New("user with that ID already exists")
var ErrorUserDoesNotExist = errors.New("no user with that ID exists in the database")
var ErrorUserNotEnoughPoints = errors.New("user doesn't have enough points")

//GetUser returns a valid *DBUser object for parameter discordId, which must be a valid discord user id.
//Will return nil if user is not found in the database
func GetUser(discordId string) (user *DBUser) {
	user = getUserDb(discordId)
	return
}

//CreateUser will create a new entry for discordId with sane defaults and return that a new *DBUser representing that.
//Will return an additional error if a user with that ID already exists
func CreateUser(discordId string) (user *DBUser, err error) {
	user = getUserDb(discordId)
	if user != nil {
		return user, ErrorUserAlreadyExists
	}
	err = createUserDb(discordId)
	if err != nil {
		log.Errorf("failed to create db user for id %s", discordId)
		return nil, err
	}

	//return new user object
	user = &DBUser{
		ID:     discordId,
		Points: 0,
	}
	return
}

//WriteToDB writes the user to the database by updating the entry for it.
//Will return an error if no user with that ID exists
func (user DBUser) WriteToDB() (err error) {
	if query := getUserDb(user.ID); query == nil {
		return ErrorUserDoesNotExist
	}
	err = updateUserDb(user)
	if err != nil {
		log.Errorf("failed to update db user for id %s", user.ID)
	}
	return
}

func (user *DBUser) AddPoints(i int64) {
	if (user.Points + i) > 0 {
		user.Points += i
	} else {
		user.Points = math.MaxInt64
	}
}

func (user *DBUser) CanAfford(i int64) bool {
	return user.Points >= i
}

func (user *DBUser) RemovePoints(i int64) (err error) {
	if !user.CanAfford(i) {
		return ErrorUserNotEnoughPoints
	}
	user.Points -= i
	return nil
}

func (user *DBUser) SetCurrencyClaimTime(t time.Time) {
	user.LastCurrencyClaimTime = t
}
