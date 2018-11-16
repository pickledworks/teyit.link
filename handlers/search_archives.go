package handlers

import (
	"github.com/araddon/dateparse"
	"gitlab.com/nod/teyit/link/database"
	"gitlab.com/nod/teyit/link/utils"
	"net/http"
	"net/url"
	"strconv"
)

func SearchArchives(w http.ResponseWriter, r *http.Request) {
	results, _ := processSearchArchives(r)
	RespondSuccessTemplate(w, r, "search", results)
}

func SearchArchivesJson(w http.ResponseWriter, r *http.Request) {
	results, _ := processSearchArchives(r)
	RespondSuccessJson(w, results)
}

func processSearchArchives(r *http.Request) (database.ArchiveSearchResults, error) {
	query := r.FormValue("q")
	offset, _ := strconv.Atoi(r.FormValue("offset"))
	after, _ := dateparse.ParseAny(r.FormValue("after"))
	before, _ := dateparse.ParseAny(r.FormValue("before"))

	// sets default value for limit if nil
	limit, _ := strconv.Atoi(r.FormValue("limit"))
	if limit == 0 {
		limit = utils.GetConfig().SearchLimit
	}

	searchParams := database.ArchiveSearchParams{
		Query:  query,
		After:  after,
		Before: before,
		Limit:  limit,
		Offset: offset,
	}

	archives, total, e := database.FindArchives(searchParams)

	nextPageUrl := url.URL{}
	// sets nextPageURL parameters only if present
	if offset+limit < total {
		nextPageQueryParams := r.URL.Query()
		nextPageQueryParams.Set("offset", strconv.Itoa(offset+limit))
		nextPageQueryParams.Set("limit", strconv.Itoa(limit))
		nextPageUrl.RawQuery = nextPageQueryParams.Encode()
		nextPageUrl.Path = r.URL.Path
	}

	previousPageUrl := url.URL{}
	// sets previousPageUrl parameters only if present
	if offset != 0 {
		previousPageQueryParams := r.URL.Query()
		previousPageQueryParams.Set("offset", strconv.Itoa(offset-limit))
		previousPageQueryParams.Set("limit", strconv.Itoa(limit))
		previousPageUrl.RawQuery = previousPageQueryParams.Encode()
		previousPageUrl.Path = r.URL.Path
	}

	return database.ArchiveSearchResults{
		Results:         archives,
		Total:           total,
		PreviousPageUrl: previousPageUrl.String(),
		NextPageUrl:     nextPageUrl.String()}, e
}
