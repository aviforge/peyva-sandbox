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

// LandingData is what the entry point renders: the site in one screen, and a
// way into chapter 0. It is not a chapter, so it carries none of a chapter's
// machinery.
type LandingData struct {
	Chapters      []content.ChapterContent
	Languages     []content.Language
	AssetPrefix   string
	StyleVersion  string
	ScriptVersion string
}

type PageData struct {
	Chapter       content.ChapterContent
	HeroAvailable bool
	Roadmap       []roadmapView
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

	// ThinkingLine is the preamble for turns that produce an answer rather than
	// a change.
	ThinkingLine template.HTML

	// Setup is rendered on the chapter that starts the project, with a copy
	// button per file, so a reader never leaves the site to fetch one.
	Setup []content.SetupFile

	Languages []content.Language

	// Prompts are the chapter's turns, in order, with {os} expanded.
	Prompts []promptView

	// AsidePrompts and AsideHeroAvailable are the sidebar's, when the chapter
	// carries one. Nil and false otherwise, and the template renders nothing.
	AsidePrompts       []promptView
	AsideHeroAvailable bool

	// HasPortalPrompt and HasThinkingPrompt say which sets of standing rules
	// this page needs to show. A chapter with no portal work has no use for the
	// portal's rules sitting above its prompts.
	HasBuildPrompt    bool
	HasPortalPrompt   bool
	HasThinkingPrompt bool

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
	SystemPickerNumber int
	SystemPickerTitle  string

	// ChapterTokens is what all of a chapter's prompts cost to send.
	ChapterTokens int

	// RunnerScripts is every operating system's script, rendered on the one
	// chapter that hands it over. All of them are in the page and the script
	// shows the selected one, so switching needs no reload and no fetch.
	RunnerScripts []content.RunnerScript

	// LanguageIsChosenHere is true only on the chapter that owns the picker.
	// The choice applies to every chapter, so offering it on all of them
	// invites a switch halfway through, and a reader who moves from Python to
	// Go at chapter 12 has eleven chapters of code that no longer fits.
	LanguageIsChosenHere bool

	// LanguageName is baked in for the chapters that only display the choice.
	LanguageName string
}

// promptView is one prompt as the page renders it.
type promptView struct {
	Label    string
	Text     template.HTML
	Portal   bool
	Thinking bool
	Step     int
	Steps    int
	// Tokens is what the prompt costs to send, which is the only part of a
	// chapter's cost this site can know.
	Tokens int
}

// Numbered reports whether the turn should show its position. A chapter with a
// single turn has nothing to count.
func (p promptView) Numbered() bool { return p.Steps > 1 }
