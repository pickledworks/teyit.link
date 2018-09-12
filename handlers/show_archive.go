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

func ShowArchiveHandler(w http.ResponseWriter, r *http.Request) {
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

func ShowArchiveApiHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	archive, err := database.GetArchive(slug)
	if err != nil {
		NotFoundPage(w, r)
	} else {
		RespondSuccessJson(w, archive.GetAsPublic())
	}
}

// Instead of hard linking the AWS S3 resource URLs, we define a redirection with a presigned links
// These links have close expire dates, discouraging use of direct links in API integrations
// This way if we change storage providers, API integrations would still work fine

// This is for the snapshot, which is the inlined HTML of the archive
func redirectToArchiveResource(file string, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]
	var download bool
	if r.FormValue("dl") == "1" {
		download = true
	}

	archive, err := database.GetArchive(slug)
	if err != nil {
		NotFoundPage(w, r)
	} else {
		url := utils.PresignArchiveResource(&utils.ArchiveResourceRequest{
			ArchiveSlug: archive.Slug,
			ArchiveID: archive.ArchiveID.String(),
			File: file,
			Download: download,
		})
		http.Redirect(w, r, url, http.StatusFound)
	}
}

// This is for the screenshot
func ShowArchiveScreenshotHandler(w http.ResponseWriter, r *http.Request) {
	redirectToArchiveResource("screenshot.png", w, r)
}

// This is for the snapshot, which is the inlined HTML of the archive
func ShowArchiveSnapshotHandler(w http.ResponseWriter, r *http.Request) {
	redirectToArchiveResource("index.html", w, r)
}
