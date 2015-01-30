package handlers

import (
	fee "frontend/errors"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)	
	
	// If we can't read the body, throw an error.
	if err != nil {
		e := fee.INVALID_POST_DATA
		data, _ := json.Marshal(e)
		w.WriteHeader(e.HttpCode)
		w.Write(data)
		
		return
	}
	
	email := strings.Trim(string(body), " ")
	// TODO: validate the email address.
	log.Println(email)
	session, serr := storage.Get(r, email)
	
	if val, ok := session.Values["test"].(string); ok {
		log.Println("Existing values: " + val)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        
	}
	
	// If the session couldn't be decoded, we've got to return an error.
	// This shouldn't happen unless something were to 
	if serr != nil {
		log.Println( serr.Error() )
		cserr := fee.CORRUPTED_SESSION
		data, _ := json.Marshal(cserr)
		
		w.WriteHeader(cserr.HttpCode)
		w.Write(data)
	}
	
	session.Values["test"] = "hey!"
	werr := session.Save(r, w)
	
	if werr != nil {
		log.Fatal(werr.Error())
	}
}
