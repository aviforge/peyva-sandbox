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
		"Config is what differs between runs: ports, addresses, secrets. Anything with one correct value is code.",
		"A missing setting stops the process at startup. A default hides it until it fails somewhere unrelated.",
		"Read settings once, in one place. Scattered reads are scattered places a default creeps in.",
		"Deploy one copy first and check health. A bad release becomes one degraded copy, not an outage.",
		"A rollback is a deploy of the previous version. The version string on the health endpoint shows who runs what.",
		"Rolling back code does not roll back data. The runbook must say whether this release can be rolled back at all.",
	},

	Concepts: []ConceptItem{
		{Term: "Health Check", Description: "An endpoint that confirms peyva and what it depends on are working, and says which version is running."},
		{Term: "Rolling Deployment", Description: "Releasing a new version one copy at a time, checking health between, instead of all at once."},
		{Term: "Release Rollback", Description: "Reverting to the previous working version when a release causes problems. Not chapter 7's transaction rollback, and it does not undo data the new version wrote."},
		{Term: "Runbook", Description: "A step-by-step guide for handling a specific, known kind of incident. Commands and expected output, not prose."},
		{Term: "Config", Description: "The component that reads every setting from outside the code and checks it. A setting is what differs between one run and the next, secrets included; anything with one correct value stays in code."},
		{Term: "Fail Fast", Description: "Refusing to start when a setting is missing, rather than guessing and failing later somewhere unrelated."},
	},

	BuildIt: BuildIt{
		Technique: "Structured output formatting",
		Why:       "A runbook read at 2am has to be commands, not prose. Hand over the exact skeleton you want back.",
		Source:    "Anthropic: Prompting best practices, Control the format of responses",
		Prompts: []Prompt{
			{Label: "Config", Text: `peyva reads settings straight from the environment in several places: the port each process listens on, the ports the proxy routes between, the Vault's and the Warden's ports, the primary a replica follows, where each Vault keeps its file, the shared secret the copies present to the Vault.

Build Config: it reads every setting once at startup, checks each, and hands them over. Nothing else reads the environment.

First sort what peyva already has. Give me a table of every setting, and for each one the rule it fell under:

  config  differs between one run or one machine and the next
  config  a secret, which never belongs in the repository
  code    only one value is ever correct, and changing it would be a bug

Two decimal places on money is not a setting. Neither is a balance that cannot go negative, nor the lease length's relationship to its renewal interval. When unsure, ask whether I should be able to change it at 3am with no review.

A missing or nonsense setting means naming it and exiting. Never a default that hides it.

Done when every setting has been sorted with the rule it fell under, and starting any process without a required one names it and stops.`},
			{Label: "Runbook", Text: `peyva runs as a Vault, a replica, a Warden, three copies and a proxy, each exposing a health endpoint, and Config now supplies every setting.

Add a version string, set at build time, reported by every health endpoint. Then write me a rollback runbook for a bad release of the copies.

Format the runbook exactly like this and nothing else:

  ## Symptom
  One line: how I know I have this problem.

  ## Check
  Numbered commands, one per line, with the output that confirms the diagnosis.

  ## Can this release be rolled back?
  One line: whether the release changed anything the previous version cannot read, and how you know.

  ## Fix
  Numbered commands, one per line, no placeholders. Anything that starts or
  stops a process goes through the runner, not a process ID I have to hunt for.

  ## Verify
  One command and the exact output that means I'm recovered.

Every command runs on {os}. Every line under Check, Fix and Verify is a command I can paste or an exact output I can compare against.

Done when deploying a deliberately broken version to one copy fails its health check before the other two are touched, and the runbook's Fix section reverts it without improvisation.`},
			{Label: "Portal", Portal: true, Text: `The Portal reaches peyva at an address written into its pages. Move it to peyva/portal/config.js: one object, loaded before anything else, holding the base URL and nothing that is not a setting.

An empty base URL means same origin, which is the ordinary case now that peyva serves the page itself. Nothing in the Portal builds an address any other way.

Done when I can point the Portal at a different port by editing that one line, without touching a page.`},
		},
	},
}
