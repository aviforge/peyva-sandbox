package content

var Chapter11 = ChapterContent{
	Number:     11,
	Slug:       "chapter-11",
	Title:      "Sharing Work: Caching",
	Subtitle:   "Caching is like keeping popular items ready in advance, so we don't make (or check) everything again.",
	Category:   "System Design",
	Difficulty: "Intermediate",
	EstTime:    "15 min",
	QuickTip:   "A stale cache is worse than no cache.",

	HeroImage:   "images/chapter-11.webp",
	HeroCaption: "Cache = keep frequently used data nearby, so it's fast to get back.",

	Concepts: []ConceptItem{
		{Term: "Cache", Description: "A fast, small store that holds recently or frequently requested data close at hand."},
		{Term: "Hit and Miss", Description: "A hit answers from the cache without asking storage. A miss fetches from storage and caches it for next time."},
		{Term: "Invalidation", Description: "Removing or updating a cached value when the underlying data changes, so the cache never lies."},
	},

	BuildIt: BuildIt{
		Technique: "Self-Refine",
		Why:       "A cache bug does not crash. It returns a confident wrong number, which a passing test will not catch.",
		Source:    "The Prompt Report: Self-Criticism, Self-Refine",
		Prompts: []Prompt{
			{Label: "Build", Text: `Every balance enquiry hits the Vault's storage, even when the same balance was read a moment ago and hasn't changed since.

Add an in-memory cache in front of the Vault's reads, invalidated whenever a payment changes that account. An in-memory map, safe for concurrent access: no Redis, no cache library.

Done when a repeated enquiry is served from cache and a payment makes the next enquiry show the new balance.`},
			{Label: "Review", Text: `You have added an in-memory cache in front of the Vault's balance reads, invalidated when a payment changes an account.

Review that implementation as if you were trying to make it serve a stale balance. Walk every path that changes a balance and check whether it invalidates. Consider a concurrent read and write, a payment that rolled back, and a payment that touches two accounts at once.

Report what you found as a list, fix it, then run that same review again on the fixed version. Keep going until a review turns up nothing.

Done when a review finds nothing left, and you have told me how many rounds it took and which staleness bug the first version had.`},
			{Label: "Portal", Portal: true, Text: `The Portal shows a balance, sends money and lists a history, and looks like each was added the week it was needed. Make it presentable.

Judge it against every rule in peyva/portal/design.md, one at a time, and against the visual idea you committed to in chapter 0. Say whether the page still has that idea or has drifted into something anonymous.

Say whether each passes before changing anything. Fix what fails, judge again, and keep going until a pass finds nothing. Tell me what you fixed each round.

Done when a critique pass finds nothing to fix.`},
		},
	},
}
