package web

import (
	"net/http"
)

func HandleCatchAll(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "", http.StatusNotFound)
}
