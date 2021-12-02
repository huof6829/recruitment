package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleAdd(t *testing.T) {
	assert.Equal(t, 2, 1+1, "wrong calculation")
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

func TestScanFile(t *testing.T) {
	leng := len(testDataScanFile)
	for i := 0; i < leng; i++ {
		e := testDataScanFile[i]
		numLines, numCommentLine, numCommentBlock := ScanFile(e.path, "")

		if e.numLines != numLines {
			t.Errorf("numLines: index %d expected %d; got %d", i, e.numLines, numLines)
		}
		if e.numCommentLine != numCommentLine {
			t.Errorf("numCommentLine: index %d expected %d; got %d", i, e.numCommentLine, numCommentLine)
		}
		if e.numCommentBlock != numCommentBlock {
			t.Errorf("numCommentBlock: index %d expected %d; got %d", i, e.numCommentBlock, numCommentBlock)
		}

	}
}

var testDataCountCommentLines = []string{
	// "testing/",
	"testing/golang",
	// "testing/c",
	// "testing/cpp/",
	// "testing/cpp/lib_json",
}

func TestCountCommentLines(t *testing.T) {
	leng := len(testDataCountCommentLines)
	for i := 0; i < leng; i++ {
		CountCommentLines(testDataCountCommentLines[i])
	}
}
