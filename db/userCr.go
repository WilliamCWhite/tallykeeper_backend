package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// from context and an email, retrieves the user_id from the database
// returns -1 if the email doesn't exist in the db
// returns -2 for other errors
func GetUserIDByEmail(ctx context.Context, email string) (int, error) {
	query := `SELECT user_id
		FROM users WHERE email = $1`

	var foundUserID int
	err := pool.QueryRow(ctx, query, email).Scan(&foundUserID)
	// "no rows in result set"
	if errors.Is(err, pgx.ErrNoRows) {
		return -1, nil
	} else if err != nil {
		return -2, fmt.Errorf("failed to query users from db with pgx: %w", err)
	}

	return foundUserID, nil
}

// from context and an email, creates a new user in the database, returning its new user_id
func CreateUserFromEmail(ctx context.Context, email string) (int, error) {
	query := `INSERT INTO users (email)
		VALUES ($1)
		RETURNING user_id`

	var newUserID int
	err := pool.QueryRow(ctx, query, email).Scan(&newUserID)

	if err != nil {
		return -1, fmt.Errorf("failed to create user with pgx: %w", err)
	}

	return newUserID, nil
}
