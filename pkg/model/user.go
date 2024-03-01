package model

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (u UserModel) Insert(user *User) error {
	query := `INSERT INTO users (name, email, password)
			VALUES($1, $2, $3)
			RETURNING id`
	args := []interface{}{user.ID, user.Name, user.Email, user.Password}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	return u.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID)

}

func (u UserModel) GetAll() ([]*User, error) {
	query := `
		SELECT id,name, email, password
		FROM users
		ORDER BY id
	`

	rows, err := u.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID,
			&user.Name, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (u UserModel) GetById(id int) (*User, error) {
	query := `
		SELECT id, name, email, password
		FROM users
		WHERE id = $1
		`
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := u.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
func (u UserModel) Update(user *User) error {
	query := `
		UPDATE users
		SET profilePhoto = $1, name = $2, username = $3, description = $4, email = $5, password = $6
		WHERE id = $7
		RETURNING updatedAt
		`

	args := []interface{}{user.ID, user.Name, user.Email, user.Password}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return u.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID)
}

func (u UserModel) Delete(id int) error {
	query := `
		DELETE FROM users
		WHERE id = $1
		`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	_, err := u.DB.ExecContext(ctx, query, id)

	return err
}
