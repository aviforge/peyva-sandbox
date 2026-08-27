package content

var Chapter02 = ChapterContent{
	Number:     2,
	Slug:       "chapter-2",
	Title:      "Finding Peyva (Processes & Ports)",
	Subtitle:   "To reach someone, you need the right address and the receptionist lets you in.",
	Category:   "Foundations",
	Difficulty: "Beginner",
	QuickTip:   "One process holds a given address and port at a time. That is what makes it somewhere you can reliably knock.",

	HeroImage:   "images/chapter-2.webp",
	HeroCaption: "A process is the running app. A port is the door callers knock on.",

	Why: []string{
		"A port is a rendezvous, not a pipe. The OS hands an address and port pair to one listening process, so a caller who knows the pair knows exactly which process will answer, and a second copy that asks for the same pair is refused.",
		"That refusal is a feature. 'Address already in use' is the OS telling you that two processes cannot both be the front door, which is the first hint that running more than one copy needs someone in front of them.",
		"Binding to 127.0.0.1 and binding to 0.0.0.0 are different promises. The first is reachable only from this machine; the second from anything that can route to it, which includes the firewall's opinion.",
		"A connection is a pair of endpoints and some state the OS holds on both sides. Accepting one costs memory and a file descriptor, and a process that never closes them runs out of both.",
		"The one front door is a design decision, not an accident. Everything outside reaches the system through it, so it is the one place that can later count, refuse, or redirect every request.",
	},

	Concepts: []ConceptItem{
		{Term: "Process ID (PID)", Description: "A number the OS assigns to identify one running program among many."},
		{Term: "Port", Description: "A number (like 9310) a process binds to, so the OS knows where to route calls."},
		{Term: "Binding", Description: "A process claiming a port on an address. The OS hands that pair to one process at a time, so a second copy is refused."},
		{Term: "Loopback (127.0.0.1)", Description: "The address meaning 'this same machine', used to reach a process running locally."},
		{Term: "Gateway", Description: "The component that holds the port and takes requests from outside, the system's one front door."},
	},

	BuildIt: BuildIt{
		Technique: "Constrain the scope",
		Why:       "Ask for a listener without stating a ceiling and you get an HTTP server with routes and JSON.",
		Source:    "Anthropic: Prompting best practices, Overeagerness",
		Prompts: []Prompt{
			{Label: "Build", Text: `peyva is one process holding the Vault, and the only way to reach it is to run it yourself. Build the Gateway, the way payment requests reach the system from outside.

For now it does one thing: claim TCP port 9310, accept connections, log one line per connection with the caller's address, then close it.

Do not add HTTP handling, routes, JSON parsing, request or response types, graceful shutdown, or any third-party package. The Gateway only owns a door right now. Everything else arrives in later chapters and would get in the way.

The right amount of complexity is the minimum that does what I just described.

Done when a request to port 9310 makes the Gateway log a connection, and starting a second copy fails with 'address already in use'.`},
		},
	},
}
