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
			{Label: "Config", Text: `peyva reads settings straight from the environment in several places: the ports, the primary a replica follows, where each Vault keeps its file, the secret the copies present.

Build Config: it reads every setting once at startup, checks each, and hands them over. Nothing else reads the environment.

First sort what peyva already has. Give me a table of every setting and the rule it fell under:

  config  differs between one run or one machine and the next
  config  a secret, which never belongs in the repository
  code    only one correct value exists, and changing it would be a bug

Two decimal places on money is not a setting. Neither is a balance that cannot go negative. When unsure, ask whether I should be able to change it at 3am with no review.

A missing or nonsense setting means naming it and stopping. Never a default that hides it.

Done when every setting is sorted with the rule it fell under, and starting any process without a required one names it and stops.`},
			{Label: "Runbook", Text: `peyva runs as a Vault, a replica, a Warden, three copies and a proxy, each with a health address.

Add a version string, set at build time, reported by every health address. Then write me a rollback runbook for a bad release of the copies.

Use exactly this shape and nothing else:

  ## Symptom
  One line: how I know I have this problem.

  ## Check
  Numbered commands, one per line, with the output that confirms it.

  ## Can this release be rolled back?
  One line: whether it changed anything the previous version cannot read.

  ## Fix
  Numbered commands, one per line, no placeholders. Anything that starts or
  stops a process goes through the runner, not a process id I hunt for.

  ## Verify
  One command and the exact output that means I am recovered.

Every command runs on {os}, and every line is either a command I can paste or an output I can compare against.

Done when releasing a broken version to one copy fails its health check before the other two are touched, and the Fix section reverts it without improvising.`},
			{Label: "Portal", Portal: true, Text: `The Portal reaches peyva at an address written into its pages. Move it to peyva/portal/config.js: one object, loaded first, holding the base URL and nothing that is not a setting.

An empty base URL means same origin, which is the ordinary case now that peyva serves the page itself. Nothing builds an address any other way.

Done when I can point the Portal at a different port by editing that one line, without touching a page.`},
		},
	},
}
