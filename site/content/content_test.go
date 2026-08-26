package content

import (
	"strings"
	"testing"
)

func TestRoadmapHasTwentyOneChapters(t *testing.T) {
	if len(Roadmap) != 21 {
		t.Errorf("len(Roadmap) = %d, want 21", len(Roadmap))
	}
}

func TestAllHasTwentyOneChapters(t *testing.T) {
	if len(All) != 21 {
		t.Errorf("len(All) = %d, want 21", len(All))
	}
}

func TestEveryChapterHasRequiredFields(t *testing.T) {
	seenNumbers := map[int]bool{}
	seenSlugs := map[string]bool{}

	for _, c := range All {
		if c.Title == "" {
			t.Errorf("chapter %d: Title is empty", c.Number)
		}
		if c.Slug == "" {
			t.Errorf("chapter %d: Slug is empty", c.Number)
		}
		if c.HeroImage == "" {
			t.Errorf("chapter %d: HeroImage is empty", c.Number)
		}
		if len(c.Intuition) == 0 {
			t.Errorf("chapter %d: Intuition is empty", c.Number)
		}
		if len(c.Concepts) == 0 {
			t.Errorf("chapter %d: expected at least one concept", c.Number)
		}
		if len(c.UnderTheHood) == 0 {
			t.Errorf("chapter %d: expected at least one under-the-hood point", c.Number)
		}
		if c.BuildIt.Technique == "" {
			t.Errorf("chapter %d: BuildIt.Technique is empty", c.Number)
		}
		if c.BuildIt.Why == "" {
			t.Errorf("chapter %d: BuildIt.Why is empty", c.Number)
		}
		if !strings.Contains(c.BuildIt.Prompt, "Done when") {
			t.Errorf("chapter %d: BuildIt.Prompt is missing a 'Done when' line", c.Number)
		}
		// The prompt is copied out of the page on its own, so it must not
		// depend on anything only visible here.
		for _, dangling := range []string{"steps above", "the steps", "above in"} {
			if strings.Contains(c.BuildIt.Prompt, dangling) {
				t.Errorf("chapter %d: BuildIt.Prompt refers to off-page content (%q) — it must stand alone", c.Number, dangling)
			}
		}
		if len(c.BreakIt.Exercises) == 0 {
			t.Errorf("chapter %d: expected at least one break-it exercise", c.Number)
		}

		if seenNumbers[c.Number] {
			t.Errorf("duplicate chapter Number: %d", c.Number)
		}
		seenNumbers[c.Number] = true

		if seenSlugs[c.Slug] {
			t.Errorf("duplicate chapter Slug: %q", c.Slug)
		}
		seenSlugs[c.Slug] = true
	}
}

// A technique nobody else has named is a technique we made up, and a made-up
// name is worse than no name: the reader carries it into rooms where it means
// nothing. Every technique must cite a corpus the reader can go and read.
func TestEveryTechniqueCitesARecognisedSource(t *testing.T) {
	for _, c := range All {
		if c.BuildIt.Source == "" {
			t.Errorf("chapter %d: BuildIt.Source is empty — technique %q must cite where it is documented",
				c.Number, c.BuildIt.Technique)
			continue
		}
		recognised := false
		for _, corpus := range RecognisedSources {
			if strings.HasPrefix(c.BuildIt.Source, corpus) {
				recognised = true
				break
			}
		}
		if !recognised {
			t.Errorf("chapter %d: BuildIt.Source %q cites no recognised corpus (want one of %v)",
				c.Number, c.BuildIt.Source, RecognisedSources)
		}
	}
}

// Teaching a range of prompting techniques is a goal of the book, so a
// repeated technique means a chapter is missing its lesson.
func TestEveryChapterTeachesADistinctTechnique(t *testing.T) {
	seen := map[string]int{}
	for _, c := range All {
		if prev, ok := seen[c.BuildIt.Technique]; ok {
			t.Errorf("chapters %d and %d both use technique %q", prev, c.Number, c.BuildIt.Technique)
		}
		seen[c.BuildIt.Technique] = c.Number
	}
}

func TestChapterNumbersMatchRoadmap(t *testing.T) {
	roadmapTitles := map[int]string{}
	for _, r := range Roadmap {
		roadmapTitles[r.Number] = r.Title
	}

	for _, c := range All {
		want, ok := roadmapTitles[c.Number]
		if !ok {
			t.Errorf("chapter %d (%q) has no matching Roadmap entry", c.Number, c.Title)
			continue
		}
		if c.Title != want {
			t.Errorf("chapter %d: Title %q does not match Roadmap title %q", c.Number, c.Title, want)
		}
	}
}

// componentChapters maps each named component of the system to the chapter
// that introduces it. The names were chosen deliberately so the reader builds
// one architecture across the book rather than twenty-one unrelated exercises.
var componentChapters = map[string]int{
	"Vault":      0,
	"Gateway":    2,
	"Teller":     4,
	"Ledger":     7,
	"Courier":    12,
	"Reconciler": 20,
}

// A component named in Build It but absent from Concepts is a word the reader
// is asked to build without ever being taught. Concepts is where this book
// defines its vocabulary, so that is where a component has to be introduced —
// not only inside the prompt, which is written to be copied out of the page.
func TestEveryComponentIsDefinedInTheConceptsOfItsChapter(t *testing.T) {
	byNumber := map[int]ChapterContent{}
	for _, c := range All {
		byNumber[c.Number] = c
	}

	for component, number := range componentChapters {
		chapter, ok := byNumber[number]
		if !ok {
			t.Errorf("component %q claims chapter %d, which does not exist", component, number)
			continue
		}
		defined := false
		for _, concept := range chapter.Concepts {
			if concept.Term == component {
				defined = true
				if concept.Description == "" {
					t.Errorf("chapter %d: Concept %q has an empty description", number, component)
				}
				break
			}
		}
		if !defined {
			t.Errorf("chapter %d introduces %q in Build It but never defines it in Concepts",
				number, component)
		}
	}
}

// The component belongs to the chapter that builds it. Defining it earlier
// spends the name before the reader has anything to attach it to.
func TestNoChapterNamesAComponentBeforeItsOwnChapter(t *testing.T) {
	for _, c := range All {
		for _, concept := range c.Concepts {
			if intro, ok := componentChapters[concept.Term]; ok && intro != c.Number {
				t.Errorf("chapter %d defines %q, but %q is introduced in chapter %d",
					c.Number, concept.Term, concept.Term, intro)
			}
		}
	}
}
