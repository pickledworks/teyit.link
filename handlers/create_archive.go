package handlers

import (
	"fmt"
	"gitlab.com/nod/teyit/link/database"
	"gitlab.com/nod/teyit/link/utils"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func CreateArchive(w http.ResponseWriter, r *http.Request) {
	archive, fresh, err := createArchiveAction(r)

	if err != nil {
		log.Println("Error while creating archive", err)
		RespondInternalServerError(w, err)
	}

	config := utils.GetConfig()
	returnUrl := strings.Builder{}
	returnUrl.WriteString(fmt.Sprintf("%s/%s", config.BaseUrl, archive.Slug))

	if fresh == true {
		returnUrl.WriteString("?fresh=false")
	}

	http.Redirect(w, r, returnUrl.String(), http.StatusFound)
}

func CreateArchiveJson(w http.ResponseWriter, r *http.Request) {
	archive, fresh, err := createArchiveAction(r)

	if err != nil {
		log.Println("Error while creating archive", err)
		RespondInternalServerErrorJson(w, err)
	}

	if fresh == true {
		w.WriteHeader(http.StatusCreated)
	}

	RespondJson(w, archive.GetAsPublic())
}

func createArchiveAction(r *http.Request) (*database.Archive, bool, error) {
	requestUrl := r.FormValue("request_url")
	force, err := strconv.ParseBool(r.FormValue("force"))

	if force != true {
		previous, err := database.CountArchivesByRequestUrl(requestUrl)

		if err != nil {
			return nil, false, err
		}

		if previous.Count > 0 {
			return &previous.LastArchive, false, nil
		}
	}

	archive, err := database.CreateArchive(requestUrl)
	if err != nil {
		log.Println("Error while creating archive", err, requestUrl, &archive)
		return nil, false, err
	}

	return archive, true, nil
}
