package db

import (
	"fmt"
	"context"
)

func GetUserIDByEmail(ctx context.Context, email string) (int, error) {
	query := `SELECT user_id
		FROM users WHERE email = $1`

	var foundUserID int
	err := pool.QueryRow(ctx, query, email).Scan(&foundUserID)
	if err != nil {
		return -1, fmt.Errorf("failed to query users from db with pgx: %w", err)
	}

	return foundUserID, nil
}

func CreateUserFromEmail(ctx context.Context, email string) (int, error) {
	id, err := GetUserIDByEmail(ctx, email)
	if err == nil {
		return -1, fmt.Errorf("User already exists with id %d", id)
	}

	query := `INSERT INTO users (email)
		VALUES ($1)
		RETURNING user_id`

	var newUserID int
	err = pool.QueryRow(ctx, query, email).Scan(&newUserID)

	if err != nil {
		return -1, fmt.Errorf("failed to create user with pgx: %w", err)
	}

	return newUserID, nil
}
