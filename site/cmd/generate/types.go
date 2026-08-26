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

	// LanguageIsChosenHere is true only on the chapter that owns the picker.
	// The choice applies to every chapter, so offering it on all of them
	// invites a switch halfway through, and a reader who moves from Python to
	// Go at chapter 12 has eleven chapters of code that no longer fits.
	LanguageIsChosenHere bool

	// LanguageName is baked in for the chapters that only display the choice.
	LanguageName string
}
