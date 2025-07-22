package db

import (
	"fmt"
	"context"
	"time"
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

// from context and a list containing a title and userId, creates a list in the database, returning its list_id
func CreateList(ctx context.Context, list List) (int, error) {
	query := `INSERT INTO lists (title, time_created, time_modified, user_id)
		VALUES ($1, $2, $3, $4)
		RETURNING list_id`
	
	var newListID int
	err := pool.QueryRow(ctx, query,
		list.Title,
		time.Now(),
		time.Now(),
		list.UserID,
	).Scan(&newListID)

	if err != nil {
		return -1, fmt.Errorf("failed to create list with pgx: %w", err)
	}

	return newListID, nil
}

// from context and a list containing a title, list_id, and user_id, updates that list in the database
func UpdateList(ctx context.Context, list List) (error) {
	query := `UPDATE lists SET title = $1, time_modified = $2
		WHERE list_id = $3 AND user_id = $4`

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
