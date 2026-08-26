package main

import "site/content"

type roadmapView struct {
	Number int
	Title  string
	Built  bool
}

type PageData struct {
	Chapter       content.ChapterContent
	HeroAvailable bool
	Roadmap       []roadmapView
	Labs          []content.RoadmapEntry
	AssetPrefix   string

	// LanguageLine is baked into the page so a prompt is never handed to an
	// assistant without naming a language. Script swaps it when the reader
	// chooses; it is never absent.
	LanguageLine string
	Languages    []content.Language
}
