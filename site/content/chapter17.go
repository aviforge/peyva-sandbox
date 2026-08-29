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
		Why:       "Name the reader and the voice you want. Writing logs for one on-call engineer at 3am changes what gets logged, not only how it reads.",
		Source:    "The Prompt Report: Zero-Shot, Style Prompting",
		Prompts: []Prompt{
			{Label: "Build", Text: `The system prints to the console and counts nothing. When something breaks, there is no way to find out what happened.

Write it for one reader: an on-call engineer, woken at 3am, who did not write this code, looking at yesterday's records for one failed payment.

One event per line, with the payment reference on every line that payment touches, and timestamps fine enough to show where the time went. Counts of payments, failures, unknown outcomes, jobs waiting, and how far behind the follower is. A health address on every process that reports what is actually true, never a hardcoded ok.

If a line would not help that person at 3am, leave it out.

Done when one failed reference rebuilds its whole path from the logs alone, and health reports trouble when the follower is unreachable or the lease is not held.`},
		},
	},
}
