package main

type TemplateData struct {
	Previous Page
	Current  Page
	Next     Page
}

type Page struct {
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
