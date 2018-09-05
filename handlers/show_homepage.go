package handlers

import (
	"net/http"
)

func ShowHomepage(w http.ResponseWriter, r *http.Request) {
	RespondSuccessTemplate(w, "homepage", nil)
}
