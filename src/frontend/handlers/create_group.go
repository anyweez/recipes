package handlers

import (
	"encoding/json"
	fee "frontend/errors"
	"net/http"
)

func CreateGroup(w http.ResponseWriter, r *http.Request) {
	err := fee.HANDLER_NOT_IMPLEMENTED
	data, _ := json.Marshal(err)

	w.WriteHeader(err.HttpCode)
	w.Write(data)
}

