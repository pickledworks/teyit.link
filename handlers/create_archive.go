package handlers

import (
	"fmt"
	"gitlab.com/nod/teyit/link/database"
	"log"
	"net/http"
)

type CreateArchiveResponse struct {

}
func CheckPreviousArchives(w http.ResponseWriter, r *http.Request) {
	requestUrl := r.FormValue("request_url")

	resp, err := database.CountArchivesByRequestUrl(requestUrl)
	if err != nil {
		log.Println("Error in counting archives", err)
		return
	}

	RespondSuccessJson(w, resp)
}

func CreateArchive(w http.ResponseWriter, r *http.Request) {
	requestUrl := r.PostFormValue("request_url")

	archive, err := database.CreateArchive(requestUrl)
	if err != nil {
		log.Println("Error", archive, archive.ArchiveID.String())
		return
	}

	redirectTo := fmt.Sprintf("/%s", archive.Slug)
	http.Redirect(w, r, redirectTo, http.StatusFound)
}

func CreateArchiveJson(w http.ResponseWriter, r *http.Request) {
	requestUrl := r.PostFormValue("request_url")

	archive, err := database.CreateArchive(requestUrl)
	if err != nil {
		log.Println("Error", err, requestUrl, &archive)

		RespondInvalidRequestJson(w, err)
		return
	}

	RespondSuccessJson(w, database.GetArchiveAsArchivePublic(archive))
}

func CreateArchiveLegacy(w http.ResponseWriter, r *http.Request) {

}

