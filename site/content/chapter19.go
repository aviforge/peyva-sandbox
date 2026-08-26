package content

var Chapter19 = ChapterContent{
	Number:     19,
	Slug:       "chapter-19",
	Title:      "Operating in Production",
	Subtitle:   "Building peyva is one thing. Running it reliably every day is another skill entirely.",
	Category:   "Operations",
	Difficulty: "Advanced",
	EstTime:    "25 min",
	QuickTip:   "Deploy to one copy first: a health check catching a bad release there is far cheaper than catching it everywhere.",

	HeroImage:   "images/chapter-19.webp",
	HeroCaption: "Great restaurants aren't just built well. They are operated well every day. That's how customers stay happy.",

	Intuition: []string{
		"Keeping peyva running is a different job from building it.",
		"Rolling out changes safely, and noticing problems before customers do.",
		"And having a plan for when it goes wrong anyway.",
	},

	Concepts: []ConceptItem{
		{Term: "Health Check", Description: "An endpoint that confirms peyva and what it depends on are working."},
		{Term: "Rolling Deployment", Description: "Releasing a new version one copy at a time, instead of all at once."},
		{Term: "Rollback", Description: "Reverting to the previous working version when a release causes problems."},
		{Term: "Runbook", Description: "A step-by-step guide for handling a specific, known kind of incident."},
		{Term: "Config", Description: "The component that reads every setting from outside the code, checks it, and hands it over. Nothing else reads the environment."},
		{Term: "Secret", Description: "A setting that must never be in the repository: a password, a key, a token. Always config, never code."},
		{Term: "Fail Fast", Description: "Refusing to start when a setting is missing, rather than guessing and failing later somewhere unrelated."},
	},

	UnderTheHood: []string{
		"Users -> Load Balancer -> peyva Instances -> Database, continuously watched by Health Checks and Metrics & Logs feeding back into the loop.",
		"Day to day: Deploy Change -> Health Check -> Verify Metrics -> All Good? Yes: done. No: Rollback & Fix, then Postmortem & Improve.",
		"Config is what differs between one run and the next: ports, addresses, file paths. One build then runs anywhere.",
		"Code is what has one correct value: money to two decimal places, a balance that cannot go negative. Could someone change it at 3am with no review? If not, it is code.",
	},

	BuildIt: BuildIt{
		Intro:     "Build Config, then give the system a deployment and rollback story.",
		Technique: "Structured output formatting",
		Why:       "A runbook read at 2am has to be commands, not prose. Hand over the exact skeleton you want back.",
		Source:    "Anthropic: Prompting best practices, Control the format of responses",
		Prompts: []Prompt{
			{Label: "Config", Intro: "Sort the settings, then build Config.", Text: `peyva reads settings straight from the environment in several places: the port a copy listens on, the ports the proxy routes between, where the Vault keeps its file.

Build Config: it reads every setting once at startup, checks each, and hands them over. Nothing else reads the environment.

First sort what peyva already has. Give me a table of every setting, and for each one the rule it fell under:

  config  differs between one run or one machine and the next
  config  a secret, which never belongs in the repository
  code    only one value is ever correct, and changing it would be a bug

Two decimal places on money is not a setting. Neither is a balance that cannot go negative. When unsure, ask whether I should be able to change it at 3am with no review.

A missing or nonsense setting means naming it and exiting. Never a default that hides it.

Done when every setting has been sorted with the rule it fell under, and starting without a required one names it and stops.`},
			{Label: "Runbook", Intro: "The version string and the rollback runbook.", Text: `peyva runs as three copies behind a proxy, each exposing a health endpoint, and Config now supplies every setting.

Add a version string, set at build time, reported by that health endpoint. Then write me a rollback runbook.

Format the runbook exactly like this and nothing else:

  ## Symptom
  One line: how I know I have this problem.

  ## Check
  Numbered commands, one per line, with the output that confirms the diagnosis.

  ## Fix
  Numbered commands, one per line, no placeholders. Anything that starts or
  stops a copy goes through the runner, not a process ID I have to hunt for.

  ## Verify
  One command and the exact output that means I'm recovered.

Every command runs on {os}. Every line under Check, Fix and Verify is a command I can paste or an exact output I can compare against.

Done when deploying a deliberately broken version to one copy fails its health check before the other two are touched, and the runbook's Fix section reverts it without improvisation.`},
			{Label: "Portal", Portal: true, Intro: "The Portal stops having its address written into it.", Text: `The Portal reaches peyva at an address written into its pages. Move it to peyva/portal/config.js: one object, loaded before anything else, holding the base URL and nothing that is not a setting.

An empty base URL means same origin, which is the ordinary case now that peyva serves the page itself. Nothing in the Portal builds an address any other way.

Done when I can point the Portal at a different port by editing that one line, without touching a page.`},
		},
	},

	BreakIt: BreakIt{
		Intro: "Take a setting away, then break a release, and confirm peyva says so both times.",
		Exercises: []string{
			"Start a copy with its port unset. It names the missing setting and exits, rather than defaulting and failing somewhere unrelated.",
			"Point the Portal's base URL at a dead port. It says it cannot reach peyva, rather than showing a balance of nothing.",
			"Move money's two decimal places into config and set it to three. Nothing stops you. Put it back: a setting nobody should be able to change belongs in code.",
			"Deploy a version with an intentional bug to one copy. Confirm its health check fails before the other two are touched.",
			"Follow your own runbook to revert that copy, pasting every command rather than improvising, and time it.",
		},
	},
}
