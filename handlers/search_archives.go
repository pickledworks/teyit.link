package handlers

import (
	"github.com/araddon/dateparse"
	"gitlab.com/nod/teyit/link/database"
	"net/http"
)

func SearchArchives(w http.ResponseWriter, r *http.Request) {
	results, _ := processSearchArchives(r)
	RespondSuccessTemplate(w, r, "search", results)
}

func SearchArchivesJson(w http.ResponseWriter, r *http.Request) {
	results, _ := processSearchArchives(r)
	RespondSuccessJson(w, results)
}

func processSearchArchives(r *http.Request) ([]database.Archive, error) {
	query := r.FormValue("q")
	after, _ := dateparse.ParseAny(r.FormValue("after"))
	before, _ := dateparse.ParseAny(r.FormValue("before"))

	searchParams := database.ArchiveSearchParams{
		Query: query,
		After: after,
		Before: before,
	}

	return database.FindArchives(searchParams)
}
