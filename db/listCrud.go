package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// from context and a userId, returns an array of all the users lists
func GetListsByUserID(ctx context.Context, userID int) ([]List, error) {
	query := `SELECT list_id, title, time_created, time_modified
		FROM lists WHERE user_id = $1
		ORDER BY time_created DESC`

	rows, err := pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query lists from PostgreSQL with pgx: %w", err)
	}

	defer rows.Close()

	var lists []List
	for rows.Next() {
		var list List
		err := rows.Scan(
			&list.ListID,
			&list.Title,
			&list.TimeCreated,
			&list.TimeModified,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan list row with pgx: %w", err)
		}
		lists = append(lists, list)
	}

	if len(lists) == 0 {
		return []List{}, nil // returns empty slice
	}

	return lists, nil
}

// from context and a list containing a title and userId, creates a list in the database, returning the new list
func CreateList(ctx context.Context, list List) (*List, error) {
	query := `INSERT INTO lists (title, time_created, time_modified, user_id)
		VALUES ($1, $2, $3, $4)
		RETURNING list_id, time_created, time_modified`

	var newListID int
	var newTimeCreated time.Time // the database and backend time modify differ, so must receive the database version
	var newTimeModified time.Time
	err := pool.QueryRow(ctx, query,
		list.Title,
		time.Now(),
		time.Now(),
		list.UserID,
	).Scan(&newListID, &newTimeCreated, &newTimeModified)

	newList := List {
		Title: list.Title,
		ListID: newListID,
		TimeCreated: newTimeCreated,
		TimeModified: newTimeModified,
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create list with pgx: %w", err)
	}

	return &newList, nil
}

// from context and a list containing a title, list_id, and user_id, updates that list in the database
func UpdateList(ctx context.Context, list List) (error) {
	query := `UPDATE lists SET title = $1, time_modified = $2
		WHERE list_id = $3 AND user_id = $4
		RETURNING time_modified`

	commandTag, err := pool.Exec(ctx, query,
		list.Title,
		time.Now(),
		list.ListID,
		list.UserID,
	)
	if err != nil {
		return fmt.Errorf("failed to update list with pgx: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("no list found/updated with ID %d for user %d using pgx:", list.ListID, list.UserID)
	}

	return nil
}

// from context, a list id, and a user id, deletes that list in the database
func DeleteList(ctx context.Context, listID int, userID int) error {
	query := `DELETE FROM lists WHERE list_id = $1 AND user_id = $2`

	commandTag, err := pool.Exec(ctx, query,
		listID,
		userID,
	)

	if err != nil {
		return fmt.Errorf("failed to delete note with pgx: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("no list found with ID %d and user_id %d", listID, userID)
	}

	return nil
}

func VerifyUserListOwnership(ctx context.Context, userID int, listID int) (bool, error) {
	query := `SELECT 1 FROM lists WHERE user_id = $1 AND list_id = $2`

	var exists int
	err := pool.QueryRow(ctx, query, userID, listID).Scan(&exists)
	if err == pgx.ErrNoRows {
		return false, fmt.Errorf("list does not belong to user: %w", err)
	}

	if exists <= 0 {
		return false, fmt.Errorf("Unknown error")
	}

	return true, nil
}
