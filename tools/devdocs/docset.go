package main

// Docset represents a docset retrieved from the DevDocs API.
type Docset struct {
	Name    string `json:"name"`
	Slug    string `json:"slug"`
	Release string `json:"release"`
}

func (d Docset) FullName() string {
	if d.Release == "" {
		return d.Name
	}

	return d.Name + " " + d.Release
}
