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
}
