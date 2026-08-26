package content

var Chapter17 = ChapterContent{
	Number:     17,
	Slug:       "chapter-17",
	Title:      "See Everything: Observability",
	Subtitle:   "In a big system like peyva, you can't watch everything directly. Observability shows what's happening and helps you fix issues fast.",
	Category:   "Operations",
	Difficulty: "Advanced",
	EstTime:    "20 min",
	QuickTip:   "A metric tells you something's wrong; a log tells you what happened; a trace tells you where.",

	HeroImage:   "images/chapter-17.webp",
	HeroCaption: "Observability = dashboards + alarms + logs (and traces). It helps us understand the system like a pizza shop watches its kitchen.",

	Concepts: []ConceptItem{
		{Term: "Observability", Description: "How much you can tell about what a running system is doing from the outside, without attaching a debugger to it."},
		{Term: "Metrics", Description: "Numbers over time: requests/sec, latency, queue size, error rate, CPU, memory."},
		{Term: "Logs", Description: "Detailed events: a login failed, a payment was processed. The specific story of what happened."},
		{Term: "Traces", Description: "A request's journey across services: one payment request's path through five different parts of the system."},
		{Term: "Alerts", Description: "Notifications when something is wrong: error rate over 5%, queue size too high."},
	},

	BuildIt: BuildIt{
		Intro:     "Every component learns to say what it's doing.",
		Technique: "Style Prompting",
		Why:       "Naming the reader and the moment they read it changes both what gets logged and how it reads.",
		Source:    "The Prompt Report: Zero-Shot, Style Prompting",
		Prompts: []Prompt{
			{Label: "Build", Text: `The system logs by printing to the console and has no metrics. When something breaks there's no way to find out what happened.

Write the observability for one specific reader: an on-call engineer, woken at 3am, who didn't write this code, is looking at yesterday's records for one failed payment, and cannot attach a debugger or reproduce it.

Give them structured logs carrying the payment reference so they can filter to a single payment, counters for payments, failures and Courier backlog, and a health endpoint reporting real Vault and Courier state rather than a hardcoded ok.

For each thing you log, ask whether it helps that person at 3am. If it doesn't, leave it out. Noise costs them more than a missing line.

Done when I can take one failed payment reference and reconstruct its whole path from logs alone, and the health endpoint reports degraded when the Vault's second copy is genuinely unreachable.`},
		},
	},

	BreakIt: BreakIt{
		Intro: "Use observability to actually catch a problem, not just log noise.",
		Exercises: []string{
			"Deliberately break one dependency (e.g. disconnect the replica from Chapter 15) and confirm the error rate metric visibly rises.",
			"Check /health during the failure and confirm it reflects the real degraded state, not a blind 'OK'.",
			"Read back the structured logs for a specific failed transfer and confirm you can reconstruct exactly what happened, in order.",
		},
	},
}
