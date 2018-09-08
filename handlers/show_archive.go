package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gitlab.com/nod/teyit/link/database"
	"gitlab.com/nod/teyit/link/utils"
	"html/template"
	"net/http"
)

type ShowArchiveTemplateVariables struct {
	Archive                  *database.Archive
	ShowAlreadyArchivedModal bool
	ArchiveData              template.JS
}

func ShowArchive(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	archive, err := database.GetArchive(slug)

	if err != nil {
		NotFoundPage(w, r)
	} else {
		showAlreadyArchivedModal := false
		if r.FormValue("fresh") == "false" {
			showAlreadyArchivedModal = true
		}

		archiveData, _ := json.Marshal(archive.GetAsPublic())

		data := ShowArchiveTemplateVariables{
			Archive:                  archive,
			ShowAlreadyArchivedModal: showAlreadyArchivedModal,
			ArchiveData:              template.JS(archiveData),
		}

		RespondSuccessTemplate(w, r, "archive_show", data)
	}
}

func ShowArchiveJson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	archive, err := database.GetArchive(slug)
	if err != nil {
		NotFoundPage(w, r)
	} else {
		RespondSuccessJson(w, archive.GetAsPublic())
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
