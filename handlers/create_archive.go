package handlers

import (
	"fmt"
	"gitlab.com/nod/teyitlink-web/database"
	"gitlab.com/nod/teyitlink-web/utils"
	"log"
	"net/http"
)

func CreateArchive(w http.ResponseWriter, r *http.Request) {
	requestUrl := r.PostFormValue("request_url")

	archive, err := database.CreateArchive(requestUrl)
	if err != nil {
		log.Println("Error", archive, archive.ArchiveID.String())
		return
	}

	go func() {
		utils.RunArchiveLambda(archive)
	}()

	redirectTo := fmt.Sprintf("/%s", archive.Slug)
	http.Redirect(w, r, redirectTo, http.StatusFound)
}

func CreateArchiveJson(w http.ResponseWriter, r *http.Request) {

}

func CreateArchiveLegacy(w http.ResponseWriter, r *http.Request) {

}
