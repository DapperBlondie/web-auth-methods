package handlers

import "net/http"

// EnableCORS use for enabling cors for the client
func (conf *AppConf) EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
		return
	})
}

// LoadSession use for loading and saving user session in client
func (conf *AppConf) LoadSession(next http.Handler) http.Handler {
	return conf.ScsManager.LoadAndSave(next)
}
