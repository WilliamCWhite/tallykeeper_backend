package handlers

import(
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/WilliamCWhite/tallykeeper_backend/db"
)

func ListsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	if userID == 0 {
		fmt.Println("Error retrieivng userID from context in ListsHandler")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	lists, err := db.GetListsByUserID(r.Context(), userID)
	if err != nil {
		fmt.Printf("Error from GetListsByUserID: %v", err)
		http.Error(w, "Internal database request error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(lists)
	if err != nil {
		fmt.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Failed encoding JSON", http.StatusInternalServerError)
		return
	}
}


