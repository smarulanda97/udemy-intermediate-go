package models

import (
	"context"
	"crypto/sha256"
	"strings"
	"time"
)

// User is the type for users
type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// GetUserByEmail gets a user by email address
func (dbm *DBModels) GetUserByEmail(email string) (User, error) {
	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	email = strings.ToLower(email)

	row := dbm.DB.QueryRowContext(ctx, `
		select 
			id, first_name, last_name, email, password, created_at, updated_at
		from 
			users
		where email = ?
	`, email)

	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (dbm *DBModels) GetUserByToken(token string) (*User, error) {
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tokenHash := sha256.Sum256([]byte(token))

	stmt := `
		select
			u.id, u.first_name, u.last_name, u.email
		from 
			users u
		inner join tokens t
			on (u.id = t.user_id)
		where
			t.token_hash = ?
			and t.expiry > ?
	`

	err := dbm.DB.QueryRowContext(ctx, stmt, tokenHash[:], time.Now()).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
