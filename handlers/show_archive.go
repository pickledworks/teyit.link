package handlers

import (
	"github.com/gorilla/mux"
	"gitlab.com/nod/teyit/link/database"
	"net/http"
)

func ShowArchive(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	archive, err := database.GetArchive(slug)
	if err != nil {
		// result in error
	}

	RespondSuccessTemplate(w, "archive_show", archive)
}

func ShowArchiveJson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

}
