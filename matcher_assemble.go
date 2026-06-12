package matchs

import "strings"

const (
	ruleTypeExclude = iota
	ruleTypeOrdered
)

type rule struct {
	rawRule  string
	words    []string
	ruleType int
}

func (r *rule) match(text string) bool {
	if r.ruleType == ruleTypeOrdered {
		index := 0
		for _, w := range r.words {
			last := strings.Index(text[index:], w)
			if last < 0 {
				return false
			}
			index += last + len(w)
		}
		return true
	} else if r.ruleType == ruleTypeExclude {
		if strings.Contains(text, r.words[0]) {
			for i := 1; i < len(r.words); i++ {
				if strings.Contains(text, r.words[i]) {
					return false
				}
			}
			return true
		}
	}

	return false
}

// AssembleMatcher matches composite rules.
//
// Rules containing "|" require words to appear in an ordered relationship.
// Rules containing "#" require the first word to appear and later words to be
// absent. This matcher reports matched rule strings but does not perform text
// replacement.
type AssembleMatcher struct {
	rules []*rule
}

// AssembleMather is the old name for AssembleMatcher.
//
// Deprecated: use AssembleMatcher.
type AssembleMather = AssembleMatcher

// NewAssembleMatcher creates an empty AssembleMatcher.
func NewAssembleMatcher() *AssembleMatcher {
	return &AssembleMatcher{}
}

// NewAssembleMather creates an empty AssembleMatcher.
//
// Deprecated: use NewAssembleMatcher.
func NewAssembleMather() *AssembleMatcher {
	return NewAssembleMatcher()
}

// Build adds composite rules to the matcher.
//
// Rules containing "|" are treated as ordered rules. Rules containing "#" are
// treated as exclusion rules. Other rules are ignored by this matcher.
func (a *AssembleMatcher) Build(words []string) {

	for _, w := range words {

		if strings.Contains(w, "|") {
			a.rules = append(a.rules, &rule{
				rawRule:  w,
				words:    strings.Split(w, "|"),
				ruleType: ruleTypeOrdered,
			})
		} else if strings.Contains(w, "#") {
			a.rules = append(a.rules, &rule{
				rawRule:  w,
				words:    strings.Split(w, "#"),
				ruleType: ruleTypeExclude,
			})
		}
	}
}

// Match returns composite rules matched by text.
//
// onlyOne stops after the first matched composite rule. repl is accepted to
// satisfy Matcher but is not used because AssembleMatcher does not perform text
// replacement. The replacement text return value keeps the original text.
func (a *AssembleMatcher) Match(text string, onlyOne bool, repl rune) (word []string, desensitization string) {
	desensitization = text
	for _, rule := range a.rules {
		if rule.match(text) {
			word = append(word, rule.rawRule)
			if onlyOne {
				return
			}
		}
	}
	return
}
