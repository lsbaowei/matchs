package matchs

import (
	"reflect"
	"testing"
)

var (
	_ Matcher = (*DFAMatcher)(nil)
	_ Matcher = (*AssembleMatcher)(nil)
	_ Matcher = (*AssembleMather)(nil)
	_ Matcher = (*RegexpMatcher)(nil)
	_ Matcher = (*RegexpMather)(nil)
)

func TestDFAMatcherConstructorsMatchSame(t *testing.T) {
	newMatcher := NewDFAMatcher()
	oldMatcher := NewDFAMather()

	for _, matcher := range []*DFAMatcher{newMatcher, oldMatcher} {
		matcher.Build([]string{"keyword"})
	}

	assertSameWords(t, newMatcher, oldMatcher, "has keyword")
}

func TestAssembleMatcherConstructorsMatchSame(t *testing.T) {
	newMatcher := NewAssembleMatcher()
	oldMatcher := NewAssembleMather()

	for _, matcher := range []*AssembleMatcher{newMatcher, oldMatcher} {
		matcher.Build([]string{"A#B"})
	}

	assertSameWords(t, newMatcher, oldMatcher, "xxAxx")
}

func TestRegexpMatcherConstructorsMatchSame(t *testing.T) {
	newMatcher := NewRegexpMatcher()
	oldMatcher := &RegexpMather{}

	for _, matcher := range []*RegexpMatcher{newMatcher, oldMatcher} {
		matcher.Build([]string{`1\d{10}`})
	}

	assertSameWords(t, newMatcher, oldMatcher, "phone 13800138000")
}

func assertSameWords(t *testing.T, left Matcher, right Matcher, text string) {
	t.Helper()

	leftWords, _ := left.Match(text, false, '*')
	rightWords, _ := right.Match(text, false, '*')

	if !reflect.DeepEqual(leftWords, rightWords) {
		t.Fatalf("words mismatch: left=%#v right=%#v", leftWords, rightWords)
	}
}
