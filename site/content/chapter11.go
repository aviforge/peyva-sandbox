package content

var Chapter11 = ChapterContent{
	Number:     11,
	Slug:       "chapter-11",
	Title:      "Sharing Work: Caching",
	Subtitle:   "Caching is like keeping popular items ready in advance, so we don't make (or check) everything again.",
	Category:   "System Design",
	Difficulty: "Intermediate",
	QuickTip:   "A stale cache is worse than no cache.",

	HeroImage:   "images/chapter-11.webp",
	HeroCaption: "Cache = keep frequently used data nearby, so it's fast to get back.",

	Why: []string{
		"A cache is a second copy, and any second copy can disagree with the first.",
		"Only the program that sees every write can clear the cache correctly. So the cache lives in the Vault, not in the copies.",
		"Clear it after the write is saved, never before. Clear it early and another read puts the old value straight back.",
		"A payment changes two accounts. Clear both.",
		"A balance on a page may be a moment old. The balance checked before taking money may not be cached at all.",
		"A cache bug never crashes. It hands back a fast, confident, wrong number.",
	},

	Concepts: []ConceptItem{
		{Term: "Cache", Description: "A small fast store holding what was asked for recently, so the next ask is quick."},
		{Term: "Hit and Miss", Description: "A hit answers from the cache. A miss goes to storage and keeps the answer for next time."},
		{Term: "Invalidation", Description: "Clearing or updating a cached value when the real one changes, so the cache never lies."},
		{Term: "Cache Placement", Description: "Which program holds the cache. Only one that sees every write can clear it correctly, so balances are cached in the Vault, not the copies."},
		{Term: "Read-Your-Writes", Description: "Someone who just paid must see the new balance. A cache that can hand them the old one has broken a promise they will notice."},
	},

	BuildIt: BuildIt{
		Technique: "Self-Refine",
		Why:       "A cache bug does not crash. It returns a confident wrong number, which a passing test will not catch.",
		Source:    "The Prompt Report: Self-Criticism, Self-Refine",
		Prompts: []Prompt{
			{Label: "Build", Text: `Every balance enquiry reaches the Vault's storage, even when the same balance was read a moment ago and hasn't changed since. Three copies of the Gateway and Teller sit in front of the Vault, which runs as its own process.

Add an in-memory cache of balances inside the Vault's process, in front of its reads, invalidated whenever a payment changes an account. An in-memory map, safe for concurrent access: no Redis, no cache library. The cache serves enquiries only; the check before a debit reads inside the transaction, never from the cache.

Before writing it, say why the cache cannot live in the copies, in two sentences.

Done when a repeated enquiry is served from cache, a payment makes the next enquiry from any copy show the new balance, and the balance check inside a payment never reads the cache.`},
			{Label: "Review", Text: `You have added an in-memory cache inside the Vault's process in front of its balance reads, invalidated when a payment changes an account.

Review that implementation as if you were trying to make it serve a stale balance. Walk every path that changes a balance and check whether it invalidates. Consider a concurrent read and write, a payment that rolled back after the cache was dropped, a payment that touches two accounts at once, and a read that refills the cache between the invalidation and the commit.

Report what you found as a list, fix it, then run that same review again on the fixed version. Keep going until a review turns up nothing.

Done when a review finds nothing left, and you have told me how many rounds it took and which staleness bug the first version had.`},
			{Label: "Portal", Portal: true, Text: `The Portal shows a balance, sends money and lists a history, and looks like each was added the week it was needed. Make it presentable.

Judge it against every rule in peyva/portal/design.md, one at a time, and against the visual idea you committed to in chapter 0. Say whether the page still has that idea or has drifted into something anonymous.

Say whether each passes before changing anything. Fix what fails, judge again, and keep going until a pass finds nothing. Tell me what you fixed each round.

Done when a critique pass finds nothing to fix.`},
		},
	},
}
