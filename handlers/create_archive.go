package handlers

import (
	"fmt"
	"gitlab.com/nod/teyit/link/database"
	"gitlab.com/nod/teyit/link/utils"
	"log"
	"net/http"
	"time"
)

func CreateArchive(w http.ResponseWriter, r *http.Request) {
	requestUrl := r.PostFormValue("request_url")

	archive, err := database.CreateArchive(requestUrl)
	if err != nil {
		log.Println("Error", archive, archive.ArchiveID.String())
		return
	}

	go func() {
		now := time.Now()
		result, err := utils.RunArchiveLambda(archive.ArchiveID, archive.RequestUrl)

		if err != nil {
			log.Println("Error", err)
			archive.FailedAt = &now
		} else {
			archive.MetaTitle = result.Title
			archive.MetaDescription = result.Description
			archive.Image = result.Image
			archive.ArchivedAt = &now
		}

		database.SaveArchive(archive)
	}()

	redirectTo := fmt.Sprintf("/%s", archive.Slug)
	http.Redirect(w, r, redirectTo, http.StatusFound)
}

func CreateArchiveJson(w http.ResponseWriter, r *http.Request) {

}

func CreateArchiveLegacy(w http.ResponseWriter, r *http.Request) {

}
