package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/WilliamCWhite/tallykeeper_backend/db"
)

func ListsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	if userID == 0 {
		fmt.Println("Error retrieivng userID from context in ListsHandler")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var err error

	switch r.Method {
		case http.MethodGet:
			err = ListGet(w, r, userID)
		case http.MethodPost:
			err = ListPost(w, r, userID)
	}

	if err != nil {
		fmt.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Failed encoding JSON", http.StatusInternalServerError)
		return
	}
}

func ListGet(w http.ResponseWriter, r *http.Request, userID int) error {
	lists, err := db.GetListsByUserID(r.Context(), userID)
	if err != nil {
		fmt.Printf("Error from ListGet: %v", err)
		http.Error(w, "Internal datbase request error", http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(lists)
	if err != nil {
		fmt.Printf("Error encoding from ListGet: %v", err)
		http.Error(w, "Internal datbase request error", http.StatusInternalServerError)
		return err
	}
	return nil
}

func ListPost(w http.ResponseWriter, r *http.Request, userID int) error {
	var list db.List
	
	err := json.NewDecoder(r.Body).Decode(&list)
	if err != nil {
		fmt.Printf("error decoding json in listpost: %v", err)
		http.Error(w, "bad json in request", http.StatusBadRequest)
		return err
	}

	list.UserID = userID
	
	listID, err := db.CreateList(r.Context(), list)
	if err != nil {
		fmt.Printf("error creating list in db: %v", err)
		http.Error(w, "Internal datbase request error", http.StatusInternalServerError)
		return err
	}

	list.ListID = listID

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(list)
	if err != nil {
		fmt.Printf("Error encoding from ListPost: %v", err)
		http.Error(w, "Internal datbase request error", http.StatusInternalServerError)
		return err
	}
	return nil
}

func ListPut(w http.ResponseWriter, r *http.Request, userID int) error {
	var list db.List
	err := json.NewDecoder(r.Body).Decode(&list)
	if err != nil {
		fmt.Printf("error decoding json in listput: %v", err)
		http.Error(w, "bad json in request", http.StatusBadRequest)
		return err
	}
	list.UserID = userID

	err = db.UpdateList(r.Context(), list)
	if err != nil {
		fmt.Printf("error updating list in db: %v", err)
		http.Error(w, "Internal database request error", http.StatusInternalServerError)
		return err
	}

	list.TimeModified = time.Now()

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(list)
	if err != nil {
		fmt.Printf("Error encoding from ListPut: %v", err)
		http.Error(w, "Internal datbase request error", http.StatusInternalServerError)
		return err
	}
	return nil
}

func ListDelete(w http.ResponseWriter, r *http.Request, userID int) error {
	var list db.List
	err := json.NewDecoder(r.Body).Decode(&list)
	if err != nil {
		fmt.Printf("error decoding json in listdelete: %v", err)
		http.Error(w, "bad json in request", http.StatusBadRequest)
		return err
	}

	listID := list.ListID

	err = db.DeleteList(r.Context(), listID, userID)
	if err != nil {
		fmt.Printf("error deleting list in db: %v", err)
		http.Error(w, "Internal database request error", http.StatusInternalServerError)
		return err
	}

	return nil
}
