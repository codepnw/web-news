package models

import (
	"errors"
	"time"

	"github.com/upper/db/v4"
	"golang.org/x/crypto/bcrypt"
)

const passwordCost = 12

type User struct {
	ID        int       `db:"id,omitempty"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Password  string    `db:"password_hash"`
	CreatedAt time.Time `db:"created_at"`
	Activated bool      `db:"activated"`
}

type UsersModel struct {
	db db.Session
}

func (u *UsersModel) Table() string {
	return "users table"
}

func (u *UsersModel) Get(id int) (*User, error) {
	var user User

	err := u.db.Collection(u.Table()).Find(db.Cond{"id": id}).One(&user)
	if err != nil {
		if errors.Is(err, db.ErrNoMoreRows) {
			return nil, ErrNoMoreRows
		}
		return nil, err
	}

	return &user, nil
}

func (u *UsersModel) FindByEmail(email string) (*User, error) {
	var user User

	err := u.db.Collection(u.Table()).Find(db.Cond{"email": email}).One(&user)
	if err != nil {
		if errors.Is(err, db.ErrNoMoreRows) {
			return nil, ErrNoMoreRows
		}
		return nil, err
	}

	return &user, nil
}

func (u *UsersModel) Insert(user *User) error {
	newHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), passwordCost)
	if err != nil {
		return err
	}

	user.Password = string(newHash)
	user.CreatedAt = time.Now()

	col := u.db.Collection(u.Table())
	res, err := col.Insert(user)

	if err != nil {
		switch {
		case errHashDuplicate(err, "users_email_key"):
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	user.ID = convertUpperIDtoInt(res.ID())
	return nil
}

func (user *User) ComparePassword(plainPassword string) (bool, error) {

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(plainPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func (u *UsersModel) Authenticate(email, password string) (*User, error) {
	user, err := u.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	if !user.Activated {
		return nil, ErrUserNotActive
	}

	match, err := user.ComparePassword(password)
	if err != nil {
		return nil, err
	}

	if !match {
		return nil, ErrInvalidLogin
	}

	return user, nil
}