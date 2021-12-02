package countcomment

type CountCommentIF interface {
	ScanLine(line string, isCommentLine bool, isCommentBlock bool, isString bool, isBreak bool) (
		numCommentLineRet int,
		numCommentBlockRet int,
		isCommentLineRet bool,
		isCommentBlockRet bool,
		isStringRet bool,
		isBreakRet bool,
	)
}
