package content

// Language is one of the languages a reader can build peyva in. The prompts
// themselves never name a language: they describe what the component has to
// do, and the reader's assistant writes it idiomatically for whichever of
// these is selected.
//
// Adding one costs a single entry here, which is the whole point of keeping
// the prompts neutral. A design that forked the prompt text per language
// would need twenty-one rewrites for every language added, and they would
// drift apart as chapters changed.
type Language struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// DefaultLanguage is baked into every generated page. A prompt that reaches
// an assistant without naming a language is worse than useless: the assistant
// picks one on its own, and can pick differently on different chapters, so
// chapter five's code will not build on chapter four's. The reader finds out
// hours later. There is deliberately no "unselected" state.
const DefaultLanguage = "go"

// Languages are offered in the sidebar. Order is roughly by how likely a
// reader is to reach for them, not alphabetical.
var Languages = []Language{
	{ID: "go", Name: "Go"},
	{ID: "python", Name: "Python"},
	{ID: "javascript", Name: "JavaScript"},
	{ID: "typescript", Name: "TypeScript"},
	{ID: "java", Name: "Java"},
	{ID: "csharp", Name: "C#"},
	{ID: "cpp", Name: "C++"},
	{ID: "rust", Name: "Rust"},
	{ID: "ruby", Name: "Ruby"},
	{ID: "php", Name: "PHP"},
	{ID: "kotlin", Name: "Kotlin"},
	{ID: "swift", Name: "Swift"},
}

// LanguageByID returns the named language, falling back to the default rather
// than failing. A missing id means a reader's saved choice no longer exists,
// which should quietly return them to Go, not leave them with no language.
func LanguageByID(id string) Language {
	for _, l := range Languages {
		if l.ID == id {
			return l
		}
	}
	for _, l := range Languages {
		if l.ID == DefaultLanguage {
			return l
		}
	}
	return Languages[0]
}
