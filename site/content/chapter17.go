package content

var Chapter17 = ChapterContent{
	Number:     17,
	Slug:       "chapter-17",
	Title:      "See Everything: Observability",
	Subtitle:   "In a big system like peyva, you can't watch everything directly. Observability shows what's happening and helps you fix issues fast.",
	Category:   "Operations",
	Difficulty: "Advanced",
	QuickTip:   "A metric tells you something's wrong; a log tells you what happened; a trace tells you where.",

	HeroImage:   "images/chapter-17.webp",
	HeroCaption: "Observability = dashboards + alarms + logs (and traces). It helps us understand the system like a pizza shop watches its kitchen.",

	Why: []string{
		"You cannot debug this by stopping it. The failure was an hour ago, on one of three copies, in a request that has finished.",
		"Every line about a payment carries its reference. One search should then show its whole path.",
		"Numbers tell you something is wrong. Lines tell you what happened. Their timestamps tell you where the time went.",
		"A health check that always says ok is a lie the load balancer believes. 'Struggling' is a real answer.",
		"Watch what warns you early: jobs waiting, how far behind the copy is, payments with an unknown outcome.",
		"Every useless line is one someone reads past at 3am. Log less, and put more in each line.",
	},

	Concepts: []ConceptItem{
		{Term: "Observability", Description: "How much you can tell about a running system from the outside, without stopping it or attaching a debugger."},
		{Term: "Structured Logs", Description: "One event per line with named fields instead of sentences, so you can filter by payment reference."},
		{Term: "Correlation ID", Description: "The value every line about one request carries, here the payment reference, so you can piece its path together."},
		{Term: "Metrics", Description: "Numbers over time: payments, failures, unknown outcomes, jobs waiting, how far behind the second copy is."},
		{Term: "Traces", Description: "Where a request went and how long each hop took. Lines sharing a reference, with timestamps, are the simple version."},
		{Term: "Health Endpoint", Description: "Says whether this program and the things it needs are working. Never a hardcoded ok."},
	},

	BuildIt: BuildIt{
		Technique: "Style Prompting",
		Why:       "Naming the reader and the moment they read it changes both what gets logged and how it reads.",
		Source:    "The Prompt Report: Zero-Shot, Style Prompting",
		Prompts: []Prompt{
			{Label: "Build", Text: `The system logs by printing to the console and has no metrics. When something breaks there's no way to find out what happened.

Write the observability for one specific reader: an on-call engineer, woken at 3am, who didn't write this code, is looking at yesterday's records for one failed payment, and cannot attach a debugger or reproduce it.

Give them structured logs carrying the payment reference on every line that payment touches, in the proxy, the copy, the Vault and the replica, each with a timestamp precise enough to see where the time went. Counters for payments, failures, unknown outcomes, Courier backlog and replication lag. A health endpoint on every process reporting real state: for a copy, whether the Vault it last used answered; for a Vault, whether it holds the lease and how far its replica is behind; never a hardcoded ok.

For each thing you log, ask whether it helps that person at 3am. If it doesn't, leave it out. Noise costs them more than a missing line.

Done when I can take one failed payment reference and reconstruct its whole path and timing from logs alone, and the health endpoint reports degraded when the replica is genuinely unreachable or the lease is not held.`},
		},
	},
}
