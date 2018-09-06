package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gitlab.com/nod/teyit/link/views"
	"net/http"
)

func CreateRoutes() *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	r.HandleFunc("/", ShowHomepage).Methods("GET")
	r.HandleFunc("/search", SearchArchives).Methods("GET")
	r.HandleFunc("/api/search", SearchArchivesJson).Methods("GET")

	r.HandleFunc("/archive", CreateArchive).Methods("POST")
	r.HandleFunc("/api/archive", CreateArchiveJson).Methods("POST")
	r.HandleFunc("/api/count-previous-archives", CheckPreviousArchives).Methods("GET")
	r.HandleFunc("/{slug}", ShowArchive).Methods("GET")
	r.HandleFunc("/{slug}", ShowArchiveJson).Methods("GET").HeadersRegexp("Content-Type", "application/json")
	r.HandleFunc("/api/archives/{slug}", ShowArchiveJson).Methods("GET")

	// below are legacy links from v1, we plan to phase these out
	// but we can't immediately because we suspect programmatic usage
	r.HandleFunc("/new", CreateArchiveLegacy).Methods("POST", "GET")
	r.HandleFunc("/bookmark", CreateArchiveLegacy).Methods("POST", "GET")
	r.HandleFunc("/add", CreateArchiveLegacy).Methods("POST", "GET")

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

func RespondInvalidRequestJson(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func RespondSuccessTemplate(w http.ResponseWriter, r *http.Request, page string, data interface{}) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	view := views.NewView("default", page)
	view.Render(w, r, data)
}
