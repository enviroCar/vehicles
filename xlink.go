package main

import (
	"encoding/json"
	"net/url"
)

// Link is a link.
type Link struct {
	URL       *url.URL
	MediaType string
	Title     string
	Relation  string
}

// NewLink creates a new link.
func NewLink(url *url.URL, relation, mediaType, title string) *Link {
	return &Link{
		URL:       url,
		MediaType: mediaType,
		Title:     title,
		Relation:  relation,
	}
}

type jsonLink struct {
	URL       string `json:"href,omitempty"`
	MediaType string `json:"type,omitempty"`
	Title     string `json:"title,omitempty"`
	Relation  string `json:"rel,omitempty"`
}

// UnmarshalJSON is required by json.Unmarshaler
func (l *Link) UnmarshalJSON(bytes []byte) error {
	var s *jsonLink
	err := json.Unmarshal(bytes, s)
	if err != nil {
		return err
	}
	return l.fromJSON(s)
}

func (l *Link) fromJSON(s *jsonLink) (err error) {
	l.URL, err = url.Parse(s.URL)
	l.MediaType = s.MediaType
	l.Title = s.Title
	l.Relation = s.Relation
	return
}

// MarshalJSON is required by json.Marshaler
func (l *Link) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.toJSON())
}

func (l *Link) toJSON() *jsonLink {
	return &jsonLink{
		URL:       l.URL.String(),
		MediaType: l.MediaType,
		Title:     l.Title,
		Relation:  l.Relation,
	}
}

// Linked is a linked object.
type Linked struct {
	Links []*Link `json:"links,omitempty"`
}

// AddLink adds a link to the Linked object
func (l *Linked) AddLink(link *Link) {
	if link != nil {
		l.Links = append(l.Links, link)
	}
}

var (
	_ json.Marshaler   = (*Link)(nil)
	_ json.Unmarshaler = (*Link)(nil)
)
