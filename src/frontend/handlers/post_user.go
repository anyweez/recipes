package handlers

import (
	gproto "code.google.com/p/goprotobuf/proto"
	"encoding/json"
//	"lib/fetch"
	"net/http"
	proto "proto"
)

func PostUser(w http.ResponseWriter, r *http.Request) {
	// Get the current user
	// Encode as JSON
	
	
	user := proto.User{
		Name: gproto.String("luke"),
		Id: gproto.Uint64(1),
	}
	
	data, _ := json.Marshal(user)
	w.Write( data )
}

