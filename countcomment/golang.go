package countcomment

type GOlang struct {
	cpp Cpp
}

func NewGolang() *GOlang {
	return &GOlang{
		cpp: *NewCpp(),
	}
}

func (p *GOlang) ScanLine(line string, isCommentLine bool, isCommentBlock bool, isString bool, isBreak bool) (
	numCommentLineRet int,
	numCommentBlockRet int,
	isCommentLineRet bool,
	isCommentBlockRet bool,
	isStringRet bool,
	isBreakRet bool,
) {

	numCommentLineRet, numCommentBlockRet, isCommentLineRet, isCommentBlockRet, isStringRet, isBreakRet = p.cpp.ScanLine(line, isCommentLine, isCommentBlock, isString, isBreak)
	return
}
