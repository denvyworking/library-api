package api

import "net/http"

func (api *api) RightAuth(w http.ResponseWriter, r *http.Request) bool {
	str := r.Header.Get("Authorization")
	if str != api.authToken {
		api.logger.Error("Invalid Auth Token", "have: ", str, "need: ", api.authToken)
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return false
	}
	return true

}
