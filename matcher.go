package matchs

// Matcher defines the common behavior for keyword matchers.
//
// Build loads matcher rules into the receiver. Implementations in this package
// append to their current state, so callers should create a new matcher when
// they need to rebuild from a clean rule set.
//
// Match scans text and returns matched rule names plus a replacement text.
// onlyOne stops after the first matched rule when the implementation supports
// early return. repl is used by matchers that support text replacement.
type Matcher interface {
	// Build loads words or rules into the matcher.
	Build(words []string)

	// Match returns matched rules and replacement text for text.
	Match(text string, onlyOne bool, repl rune) ([]string, string)
}
