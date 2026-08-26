package content

// System is an operating system a reader can follow the book on.
//
// It exists for one reason: the runner in chapter 10 is a shell script, and a
// shell script is the one thing in peyva that cannot be written once and work
// everywhere. Everything else the reader builds is in their chosen language and
// does not care what it runs on.
type System struct {
	ID string `json:"id"`
	// Name is what the picker shows.
	Name string `json:"name"`
	// Prompt is what gets substituted into a prompt at {os}. It names the shell
	// as well as the system, because "a script for Windows" is still a choice
	// between PowerShell and a batch file, and the two kill processes
	// differently.
	Prompt string `json:"prompt"`
}

// DefaultSystem is baked into every generated page, for the same reason
// DefaultLanguage is: a prompt that reaches an assistant without naming one
// gets a script for whatever the assistant assumed, and the reader finds out
// when it fails to run.
const DefaultSystem = "windows"

var Systems = []System{
	{ID: "windows", Name: "Windows", Prompt: "Windows, in PowerShell"},
	{ID: "macos", Name: "macOS", Prompt: "macOS, in bash or zsh"},
	{ID: "linux", Name: "Linux", Prompt: "Linux, in bash"},
}

// SystemByID returns the named system, falling back to the default rather than
// failing, so a saved choice that no longer exists returns the reader to the
// default instead of leaving them with none.
func SystemByID(id string) System {
	for _, s := range Systems {
		if s.ID == id {
			return s
		}
	}
	for _, s := range Systems {
		if s.ID == DefaultSystem {
			return s
		}
	}
	return Systems[0]
}
