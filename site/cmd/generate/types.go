package main

import (
	"html/template"

	"site/content"
)

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

	// StyleVersion and ScriptVersion are content digests appended to the asset
	// URLs, so a deploy is not invisible to a browser holding a cached copy.
	StyleVersion  string
	ScriptVersion string

	// LanguageLine is baked into the page so a prompt is never handed to an
	// assistant without naming a language. Script swaps it when the reader
	// chooses; it is never absent.
	LanguageLine template.HTML

	// UILine is the portal prompt preamble, rendered only where a chapter has
	// portal work.
	UILine template.HTML

	// Setup is rendered on the chapter that starts the project, with a copy
	// button per file, so a reader never leaves the site to fetch one.
	Setup []content.SetupFile

	Languages []content.Language

	// Prompt and UIPrompt are the chapter's prompts with {os} already expanded,
	// so they are HTML rather than the raw strings on the chapter.
	Prompt   template.HTML
	UIPrompt template.HTML

	// Systems and SystemName mirror Languages: the choice is offered once, in
	// setup, and every prompt that needs it reads it.
	Systems    []content.System
	SystemName string

	// SystemMatters is true on the chapters whose prompts actually name an
	// operating system. The rest never mention one, so telling those readers
	// which is in force would be noise about a choice that changes nothing.
	SystemMatters bool

	// SystemIsChosenHere is true on the first chapter that needs a system, which
	// is the one that offers the picker. Later chapters read the choice, so a
	// reader cannot get a PowerShell runner in chapter 10 and bash commands to
	// operate it in chapter 19.
	SystemIsChosenHere bool
	SystemPickerHref   string
	SystemPickerTitle  string

	// LanguageIsChosenHere is true only on the chapter that owns the picker.
	// The choice applies to every chapter, so offering it on all of them
	// invites a switch halfway through, and a reader who moves from Python to
	// Go at chapter 12 has eleven chapters of code that no longer fits.
	LanguageIsChosenHere bool

	// LanguageName is baked in for the chapters that only display the choice.
	LanguageName string
}
