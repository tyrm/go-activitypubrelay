package models

type PublicKey struct {
	ID           string `json:"id,omitempty"`
	Owner        string `json:"owner,omitempty"`
	PublicKeyPem string `json:"publicKeyPem,omitempty"`
}

type Endpoints struct {
	SharedInbox string `json:"sharedInbox,omitempty"`
}

type Actor struct {
	Context           interface{} `json:"@context,omitempty"`
	Endpoints         Endpoints   `json:"endpoints,omitempty"`
	Followers         string      `json:"followers,omitempty"`
	Following         string      `json:"following,omitempty"`
	Inbox             string      `json:"inbox,omitempty"`
	Name              string      `json:"name,omitempty"`
	Type              string      `json:"type,omitempty"`
	ID                string      `json:"id,omitempty"`
	PublicKey         PublicKey   `json:"publicKey,omitempty"`
	Summary           string      `json:"summary,omitempty"`
	PreferredUsername string      `json:"preferredUsername,omitempty"`
	URL               string      `json:"url,omitempty"`
}
