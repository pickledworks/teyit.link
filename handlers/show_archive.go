package handlers

import (
	"github.com/gorilla/mux"
	"gitlab.com/nod/teyit/link/database"
	"gitlab.com/nod/teyit/link/utils"
	"net/http"
)

func ShowArchive(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	archive, err := database.GetArchive(slug)
	if err != nil {
		NotFoundPage(w, r)
	} else {
		RespondSuccessTemplate(w, r, "archive_show", archive)
	}
}

func ShowArchiveJson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	archive, err := database.GetArchive(slug)
	if err != nil {
		NotFoundPage(w, r)
	} else {
		RespondSuccessJson(w, database.GetArchiveAsArchivePublic(archive))
	}
}

func RedirectToArchiveScreenshot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	archive, err := database.GetArchive(slug)
	if err != nil {
		NotFoundPage(w, r)
	} else {
		url := utils.PresignArchiveResource(archive.ArchiveID, "screenshot.png")
		http.Redirect(w, r, url, http.StatusFound)
	}
}

func RedirectToArchiveSnapshot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	archive, err := database.GetArchive(slug)
	if err != nil {
		NotFoundPage(w, r)
	} else {
		url := utils.PresignArchiveResource(archive.ArchiveID, "index.html")
		http.Redirect(w, r, url, http.StatusFound)
	}
}
