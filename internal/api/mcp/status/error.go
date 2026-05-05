package status

import "net/http"

func Unauthorized(w http.ResponseWriter, err error) {
	w.Header().Set("WWW-Authenticate", `Bearer`)

	http.Error(w, err.Error(), http.StatusUnauthorized)
}
