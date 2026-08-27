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
		"You cannot debug a distributed system by attaching to it. The failure happened an hour ago on one of three copies.",
		"Every log line about a payment carries its reference. One grep should reconstruct its whole path.",
		"Metrics say something is wrong. Logs say what happened. Correlated timestamps say where the time went.",
		"A health check that always says ok is a lie the load balancer believes. Degraded is a real answer.",
		"Watch the numbers that predict failure: outbox backlog, replication lag, unknown outcomes.",
		"Every useless line is one the on-call engineer reads past at 3am. Log less, with more per line.",
	},

	Concepts: []ConceptItem{
		{Term: "Observability", Description: "How much you can tell about what a running system is doing from the outside, without attaching a debugger to it."},
		{Term: "Structured Logs", Description: "One event per line, with named fields rather than prose, so a payment reference can be filtered on rather than searched for."},
		{Term: "Correlation ID", Description: "The value every line about one request carries, here the payment reference, so its path across processes can be reassembled."},
		{Term: "Metrics", Description: "Numbers over time: payments, failures, outbox backlog, replication lag, unknown outcomes."},
		{Term: "Traces", Description: "A request's journey across processes: where it went and how long each hop took. Correlated timestamps are the minimum version."},
		{Term: "Health Endpoint", Description: "Reports whether this process and what it depends on are working. Degraded when the Vault, the lease or the replica are not what they should be."},
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
