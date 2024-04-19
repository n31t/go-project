package model

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserModel struct {
	DB       *sql.DB
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}

func (u *UserModel) Insert(user *User) error {
	query := `INSERT INTO users (username, password, email)
	VALUES ($1, $2, $3)
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := u.DB.QueryRowContext(ctx, query, user.Username, user.Password, user.Email).Scan(&user.Id)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserModel) Select(id int) (*User, error) {
	query := `
	SELECT id, username, password, email
	FROM users
	WHERE id = $1`
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := u.DB.QueryRowContext(ctx, query, id).Scan(&user.Id, &user.Username, &user.Password, &user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserModel) SelectAll() ([]*User, error) {
	query := `
	SELECT * FROM users`
	rows, err := u.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.Id, &user.Username, &user.Password, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (u *UserModel) Update(user *User) error {
	query := `UPDATE users SET username = $1, password = $2, email = $3 WHERE id = $4`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := u.DB.ExecContext(ctx, query, user.Username, user.Password, user.Email, user.Id)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserModel) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := u.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
