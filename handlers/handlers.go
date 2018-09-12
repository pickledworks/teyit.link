package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	_ "gitlab.com/nod/teyit/link/statik"
	"gitlab.com/nod/teyit/link/utils"
	"log"
	"net/http"
)

func CreateRoutes() *mux.Router {
	r := mux.NewRouter()

	// Handle homepage inline
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		RespondSuccessTemplate(w, r, "homepage", nil)
	}).Methods("GET")

	r.HandleFunc("/search", SearchArchives).Methods("GET")
	r.HandleFunc("/api/search", SearchArchivesJson).Methods("GET")

	// below are legacy links from v1, we plan to phase these out
	// but we can't immediately because we suspect programmatic usage
	r.HandleFunc("/new", CreateArchiveHandler).Methods("POST", "GET")
	r.HandleFunc("/bookmark", CreateArchiveHandler).Methods("POST", "GET")
	r.HandleFunc("/add", CreateArchiveHandler).Methods("POST", "GET")
	// Next up, the current API handler for creating archives
	r.HandleFunc("/api/archive", CreateArchiveApiHandler).Methods("POST", "GET")

	// Display the archive
	r.HandleFunc("/{slug}", ShowArchiveHandler).Methods("GET")
	r.HandleFunc("/{slug}", ShowArchiveApiHandler).Methods("GET").HeadersRegexp("Content-Type", "application/json")
	r.HandleFunc("/api/archives/{slug}", ShowArchiveApiHandler).Methods("GET")

	r.HandleFunc("/{slug}/screenshot", ShowArchiveScreenshotHandler).Methods("GET")
	r.HandleFunc("/{slug}/snapshot", ShowArchiveSnapshotHandler).Methods("GET")

	// Handle static files
	var staticServer http.Handler

	// If we are in development, just bind the directory
	if utils.GetConfig().Env == "development" {
		staticServer = http.FileServer(http.Dir("./public"))
	} else {
		// We use statik file system in production, meaning all the assets are bundled inside the binary
		statikFS, err := fs.New()
		if err != nil {
			log.Fatal("statik fs", err)
		}
		staticServer = http.FileServer(statikFS)
	}

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticServer))

	r.NotFoundHandler = http.HandlerFunc(NotFoundPage)

	return r
}

func NotFoundPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "<h2>Sorry Could not Find Resource. 404 Error</h2>")
}

func RespondSuccessJson(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func RespondJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func RespondError(w http.ResponseWriter, r *http.Request, err error) {
	var errText string
	var statusCode int

	switch err {
	case utils.ErrorUrlInvalid, utils.ErrorUrlBlocked, utils.ErrorUrlSchemeInvalid, utils.ErrorUrlEmptyUrl:
		statusCode = http.StatusBadRequest
		errText = err.Error()
	default:
		statusCode = http.StatusInternalServerError
		errText = "internal server error"
	}

	w.WriteHeader(statusCode)

	if r.Header.Get("Content-Type") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": errText,
		})
		return
	}
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "Error occurred: %s", errText)
}

func RespondSuccessTemplate(w http.ResponseWriter, r *http.Request, page string, data interface{}) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	view := utils.NewView("default", page)
	view.Render(w, r, data)
}
