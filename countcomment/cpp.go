package countcomment

type Cpp struct {
}

func NewCpp() *Cpp {
	return &Cpp{}
}

// overwrite IF

/* Scan a line and analyse the special character. Return the states of a line and counting numbers. The rule is
content                 isCommentLineRet  isCommentBlockRet  isStringRet isBreakRet
normal                  false             false              false        false
/// comment inline      false             false              false        false
/// comment inline \    true              false              false        true
/* comment inblock      false             true               false        false
/* comment inblock*\/   false             false              false        false
" in string in "        false             false              false        false
" instring //  /*       false             false              true         false
" instring \            false             false              true         true
R"xxx(raw string        false             false              true         false
*/
func (p *Cpp) ScanLine(line string, isCommentLine bool, isCommentBlock bool, isString bool, isBreak bool) (
	numCommentLineRet int,
	numCommentBlockRet int,
	isCommentLineRet bool,
	isCommentBlockRet bool,
	isStringRet bool,
	isBreakRet bool,
) {
	// debug
	// fmt.Printf("%t %t %t %t %s\n", isCommentLine, isCommentBlock, isString, isBreak, line)

	// state of pre-comment
	isCommentBegin := isCommentLine || isCommentBlock
	// state of pre-finishing for commentblock
	isCommentBlockEnd := false
	// length of one line
	leng := len(line)

	/* last line is a comment line and the ending has a break,  eg:
	"// This line has a trailing comment" // Comment: \
	inline (the line-break makes this one counted as 2 lines)
	*/
	if isCommentLine && isBreak {
		numCommentLineRet++
		if leng > 0 && line[leng-1] == '\\' {
			isBreakRet = true
		} else {
			isCommentLine = false
		}
		isCommentLineRet = isCommentLine
		return
	}

	// mark the duplicating lines.  eg:   */ int y = 1;  /*
	isCountedLine := false
	if isCommentBlock {
		numCommentBlockRet++
		isCountedLine = true
	}

	// occurences of the break, eg: "\\\\\"
	numBreak := 0
	indexLastBreak := 0

	// raw string, eg: R"xxx(yyy)xxx"
	isRawString := false
	isRawHead := false
	rowHead := make(map[string]bool)
	indexLastRawQuato := 0

	/* scan from the beginning and detect the target character including: / * " \ ( ) , then analyse the current state of line is
	in string or in comment-line or in comment-block, then set the right states and return states for analysing the next line
	*/
	for i := 0; i < leng; i++ {

		if isCommentLine {
			break
		}

		// current character is / and not in string
		if line[i] == '/' && !isString {
			if !isCommentBegin {
				isCommentBegin = true
			} else {
				if !isCommentLine && !isCommentBlock && i > 0 && line[i-1] == '/' {
					isCommentLine = true
					numCommentLineRet++
				}
				if isCommentBlockEnd && i > 0 && line[i-1] == '*' {
					isCommentBlock = false
					isCommentBegin = false
				}
			}

			// current character is * and not in string
		} else if line[i] == '*' && !isString {
			if isCommentBegin {
				if !isCommentBlock && i > 0 && line[i-1] == '/' {
					isCommentBlock = true
					isCommentBlockEnd = false

					// remove duplicating lines.  eg:   */ int y = 1;  /*
					if !isCountedLine {
						isCountedLine = true
						numCommentBlockRet++
					}
				}

				if isCommentBlock && !isCommentBlockEnd {
					isCommentBlockEnd = true
				}
			}
			// current character is " and not in comment
		} else if line[i] == '"' && !isCommentLine && !isCommentBlock {
			if !isString {
				// escape character, eg: a='"' , b='\"'
				if i > 0 && line[i-1] == '\'' && i < (leng-1) && line[i+1] == '\'' {
					isString = false
				} else if i > 1 && line[i-2] == '\'' && line[i-1] == '\\' && i < (leng-1) && line[i+1] == '\'' {
					isString = false
				} else {
					// raw string in one line, eg:   R"({ "pr)"
					if i > 0 && line[i-1] == 'R' {
						isRawString = true
						indexLastRawQuato = i
						isRawHead = true
					}
					isString = true
				}
			} else {
				/// analyse raw-string
				if isRawString {
					s := line[indexLastRawQuato : i+1] //  xxx"
					str := s[1:]
					if len(str) > 1 {
						str = str[len(str)-1:] + str[:len(str)-1]
					}
					if v, ok := rowHead[str]; !ok || !v {
						continue
					}

					isRawString = false
					rowHead[str] = false
					isString = false
				}

				/* special cases, eg:
					"xxxxx \
				    "yyyyy
				*/
				if isBreak { // the end of last line is \
					isString = false
				}

				// analyse breaks
				if (i - indexLastBreak) == 1 {
					if numBreak%2 == 1 { //  eg: "\\\\"
						isString = true
					} else {
						isString = false //  eg: "\"”, "\\\""
						numBreak = 0
					}
				} else {
					isString = false
					numBreak = 0
				}
			}

			// current character is \ and in string
		} else if line[i] == '\\' && isString { // count breaks
			if numBreak == 0 { // first occur
				numBreak++
			} else if (i - indexLastBreak) == 1 { // continuous occur
				numBreak++
			}
			indexLastBreak = i

			// current character is ( and in raw-string
		} else if line[i] == '(' && isRawHead { // get "xxx from R"xxx(yyy)xxx"
			rowHead[line[indexLastRawQuato:i]] = true
			isRawHead = false

			// current character is ) and in raw-string
		} else if line[i] == ')' && isRawString { // get index of ) to analyse the ending of raw-string
			indexLastRawQuato = i
		}

	}

	// check the ending character
	if leng > 0 && line[leng-1] == '\\' {
		isBreak = true
	} else {
		isBreak = false
		if isCommentLine {
			isCommentLine = false
		}
	}

	// save all states
	isCommentLineRet, isCommentBlockRet, isStringRet, isBreakRet = isCommentLine, isCommentBlock, isString, isBreak
	return

}
