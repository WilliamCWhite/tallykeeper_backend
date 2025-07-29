package db

import (
	"fmt"
	"context"
	"time"
)

// From context and a listID, returns an array of entries from that list
func GetEntries(ctx context.Context, listID int) ([]Entry, error) {
	query := `SELECT entry_id, name, score, time_created, time_modified
		FROM entries WHERE list_id = $1
		ORDER BY time_created`

	rows, err := pool.Query(ctx, query, listID)
	if err != nil {
		return nil, fmt.Errorf("failed to query entries from db with pgx: %w", err)
	}

	defer rows.Close()

	var entries []Entry
	for rows.Next() {
		var entry Entry
		err := rows.Scan(
			&entry.EntryID,
			&entry.Name,
			&entry.Score,
			&entry.TimeCreated,
			&entry.TimeModified,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan entry row with pgx: %w", err)
		}
		entries = append(entries, entry)
	}

	if len(entries) == 0 {
		return []Entry{}, nil
	}

	return entries, nil
}

// from context and an entry containing name, score, and a list id, adds that entry to the database, returning its entry_id
func CreateEntry(ctx context.Context, entry Entry) (*Entry, error) {
	query := `INSERT INTO entries (name, score, time_created, time_modified, list_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING entry_id, time_created, time_modified`

	var newEntryID int
	var newTimeCreated time.Time // the database and backend time modify differ, so must receive the database version
	var newTimeModified time.Time
	err := pool.QueryRow(ctx, query,
		entry.Name,
		entry.Score,
		time.Now(),
		time.Now(),
		entry.ListID,
	).Scan(&newEntryID, &newTimeCreated, &newTimeModified)

	newEntry := Entry {
		Name: entry.Name,
		Score: entry.Score,
		EntryID: newEntryID,
		TimeCreated: newTimeCreated,
		TimeModified: newTimeModified,
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create entry with pgx: %w", err)
	}

	return &newEntry, nil
}

// from context and an entry containing name, score, entry_id, and list_id, updates the entry in the database
func UpdateEntry(ctx context.Context, entry Entry) (error) {
	query := `UPDATE entries SET name = $1, score = $2, time_modified = $3
		WHERE entry_id = $4 AND list_id = $5`

	commandTag, err := pool.Exec(ctx, query,
		entry.Name,
		entry.Score,
		time.Now(),
		entry.EntryID,
		entry.ListID,
	)
	if err != nil {
		return fmt.Errorf("failed to update entry with pgx: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("no entry found/updated with ID %d for list %d using pgx", entry.EntryID, entry.ListID)
	}

	return nil
}

// from context, an entry id, and a list id, deletes an entry in the database
func DeleteEntry(ctx context.Context, entryID int, listID int) error {
	query := `DELETE FROM entries WHERE entry_id = $1 AND list_id = $2`

	commandTag, err := pool.Exec(ctx, query,
		entryID,
		listID,
	)

	if err != nil {
		return fmt.Errorf("failed to delete entry with pgx: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("no entry found with ID %d and list_id %d", entryID, listID)
	}

	return nil
}
