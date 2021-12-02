package countcomment

import (
	"strings"
	"testing"
)

var testDataScanLine = []struct {
	line               string
	isCommentLine      bool
	isCommentBlock     bool
	isString           bool
	isBreak            bool
	numCommentLineRet  int
	numCommentBlockRet int
	isCommentLineRet   bool
	isCommentBlockRet  bool
	isStringRet        bool
	isBreakRet         bool
}{
	// normal cases
	{`"instring" /*inblock*/ // Comment: inline`, false, false, false, false, 1, 1, false, false, false, false},
	{`normal `, false, false, false, false, 0, 0, false, false, false, false},
	{`/// comment inline `, false, false, false, false, 1, 0, false, false, false, false},
	{`/// comment inline \`, false, false, false, false, 1, 0, true, false, false, true},
	{`/* comment inblock `, false, false, false, false, 0, 1, false, true, false, false},
	{`" in string"`, false, false, false, false, 0, 0, false, false, false, false},
	{`" in string // line  /*block "`, false, false, false, false, 0, 0, false, false, false, false},
	{`" instring \ `, false, false, false, false, 0, 0, false, false, true, true},
	{`R"xxx(raw string `, false, false, false, false, 0, 0, false, false, true, false},
	{`R"({ "pr)"`, false, false, false, false, 0, 0, false, false, false, false},
	{`const char * vogon_poem = R"V0G0N("xxxxxxx")V0G0N"`, false, false, false, false, 0, 0, false, false, false, false},

	// special cases
	{`*/ int y = 1;  /*`, false, true, false, false, 0, 1, false, true, false, false},
	{`*/ int y = 1;  /*`, true, false, false, true, 1, 0, false, false, false, false},
	{`*/ int y = 1;  /* \`, true, false, false, true, 1, 0, true, false, false, true},
	{`" instring `, false, false, true, true, 0, 0, false, false, false, false},
	{`" instring \"  ss`, false, false, true, true, 0, 0, false, false, true, false},
	{`" instring \"  ss`, false, false, false, false, 0, 0, false, false, true, false},
	{`" instring \"  ss"`, false, false, false, false, 0, 0, false, false, false, false}, // end by "
	{`" instring \\"  ss`, false, false, false, false, 0, 0, false, false, false, false},
	{`" instring \\"  ss"`, false, false, false, false, 0, 0, false, false, true, false},
	{`a='"' `, false, false, false, false, 0, 0, false, false, false, false},
	{`a='\"' `, false, false, false, false, 0, 0, false, false, false, false},

	// mix case
	{`//* Comment: inline`, false, false, false, false, 1, 0, false, false, false, false},
	{`"// This is not a comment"`, false, false, false, false, 0, 0, false, false, false, false},
	{`"// This is not a comment\`, false, false, false, false, 0, 0, false, false, true, true},
	{` "// This is not a comment \" not a comment as well"`, false, false, false, false, 0, 0, false, false, false, false},
	{`"/* This is not a comment */"`, false, false, false, false, 0, 0, false, false, false, false},
	{`// This is not a comment" // Comment: inline`, false, false, false, false, 1, 0, false, false, false, false},
	{`*/ int x = 0;       // Comment: block, inline`, false, true, false, false, 1, 1, false, false, false, false},
	{`Comment: block */ // Comment: inline // Not another inline comment	`, false, true, false, false, 1, 1, false, false, false, false},
	{`"/* This is not a comment \`, false, false, false, false, 0, 0, false, false, true, true},
	{` );        /* Comment: block`, false, false, false, false, 0, 1, false, true, false, false},
}

func TestScanLine(t *testing.T) {
	eFile := NewCpp()
	leng := len(testDataScanLine)
	for i := 0; i < leng; i++ {
		e := testDataScanLine[i]
		numCommentLineRet, numCommentBlockRet, isCommentLineRet, isCommentBlockRet, isStringRet, isBreakRet := eFile.ScanLine(
			strings.TrimSpace(e.line),
			e.isCommentLine,
			e.isCommentBlock,
			e.isString,
			e.isBreak,
		)
		if numCommentLineRet != e.numCommentLineRet {
			t.Errorf("numCommentLineRet: index %d expected %d; got %d", i, e.numCommentLineRet, numCommentLineRet)
		}
		if numCommentBlockRet != e.numCommentBlockRet {
			t.Errorf("numCommentBlockRet: index %d expected %d; got %d", i, e.numCommentBlockRet, numCommentBlockRet)
		}
		if isCommentLineRet != e.isCommentLineRet {
			t.Errorf("isCommentLineRet: index %d expected %t; got %t", i, e.isCommentLineRet, isCommentLineRet)
		}
		if isCommentBlockRet != e.isCommentBlockRet {
			t.Errorf("isCommentBlockRet: index %d expected %t; got %t", i, e.isCommentBlockRet, isCommentBlockRet)
		}
		if isStringRet != e.isStringRet {
			t.Errorf("isStringRet: index %d expected %t; got %t", i, e.isStringRet, isStringRet)
		}
		if isBreakRet != e.isBreakRet {
			t.Errorf("isBreakRet: index %d expected %t; got %t", i, e.isBreakRet, isBreakRet)
		}
	}

}

var testDataScanFile = []struct {
	path            string
	numLines        int
	numCommentLine  int
	numCommentBlock int
}{
	{"testing/cpp/special_cases.cpp", 62, 6, 34},
	{"testing/cpp/test_lib_json/main.cpp", 3971, 182, 0},
	{"testing/c/json_reader.c", 1992, 134, 0},
	{"testing/c/fuzz.hpp", 14, 5, 0},
}
