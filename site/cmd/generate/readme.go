package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"site/content"
)

const (
	tocStartMarker = "<!-- toc:start -->"
	tocEndMarker   = "<!-- toc:end -->"
)

// writeREADMETOC rewrites the chapter list between the toc markers in the
// file at readmePath, leaving every byte outside them untouched.
//
// The markers are required rather than optional: a README missing them is a
// mistake worth failing on, since the alternative is guessing where the list
// belongs and silently appending a second one.
//
// The list is deliberately unordered (`- **N.**`) rather than a Markdown
// ordered list. Markdown renumbers ordered lists from the first item, so
// typed digits are discarded on render — the displayed numbers would drift
// from ChapterContent.Number the moment a chapter was inserted or reordered.
//
// The write is skipped when the rendered result already matches what is on
// disk, so generating the site twice does not dirty the working tree.
func writeREADMETOC(readmePath string, chapters []content.ChapterContent) error {
	data, err := os.ReadFile(readmePath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", readmePath, err)
	}
	existing := string(data)

	start := strings.Index(existing, tocStartMarker)
	if start < 0 {
		return fmt.Errorf("%s: missing %s marker", readmePath, tocStartMarker)
	}
	end := strings.Index(existing, tocEndMarker)
	if end < 0 {
		return fmt.Errorf("%s: missing %s marker", readmePath, tocEndMarker)
	}
	if end < start {
		return fmt.Errorf("%s: %s appears before %s", readmePath, tocEndMarker, tocStartMarker)
	}

	var block strings.Builder
	block.WriteString(tocStartMarker)
	block.WriteString("\n")
	for _, c := range chapters {
		block.WriteString("- **")
		block.WriteString(strconv.Itoa(c.Number))
		block.WriteString(".** ")
		block.WriteString(c.Title)
		block.WriteString("\n")
	}
	block.WriteString(tocEndMarker)

	updated := existing[:start] + block.String() + existing[end+len(tocEndMarker):]
	if updated == existing {
		return nil
	}
	return os.WriteFile(readmePath, []byte(updated), 0o644)
}
