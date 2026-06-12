package matchs

import (
	"fmt"

	"github.com/dlclark/regexp2"
)

type regRule struct {
	Tag      string
	metaData string
	reg      *regexp2.Regexp
}

// NewRegRule compiles str into a regexp rule.
//
// str must not include REGEXP_PREFIX. The returned rule reports matches with
// REGEXP_PREFIX added back to the rule tag.
func NewRegRule(str string) (*regRule, error) {
	if r, err := regexp2.Compile(str, 0); err == nil {
		return &regRule{
			Tag:      REGEXP_PREFIX + str,
			metaData: str,
			reg:      r,
		}, nil
	} else {
		return nil, fmt.Errorf("%s comiple regexp error:%s ", str, err.Error())
	}
}

// MatchAll returns this regexp rule's tag when text matches.
//
// Only the first regexp match is considered. The map value is the matched
// substring, but callers in this package use the rule tag as the reported word.
func (r *regRule) MatchAll(text string) map[string]string {
	var ret = make(map[string]string, 0)
	match, _ := r.reg.FindStringMatch(text)
	if match != nil {
		ret[r.Tag] = match.String()
	}
	return ret
}

// RegexpMatcher matches rules built from REGEXP_PREFIX-prefixed patterns.
//
// Build expects patterns without REGEXP_PREFIX because MatchService strips the
// prefix before dispatching to this matcher. Invalid patterns are ignored. This
// matcher reports matched regexp rule tags but does not perform text
// replacement.
type RegexpMatcher struct {
	matchers []*regRule
}

// RegexpMather is the old name for RegexpMatcher.
//
// Deprecated: use RegexpMatcher.
type RegexpMather = RegexpMatcher

// NewRegexpMatcher creates an empty RegexpMatcher.
func NewRegexpMatcher() *RegexpMatcher {
	return &RegexpMatcher{}
}

// Build compiles regexp patterns into matcher rules.
//
// Invalid regexp patterns are skipped to preserve the existing API, which does
// not return build errors.
func (a *RegexpMatcher) Build(words []string) {
	for _, w := range words {
		if m, err := NewRegRule(w); err == nil {
			a.matchers = append(a.matchers, m)
		}
	}
}

// Match returns regexp rule tags matched by text.
//
// onlyOne stops after the first regexp rule with any match. repl is accepted to
// satisfy Matcher but is not used because RegexpMatcher does not perform text
// replacement. The replacement text return value keeps the original text.
func (a *RegexpMatcher) Match(text string, onlyOne bool, repl rune) (word []string, desensitization string) {
	desensitization = text
	for _, r := range a.matchers {
		ret := r.MatchAll(text)
		for tag := range ret {
			word = append(word, tag)
		}
		if onlyOne && len(ret) > 0 {
			return
		}
	}
	return
}
