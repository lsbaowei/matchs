package matchs

// DFAMatcher matches plain keywords with a trie-like DFA tree.
//
// Composite rules and regexp rules are handled by other matcher
// implementations. The repl argument is reserved for replacement behavior in
// this matcher; callers should use the returned match list as the authoritative
// match result.
type DFAMatcher struct {
	replaceChar rune
	root        *DFANode
	whiteRoot   *DFANode
}

// NewDFAMatcher creates an empty DFAMatcher.
func NewDFAMatcher() *DFAMatcher {
	return &DFAMatcher{
		root: &DFANode{
			End: false,
		},
	}
}

// NewDFAMather creates an empty DFAMatcher.
//
// Deprecated: use NewDFAMatcher.
func NewDFAMather() *DFAMatcher {
	return NewDFAMatcher()
}

// Build adds plain keywords to the DFA tree.
//
// Calling Build multiple times appends words to the current tree.
func (d *DFAMatcher) Build(words []string) {
	for _, item := range words {
		d.root.AddWord(item)
	}
}

// Match scans text for plain keywords.
//
// onlyOne asks the matcher to return after the first match. repl is the rune
// reserved for replacing matched ranges. The return values are the matched
// keywords and the replacement text reported by the matcher.
func (d *DFAMatcher) Match(text string, onlyOne bool, repl rune) (sensitiveWords []string, replaceText string) {
	if d.root == nil {
		replaceText = text
		return
	}

	textChars := []rune(text)
	//textCharsCopy := make([]rune, len(textChars))
	//copy(textCharsCopy, textChars)

	length := len(textChars)
	for i := 0; i < length; i++ {
		//root本身是没有key的，root的下面一个节点，才算是第一个；
		temp := d.root.FindChild(textChars[i])
		if temp == nil {
			continue
		}
		j := i + 1
		for ; j < length && temp != nil; j++ {
			if temp.End {
				sensitiveWords = append(sensitiveWords, string(textChars[i:j]))
				//replaceRune(textCharsCopy, repl, i, j)
			}
			temp = temp.FindChild(textChars[j])
		}

		if j == length && temp != nil && temp.End {
			sensitiveWords = append(sensitiveWords, string(textChars[i:length]))
			//replaceRune(textCharsCopy, repl, i, length)
		}
	}
	//replaceText = string(textCharsCopy)
	return
}

// replaceRune replaces chars in [begin, end) with replaceChar.
func replaceRune(chars []rune, replaceChar rune, begin int, end int) {
	for i := begin; i < end; i++ {
		chars[i] = replaceChar
	}
}

// DFANode is a node in the DFAMatcher keyword tree.
type DFANode struct {
	End  bool
	Next map[rune]*DFANode
}

// AddWord adds word into the subtree rooted at n.
func (n *DFANode) AddWord(word string) {
	node := n
	chars := []rune(word)
	for index := range chars {
		node = node.AddChild(chars[index])
	}
	node.End = true
}

// AddChild returns the child node for c, creating it when needed.
func (n *DFANode) AddChild(c rune) *DFANode {
	if n.Next == nil {
		n.Next = make(map[rune]*DFANode)
	}

	if next, ok := n.Next[c]; ok {
		return next
	}
	n.Next[c] = &DFANode{
		End:  false,
		Next: nil,
	}
	return n.Next[c]
}

// FindChild returns the child node for c, or nil when it does not exist.
func (n *DFANode) FindChild(c rune) *DFANode {
	if n.Next == nil {
		return nil
	}

	if _, ok := n.Next[c]; ok {
		return n.Next[c]
	}
	return nil
}
