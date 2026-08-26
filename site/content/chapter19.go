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
		"Every chapter so far has been about building peyva.",
		"This one is about keeping it running: rolling out changes safely, noticing problems before users do.",
		"And having a plan for when it goes wrong anyway.",
	},

	Concepts: []ConceptItem{
		{Term: "Health Check", Description: "An endpoint or process that continuously confirms peyva and its dependencies are actually working."},
		{Term: "Rolling Deployment", Description: "Releasing a new version gradually, one copy at a time, instead of all at once."},
		{Term: "Rollback", Description: "Quickly reverting to the previous working version when a new release causes problems."},
		{Term: "Runbook", Description: "A step-by-step guide for handling a specific, known kind of incident."},
		{Term: "Config", Description: "The component that reads every setting from outside the code, checks it, and hands it over. Nothing else reads the environment, the same way nothing but the Vault changes a balance."},
		{Term: "Secret", Description: "A setting that must never be in the repository: a password, a key, a token. Always config, never code."},
		{Term: "Fail Fast", Description: "Refusing to start when a setting is missing or makes no sense, instead of starting on a guess and failing later somewhere that looks unrelated."},
	},

	UnderTheHood: []string{
		"Users -> Load Balancer -> peyva Instances -> Database, continuously watched by Health Checks and Metrics & Logs feeding back into the loop.",
		"Day to day: Deploy Change -> Health Check -> Verify Metrics -> All Good? Yes: done. No: Rollback & Fix, then Postmortem & Improve.",
		"Config is what differs between one run and the next: ports, addresses, file paths, how long to wait before giving up. It comes from outside, so the same build runs anywhere.",
		"Code is what has only one correct value. Money to two decimal places. A balance that cannot go negative. Debits matching credits. Move one of those into config and you have not made peyva configurable, you have made its invariants optional.",
		"The test, when you cannot tell: could someone change this at 3am, alone, with no review? If the answer is no, it is code.",
	},

	BuildIt: BuildIt{
		Intro:     "Build Config, then give the system a deployment and rollback story.",
		Technique: "Structured output formatting",
		Why:       "A runbook read at 2am has to be commands, not prose. Hand over the exact skeleton you want back, and phrase it as what to produce rather than what to avoid, which is what actually steers the output.",
		Source:    "Anthropic: Prompting best practices, Control the format of responses",
		Prompt: `I start peyva with the runner, which brings up three copies behind the proxy, each exposing a health endpoint.

Settings are read straight from the environment in several places by now. Build Config: it reads every setting once at startup, checks each one, and hands them over. Nothing else reads the environment after that.

First sort what peyva already has. Give me a table of every setting, and for each one the rule it fell under:

  config  differs between one run or one machine and the next
  config  a secret, which never belongs in the repository
  code    only one value is ever correct, and changing it would be a bug

Two decimal places on money is not a setting. Neither is a balance that cannot go negative. If you are unsure which side something falls, ask whether I should be able to change it at 3am with no review.

A missing or nonsense setting means naming it and exiting. Never a default that hides it.

Then add a version string, set at build time, reported by the health endpoint. Then write me a rollback runbook.

Format the runbook exactly like this and nothing else:

  ## Symptom
  One line: how I know I have this problem.

  ## Check
  Numbered commands, one per line, with the output that confirms the diagnosis.

  ## Fix
  Numbered commands, one per line, copy-pasteable, no placeholders I have to
  think about. Anything that starts or stops a copy goes through the runner,
  not a process ID I have to hunt for.

  ## Verify
  One command and the exact output that means I'm recovered.

Every command runs on {os}. Every line under Check, Fix and Verify is either a command I can paste or an exact output I can compare against. It's 2am and I'm not making judgement calls.

Done when every setting has been sorted with its reason, starting without one names it and stops, deploying a deliberately broken version to one copy fails its health check before the other two are touched, and the runbook's Fix section reverts it without improvisation.`,
		UIIntro: "The Portal stops having its address written into it.",
		UIPrompt: `The Portal reaches peyva at an address written into its pages. Move it to peyva/portal/config.js: one object, loaded before anything else, holding the base URL and nothing that is not a setting.

An empty base URL means same origin, which is the ordinary case now that peyva serves the page itself. Nothing in the Portal builds an address any other way.

Done when I can point the Portal at a different port by editing that one line, without touching a page.`,
	},

	BreakIt: BreakIt{
		Intro: "Take a setting away, then break a release, and confirm peyva says so both times.",
		Exercises: []string{
			"Start a copy with its port unset. It names the missing setting and exits, instead of starting on a default and failing somewhere that looks unrelated.",
			"Point the Portal's base URL at a port nothing is listening on. The page says it cannot reach peyva, rather than showing a balance of nothing.",
			"Move the two decimal places on money into config, then set it to three. Nothing stops you, which is the point: a setting nobody should be able to change is one the code should have kept. Put it back.",
			"Deploy a version with an intentional bug to one copy first. Confirm its health check fails before it reaches the other two.",
			"Follow your own rollback runbook to revert that one copy, pasting every command rather than improvising, and time how long it actually takes.",
		},
	},
}
