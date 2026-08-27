package content

var Chapter19 = ChapterContent{
	Number:     19,
	Slug:       "chapter-19",
	Title:      "Operating in Production",
	Subtitle:   "Building peyva is one thing. Running it reliably every day is another skill entirely.",
	Category:   "Operations",
	Difficulty: "Advanced",
	QuickTip:   "Deploy to one copy first: a health check catching a bad release there is far cheaper than catching it everywhere.",

	HeroImage:   "images/chapter-19.webp",
	HeroCaption: "Great restaurants aren't just built well. They are operated well every day. That's how customers stay happy.",

	Why: []string{
		"Settings are what change between runs: ports, addresses, passwords. Anything with only one correct value is code.",
		"A missing setting should stop the program at startup. A default hides it until it fails somewhere unrelated.",
		"Read every setting once, in one place. Read them all over the code and a default sneaks in somewhere.",
		"Release to one copy first and check it. A bad version then costs you one copy, not the whole system.",
		"A rollback is just releasing the previous version. The version on the health check tells you which copy runs what.",
		"Going back to old code does not go back to old data. The runbook must say whether this release can be undone at all.",
	},

	Concepts: []ConceptItem{
		{Term: "Health Check", Description: "An address that says whether this program and what it needs are working, and which version is running."},
		{Term: "Rolling Deployment", Description: "Releasing a new version one copy at a time, checking health in between."},
		{Term: "Release Rollback", Description: "Going back to the last working version. Not chapter 7's rollback, and it does not undo data the new version wrote."},
		{Term: "Runbook", Description: "A step by step guide for one known kind of incident. Commands and expected output, not prose."},
		{Term: "Config", Description: "The component that reads every setting from outside the code and checks it. A setting is what changes between runs; anything with one correct value stays in code."},
		{Term: "Fail Fast", Description: "Refusing to start when a setting is missing, instead of guessing and failing later somewhere else."},
	},

	BuildIt: BuildIt{
		Technique: "Structured output formatting",
		Why:       "A runbook read at 2am has to be commands, not prose. Hand over the exact skeleton you want back.",
		Source:    "Anthropic: Prompting best practices, Control the format of responses",
		Prompts: []Prompt{
			{Label: "Build", Text: `peyva reads settings from the environment in several places: ports, the primary a follower tracks, where each Vault keeps its file, the shared secret.

Build Config: read every setting once at startup, check it, hand it over. Nothing else reads the environment.

First, a table of every setting and which rule it falls under:

  config  differs between one run or one machine and the next
  config  a secret, which never belongs in the repository
  code    only one correct value exists, and changing it would be a bug

Two decimal places on money is not a setting. When unsure, ask whether I should be able to change it at 3am with no review.

A missing setting means naming it and stopping. Never a default that hides it.

Done when every setting is sorted, and starting any process without a required one names it and stops.`},
			{Label: "Build", Text: `peyva runs as a Vault, a follower, a Warden, three copies and a proxy, each with a health address.

Report a version string from every health address. Then write me a rollback runbook for a bad release of the copies, in exactly this shape:

  ## Symptom
  One line: how I know I have this problem.

  ## Check
  Numbered commands, each with the output that confirms it.

  ## Can this release be rolled back?
  One line: did it write anything the old version cannot read.

  ## Fix
  Numbered commands, no placeholders, through the runner.

  ## Verify
  One command and the exact output that means I am recovered.

Every command runs on {os} and can be pasted as is.

Done when a broken release to one copy fails its health check before the other two are touched, and Fix reverts it without improvising.`},
			{Label: "Portal", Portal: true, Text: `The Portal has peyva's address written into its pages. Move it to peyva/portal/config.js: one object, loaded first, holding the base URL and nothing else.

Empty means same origin, which is the ordinary case now.

Done when I can point the Portal at a different port by editing that one line.`},
		},
	},
}
