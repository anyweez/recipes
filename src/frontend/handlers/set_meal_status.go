package handlers

import (
	"encoding/json"
	fee "frontend/errors"
	log "logging"
	"net/http"
)

func SetMealStatus(w http.ResponseWriter, r *http.Request, le log.LogEvent) {
	err := fee.HANDLER_NOT_IMPLEMENTED
	data, _ := json.Marshal(err)

	w.WriteHeader(err.HttpCode)
	w.Write(data)
}