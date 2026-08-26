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
	},

	UnderTheHood: []string{
		"Users -> Load Balancer -> peyva Instances -> Database, continuously watched by Health Checks and Metrics & Logs feeding back into the loop.",
		"Day to day: Deploy Change -> Health Check -> Verify Metrics -> All Good? Yes: done. No: Rollback & Fix, then Postmortem & Improve.",
	},

	BuildIt: BuildIt{
		Intro:     "No new component. Give the system a deployment and rollback story.",
		Technique: "Structured output formatting",
		Why:       "A runbook read at 2am has to be commands, not prose. Hand over the exact skeleton you want back, and phrase it as what to produce rather than what to avoid, which is what actually steers the output.",
		Source:    "Anthropic: Prompting best practices, Control the format of responses",
		Prompt: `I start peyva with the runner, which brings up three copies behind the proxy, each exposing a health endpoint.

Add a version string, set at build time, reported by that endpoint. Then write me a rollback runbook.

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

Done when deploying a deliberately broken version to one copy fails its health check before the other two are touched, and the runbook's Fix section reverts it without improvisation.`,
	},

	BreakIt: BreakIt{
		Intro: "Deploy a deliberately broken version and confirm your own process catches it.",
		Exercises: []string{
			"Deploy a version with an intentional bug to one copy first. Confirm its health check fails before it reaches the other two.",
			"Follow your own rollback runbook to revert that one copy, pasting every command rather than improvising, and time how long it actually takes.",
			"Write a one-paragraph postmortem: what broke, how it was caught, what would prevent it next time.",
		},
	},
}
