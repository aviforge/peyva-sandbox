package content

var Chapter02 = ChapterContent{
	Number:     2,
	Slug:       "chapter-2",
	Title:      "Finding Peyva (Processes & Ports)",
	Subtitle:   "To reach someone, you need the right address and the receptionist lets you in.",
	Category:   "Foundations",
	Difficulty: "Beginner",
	EstTime:    "10 min",
	QuickTip:   "One process holds a given address and port at a time. That is what makes it somewhere you can reliably knock.",

	HeroImage:   "images/chapter-2.webp",
	HeroCaption: "A process is the running app. A port is the door callers knock on.",

	Intuition: []string{
		"Your computer runs many programs at once, like an office building holds many teams.",
		"Reaching one specific program needs an address (127.0.0.1) and a port. Something peyva doesn't have yet.",
		"This chapter gives it a door to knock on.",
	},

	Concepts: []ConceptItem{
		{Term: "Process ID (PID)", Description: "A number the OS assigns to identify one running program among many."},
		{Term: "Port", Description: "A number (like 9310) a process binds to, so the OS knows where to route calls."},
		{Term: "Binding", Description: "A process claiming a port on an address. By default the OS hands that pair to one process at a time, which is why the second copy is refused."},
		{Term: "Loopback (127.0.0.1)", Description: "The address meaning 'this same machine', used to reach a process running locally."},
		{Term: "Gateway", Description: "The component that holds the port and takes requests from outside, the system's one front door."},
	},

	UnderTheHood: []string{
		"Every peyva process has a PID and, once it listens, a port (e.g. peyva-api on 9310).",
		"A client reaches peyva at 127.0.0.1:9310, address plus port.",
		"The port is the door; the process behind it decides what happens when someone knocks.",
	},

	BuildIt: BuildIt{
		Intro:     "Build the Gateway: the way requests reach the system from outside.",
		Technique: "Constrain the scope",
		Why:       "Assistants overengineer by default: extra files, abstractions, flexibility nobody asked for. The remedy is to state the ceiling explicitly. Ask for a listener without one and you'll get an HTTP server with routes and JSON.",
		Source:    "Anthropic: Prompting best practices, Overeagerness",
		Prompt: "Build the Gateway, the way payment requests reach the system from outside.\n\n" +
			"For now it does one thing: claim TCP port 9310, accept connections, log one line per connection with the caller's address, then close it.\n\n" +
			"Do not add HTTP handling, routes, JSON parsing, request or response types, graceful shutdown, or any third-party package. The Gateway only owns a door right now. Everything else arrives in later chapters and would get in the way.\n\n" +
			"The right amount of complexity is the minimum that does what I just described.\n\n" +
			"Done when a request to port 9310 makes the Gateway log a connection, and starting a second copy fails with 'address already in use'.",
	},

	BreakIt: BreakIt{
		Intro: "An address and port have one owner at a time. See what happens when two processes want the same door.",
		Exercises: []string{
			"Start peyva once. It binds port 9310 successfully.",
			"Start a second copy while the first is still running. Watch it fail with 'address already in use'.",
			"Stop the first copy, then start the second again: now it succeeds. The port was free the moment its owner let go.",
		},
	},
}
