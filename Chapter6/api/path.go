package main

import "strings"

// PathSeparator ...
const PathSeparator = "/"

// Path is a struct to store Path attributes
type Path struct {
	Path string
	ID   string
}

// NewPath creates a new Path
func NewPath(p string) *Path {
	var id string
	p = strings.Trim(p, PathSeparator)
	s := strings.Split(p, PathSeparator)
	if len(s) > 1 {
		id = s[len(s)-1]
		p = strings.Join(s[:len(s)-1], PathSeparator)
	}
	return &Path{Path: p, ID: id}
}

// HasID checks whether ID has value set
func (p *Path) HasID() bool {
	return len(p.ID) > 0
}
