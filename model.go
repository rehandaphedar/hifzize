package main

import "git.sr.ht/~rehandaphedar/genanki-go-utils/v4/pkg/qul"

type TemplateData struct {
	Previous  Page
	Current   Page
	Next      Page
	Instances []qul.Instance
}

type Page struct {
	qul.VersePosition
	Type   PageType
	Number int
	Path   string
}

type PageType int

const (
	PageTypeNormal PageType = iota
	PageTypeSpecial
	PageTypeOpening
	PageTypeConclusion
)

type MediaEntry struct {
	Src string `yaml:"src"`
	As  string `yaml:"as"`
}
