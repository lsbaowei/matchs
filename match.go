package matchs

import "strings"

const (
	// DFA identifies the plain keyword matcher.
	DFA = 0
	// ASSEMBLE identifies the composite rule matcher.
	ASSEMBLE = 1
	// REGEXP identifies the regexp matcher.
	REGEXP = 2

	// REGEXP_PREFIX marks a rule as a regexp rule.
	REGEXP_PREFIX = "reg@"
)

var matcherOrder = []int{DFA, ASSEMBLE, REGEXP}

// MatchService groups all supported matcher implementations behind one API.
//
// Build dispatches rules to plain keyword, composite-rule, or regexp matchers.
// Match aggregates the results from the built matchers.
type MatchService struct {
	matchers map[int]Matcher
}

// NewMatchService creates an empty MatchService.
func NewMatchService() *MatchService {
	return &MatchService{
		matchers: make(map[int]Matcher),
	}
}

// Build classifies words into matcher-specific rule sets and builds matchers.
//
// Rules with REGEXP_PREFIX are handled as regexp rules. Rules containing "|" or
// "#" are handled as composite rules. All other rules are handled as plain
// keywords.
func (m *MatchService) Build(words []string) {
	var (
		dfaList      []string
		assembleList []string
		regexpList   []string
	)

	for i := 0; i < len(words); i++ {
		if strings.HasPrefix(words[i], REGEXP_PREFIX) {
			regexpList = append(regexpList, words[i][len(REGEXP_PREFIX):])
		} else if strings.Contains(words[i], "|") || strings.Contains(words[i], "#") {
			assembleList = append(assembleList, words[i])
		} else {
			dfaList = append(dfaList, words[i])
		}
	}

	if len(dfaList) > 0 {
		matcher := NewDFAMatcher()
		matcher.Build(dfaList)
		m.matchers[DFA] = matcher
	}

	if len(assembleList) > 0 {
		matcher := NewAssembleMatcher()
		matcher.Build(assembleList)
		m.matchers[ASSEMBLE] = matcher
	}

	if len(regexpList) > 0 {
		matcher := NewRegexpMatcher()
		matcher.Build(regexpList)
		m.matchers[REGEXP] = matcher
	}
}

// Match scans text with all built matchers.
//
// onlyOne asks each matcher to stop after its first match when supported. repl
// is passed to matchers that support replacement. It returns matched rule names
// and the replacement text reported by the matchers.
func (m *MatchService) Match(text string, onlyOne bool, repl rune) (sensitiveWords []string, replaceText string) {
	replaceText = text
	for _, matcherType := range matcherOrder {
		x, ok := m.matchers[matcherType]
		if !ok {
			continue
		}

		ret, replaced := x.Match(replaceText, onlyOne, repl)
		replaceText = replaced
		for _, word := range ret {
			sensitiveWords = append(sensitiveWords, word)
		}
		if onlyOne && len(ret) > 0 {
			return
		}
	}
	return
}

/*-------------other util-------------------*/

func isASCIISpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

// TrimString returns s without leading and trailing ASCII space.
func TrimString(s string) string {
	for len(s) > 0 && isASCIISpace(s[0]) {
		s = s[1:]
	}
	for len(s) > 0 && isASCIISpace(s[len(s)-1]) {
		s = s[:len(s)-1]
	}
	return s
}
