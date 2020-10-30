package web

import (
	"encoding/json"
	"net/http"
)

func HandleInbox(w http.ResponseWriter, r *http.Request) {
	//decoder := json.NewDecoder(r.Body)




	nodeinfo, err := json.Marshal(&wellknownNodeinfo)
	if err != nil {
		logger.Warningf("Could not marshal JSON: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/activity+json")
	_, err = w.Write(nodeinfo)
	if err != nil {
		logger.Warningf("Could not write response: %s", err.Error())
		return
	}
}
