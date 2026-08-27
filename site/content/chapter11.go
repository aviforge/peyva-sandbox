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
		"A cache is a second copy of a value, and every second copy can disagree with the first. The whole difficulty of caching is keeping that disagreement short, bounded, and invisible where it matters.",
		"Where the cache lives decides who can invalidate it. A cache inside each of three copies can only be invalidated by the copy that did the write; the other two keep serving the old balance. The only process that sees every write to a balance is the one that owns it, so the cache goes there.",
		"Invalidation must ride the same transaction as the write, or happen after it commits, never before. Drop the cached value, then have the transaction roll back, and the next read repopulates from the truth; drop it before the commit, and a concurrent read can refill the cache with the old value a moment before the new one lands.",
		"A payment touches two accounts. Invalidating one and not the other is the bug a test with one account will never find.",
		"Reads of money are not all equal. A balance shown on a page can be a few milliseconds old; a balance checked before a debit cannot be cached at all, because the check and the write have to be one unit. The cache serves the first kind and stays out of the second.",
		"The failure mode of a cache is not an error. It is a confident, fast, wrong number, which is why cache bugs survive test suites that only check for crashes.",
	},

	Concepts: []ConceptItem{
		{Term: "Cache", Description: "A fast, small store that holds recently or frequently requested data close at hand."},
		{Term: "Hit and Miss", Description: "A hit answers from the cache without asking storage. A miss fetches from storage and caches it for next time."},
		{Term: "Invalidation", Description: "Removing or updating a cached value when the underlying data changes, so the cache never lies."},
		{Term: "Cache Placement", Description: "Which process holds the cache. Only a process that sees every write can invalidate correctly, so a cache of balances lives with the Vault, not in the copies."},
		{Term: "Read-Your-Writes", Description: "A customer who just paid must see the new balance. Any cache that can serve them the old one has broken a promise they will notice."},
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
