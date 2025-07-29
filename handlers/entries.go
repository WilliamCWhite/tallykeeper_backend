package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/WilliamCWhite/tallykeeper_backend/db"
	"github.com/gorilla/mux"
)


func EntriesHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	if userID == 0 {
		fmt.Println("Error retrieivng userID from context in EntriesHandler")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	listIDString := vars["listID"]
	listID, err := strconv.Atoi(listIDString)
	if err != nil {
		fmt.Printf("error converting string: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	ok, err := db.VerifyUserListOwnership(r.Context(), userID, listID)
	if err != nil || !ok {
		fmt.Printf("error verifying list ownership: %v", err)
		http.Error(w, "user id and list id don't match", http.StatusBadRequest)
		return
	}

	var err2 error
	switch r.Method {
		case http.MethodGet:
			err = EntriesGet(w, r, listID)
		case http.MethodPost:
			err = EntriesPost(w, r, listID)
			err2 = db.UpdateListTimeModified(r.Context(), listID, userID)
		case http.MethodDelete:
			err = EntriesDelete(w, r, listID)
			err2 = db.UpdateListTimeModified(r.Context(), listID, userID)
		case http.MethodPut:
			err = EntriesPut(w, r, listID)
			err2 = db.UpdateListTimeModified(r.Context(), listID, userID)
	}

	if err != nil {
		fmt.Printf("Error handling method for entries: %v", err)
		http.Error(w, "Failed handling method", http.StatusInternalServerError)
		return
	}

	if err2 != nil {
		fmt.Printf("Error updating list time modified: %v", err)
		http.Error(w, "failed updating list time_modified", http.StatusInternalServerError)
		return
	}
}

func EntriesGet(w http.ResponseWriter, r *http.Request, listID int) error {
	entries, err := db.GetEntries(r.Context(), listID)
	if err != nil {
		fmt.Printf("Error from EntryGet: %v", err)
		http.Error(w, "Internal database request error", http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(entries)
	if err != nil {
		fmt.Printf("Error encoding from EntriesGet: %v", err)
		http.Error(w, "Internal datbase request error", http.StatusInternalServerError)
		return err
	}
	return nil
}

func EntriesPost(w http.ResponseWriter, r *http.Request, listID int) error {
	var entry db.Entry
	 
	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		fmt.Printf("error decoding json in entry post: %v", err)
		http.Error(w, "bad json in request", http.StatusBadRequest)
		return err
	}

	entry.ListID = listID

	newEntry, err := db.CreateEntry(r.Context(), entry)
	if err != nil {
		fmt.Printf("error creating entry in db: %v", err)
		http.Error(w, "Internal datbase request error", http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(newEntry)
	if err != nil {
		fmt.Printf("Error encoding from EntryPost: %v", err)
		http.Error(w, "Internal datbase request error", http.StatusInternalServerError)
		return err
	}
	return nil
}

func EntriesDelete(w http.ResponseWriter, r *http.Request, listID int) error {
	var entry db.Entry
	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		fmt.Printf("error decoding json in entry delete: %v", err)
		http.Error(w, "bad json in request", http.StatusBadRequest)
		return err
	}

	entryID := entry.EntryID

	err = db.DeleteEntry(r.Context(), entryID, listID)
	if err != nil {
		fmt.Printf("error deleting entry in db: %v", err)
		http.Error(w, "Internal database request error", http.StatusInternalServerError)
		return err
	}

	return nil
}

func EntriesPut(w http.ResponseWriter, r *http.Request, listID int) error {
	var entry db.Entry
	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		fmt.Printf("error decoding json in entryput: %v", err)
		http.Error(w, "bad json in request", http.StatusBadRequest)
		return err
	}
	entry.ListID = listID

	err = db.UpdateEntry(r.Context(), entry)
	if err != nil {
		fmt.Printf("error updating entry in db: %v", err)
		http.Error(w, "Internal database request error", http.StatusInternalServerError)
		return err
	}

	entry.TimeModified = time.Now()

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(entry)
	if err != nil {
		fmt.Printf("Error encoding from EntryPut: %v", err)
		http.Error(w, "Internal datbase request error", http.StatusInternalServerError)
		return err
	}
	return nil
}
