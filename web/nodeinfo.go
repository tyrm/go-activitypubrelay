package web

import (
	"encoding/json"
	"litepub1/models"
	"net/http"
)


type Services struct {
	Inbound  []string `json:"inbound"`
	Outbound []string `json:"outbound"`
}

type Software struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

type UsageUsers struct {
	Total int `json:"total"`
}

type Usage struct {
	LocalPosts int   `json:"localPosts"`
	Users      UsageUsers `json:"users"`
}

type Nodeinfo struct {
	Metadata         interface{} `json:"metadata,omitempty"`
	OpenRegistration bool        `json:"openRegistrations"`
	Protocols        []string    `json:"protocols"`
	Services         Services    `json:"services"`
	Software         Software    `json:"software"`
	Usage            Usage       `json:"usage"`
	Version          string      `json:"version"`
}

type WellknownNodeinfo struct {
	Links []Link
}

func HandleNodeinfo20(w http.ResponseWriter, r *http.Request) {
	nodeinto := nodeinfoTemplate //copy
	peers, err := models.ReadFollowedInstances()
	if err != nil {
		logger.Warningf("Could not get peer list from database: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var peerList []string
	for _, peer := range *peers {
		peerList = append(peerList, peer.Hostname)
	}

	metadata := make(map[string]*[]string)
	metadata["peers"] = &peerList

	nodeinto.Metadata = metadata

	actor, err := json.Marshal(&nodeinto)
	if err != nil {
		logger.Warningf("Could not marshal JSON: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/activity+json")
	_, err = w.Write(actor)
	if err != nil {
		logger.Warningf("Could not write response: %s", err.Error())
		return
	}
}

func HandleWellKnownNodeInfo(w http.ResponseWriter, r *http.Request) {
	nodeinfo, err := json.Marshal(&myWellknownNodeinfo)
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
