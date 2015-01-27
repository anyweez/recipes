package handlers

import (
	"net/http"
)

type Handler func(w http.ResponseWriter, r *http.Request) 

type Registry map[string]Handler
