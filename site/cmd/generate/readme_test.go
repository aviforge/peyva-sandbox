package main

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"site/content"
)

// fixtureChapters is a deliberately tiny stand-in for content.All so the
// block-rewriting tests assert on formatting rather than on the real
// book's 21 titles, which change as chapters are edited.
func fixtureChapters() []content.ChapterContent {
	return []content.ChapterContent{
		{Number: 0, Title: "First Chapter"},
		{Number: 1, Title: "Second Chapter"},
	}
}

func writeFixtureREADME(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "README.md")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("writing fixture README: %v", err)
	}
	return path
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	return string(data)
}

func TestWriteREADMETOCReplacesBlockAndPreservesSurroundingProse(t *testing.T) {
	path := writeFixtureREADME(t, `# Peyva Sandbox

## Contents

<!-- toc:start -->
- **99.** Stale Title Nobody Wants
<!-- toc:end -->

## Prerequisites

- Go 1.21+
`)

	if err := writeREADMETOC(path, fixtureChapters()); err != nil {
		t.Fatalf("writeREADMETOC: %v", err)
	}

	want := `# Peyva Sandbox

## Contents

<!-- toc:start -->
- **0.** First Chapter
- **1.** Second Chapter
<!-- toc:end -->

## Prerequisites

- Go 1.21+
`
	if got := readFile(t, path); got != want {
		t.Errorf("README mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestWriteREADMETOCErrorsOnMissingMarkersAndLeavesFileUnchanged(t *testing.T) {
	original := `# Peyva Sandbox

No markers anywhere in this file.
`
	path := writeFixtureREADME(t, original)

	err := writeREADMETOC(path, fixtureChapters())
	if err == nil {
		t.Fatal("expected an error when the toc markers are absent, got nil")
	}
	if got := readFile(t, path); got != original {
		t.Errorf("README was modified despite the error\n--- got ---\n%s", got)
	}
}

func TestWriteREADMETOCErrorsWhenMarkersAreReversed(t *testing.T) {
	original := `# Peyva Sandbox

<!-- toc:end -->
- **0.** Out of order
<!-- toc:start -->
`
	path := writeFixtureREADME(t, original)

	err := writeREADMETOC(path, fixtureChapters())
	if err == nil {
		t.Fatal("expected an error when the end marker precedes the start marker, got nil")
	}
	if got := readFile(t, path); got != original {
		t.Errorf("README was modified despite the error\n--- got ---\n%s", got)
	}
}

func TestWriteREADMETOCIsIdempotent(t *testing.T) {
	path := writeFixtureREADME(t, `# Peyva Sandbox

<!-- toc:start -->
<!-- toc:end -->
`)

	if err := writeREADMETOC(path, fixtureChapters()); err != nil {
		t.Fatalf("first writeREADMETOC: %v", err)
	}
	first := readFile(t, path)

	if err := writeREADMETOC(path, fixtureChapters()); err != nil {
		t.Fatalf("second writeREADMETOC: %v", err)
	}
	if second := readFile(t, path); second != first {
		t.Errorf("second run changed the file\n--- first ---\n%s\n--- second ---\n%s", first, second)
	}
}

func TestWriteREADMETOCListsEveryRegistryChapterInReadingOrder(t *testing.T) {
	path := writeFixtureREADME(t, `<!-- toc:start -->
<!-- toc:end -->
`)

	if err := writeREADMETOC(path, content.All); err != nil {
		t.Fatalf("writeREADMETOC: %v", err)
	}

	want := "<!-- toc:start -->\n"
	for _, c := range content.All {
		want += "- **" + strconv.Itoa(c.Number) + ".** " + c.Title + "\n"
	}
	want += "<!-- toc:end -->\n"

	if got := readFile(t, path); got != want {
		t.Errorf("registry TOC mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestWriteREADMETOCErrorsWhenFileMissing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "does-not-exist.md")

	if err := writeREADMETOC(path, fixtureChapters()); err == nil {
		t.Fatal("expected an error for a missing README, got nil")
	}
}
