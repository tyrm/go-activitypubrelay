package models

type Link struct {
	HRef string `json:"href,omitempty"`
	Rel  string `json:"rel,omitempty"`
	Type string `json:"type,omitempty"`
}

type WebFinger struct {
	Aliases []string `json:"aliases,omitempty"`
	Links   []Link   `json:"links,omitempty"`
	Subject string   `json:"subject,omitempty"`
}
