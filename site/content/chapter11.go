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
		{Term: "Cache Hit", Description: "The data was already in the cache, so the answer comes back without the database being asked at all."},
		{Term: "Cache Miss", Description: "The data wasn't cached. peyva fetches it from the database and stores it in the cache for next time."},
		{Term: "Invalidation", Description: "Removing or updating a cached value when the underlying data changes, so the cache never lies."},
	},

	BuildIt: BuildIt{
		Intro:     "The Vault learns to answer faster: without ever answering wrong.",
		Technique: "Self-Refine",
		Why:       "A cache bug does not crash. It returns a confident wrong number, which a passing test will not catch.",
		Source:    "The Prompt Report: Self-Criticism, Self-Refine",
		Prompts: []Prompt{
			{Label: "Build", Intro: "The cache.", Text: `Every balance enquiry hits the Vault's storage, even when the same balance was read a moment ago and hasn't changed since.

Add an in-memory cache in front of the Vault's reads, invalidated whenever a payment changes that account. An in-memory map, safe for concurrent access: no Redis, no cache library.

Done when a repeated enquiry is served from cache and a payment makes the next enquiry show the new balance.`},
			{Label: "Review", Intro: "Review it for stale balances.", Text: `You have added an in-memory cache in front of the Vault's balance reads, invalidated when a payment changes an account.

Review that implementation as if you were trying to make it serve a stale balance. Walk every path that changes a balance and check whether it invalidates. Consider a concurrent read and write, a payment that rolled back, and a payment that touches two accounts at once.

Report what you found as a list, fix it, then run that same review again on the fixed version. Keep going until a review turns up nothing.

Done when a review finds nothing left, and you have told me how many rounds it took and which staleness bug the first version had.`},
			{Label: "Portal", Portal: true, Intro: "The Portal stops looking like a first draft.", Text: `The Portal shows a balance, sends money and lists a history, and looks like each was added the week it was needed. Make it presentable.

Judge what you have against these, one at a time: whose wallet this is never has to be guessed at; the menu makes it obvious where the customer is and what else they can do; it reads at a glance on a phone; money is aligned and always shows two decimals; the balance is the most prominent thing on its own screen; every action says what happened; nothing shifts as data loads; it is legible in both light and dark.

Say whether each passes before changing anything. Fix what fails, judge again, and keep going until a pass finds nothing. Tell me what you fixed each round.

Still plain HTML and CSS.

Done when a critique pass finds nothing to fix.`},
		},
	},

	BreakIt: BreakIt{
		Intro: "Prove the invalidation actually works.",
		Exercises: []string{
			"Read Alice's balance a few thousand times with the cache on, then again with it off, and compare the totals. A single hit is far too fast to see. The saving only shows up in the aggregate.",
			"Make a transfer that changes Alice's balance, then check it again. It must reflect the new balance, not the stale cached one.",
			"Deliberately skip the invalidation step and repeat the test. Watch the cache confidently return the wrong balance.",
		},
	},
}
