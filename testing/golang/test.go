package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"compass.com/go-homework/countcomment"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		printHelp()
	} else {
		dir := args[0]
		if err := CountCommentLines(dir); err != nil {
			fmt.Println(err)
		}
	}
}

func printHelp() {
	fmt.Println("usage: \n\tgo run . <directory>")
}

func CountCommentLines(dir string) error {
	// TODO: start your work here

	// find all require files
	paths := make([]string, 0, 1024)
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	}); err != nil {
		printError(err)
		return errors.New(fmt.Sprintf(`
		error:		not implemented.
		directory:	%s`, dir))
	}

	chstr := make(chan string, 10)
	strs := make([]string, 0, 10)
	var wg sync.WaitGroup
	// scanning files
ExitFor:
	for i := 0; i < len(paths); i++ {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			numLines, numCommentLine, numCommentBlock := ScanFile(path)
			chstr <- fmt.Sprintf("%-50s    total:%5d    inline:%5d    block:%5d\n", path, numLines, numCommentLine, numCommentBlock)
		}(paths[i])

		select {
		case str := <-chstr:
			strs = append(strs, str)
		case <-time.After(time.Second * 300):
			break ExitFor
		}
	}

	wg.Wait()

	close(chstr)

	sort.Strings(strs)

	// print results in console
	for i := 0; i < len(strs); i++ {
		fmt.Print(strs[i])
	}

	return nil
}

/* Open a file and read a line by ascending order. Print the counting numbers in the console. */
func ScanFile(path string) (numLines int, numCommentLine int, numCommentBlock int) {
	// specific file kind by os.Args when print cmmand in console, eg: go run . testing python
	var eFile countcomment.CountCommentIF
	if strings.HasSuffix(path, ".cpp") ||
		strings.HasSuffix(path, ".hpp") ||
		strings.HasSuffix(path, ".c") ||
		strings.HasSuffix(path, ".h") {

		eFile = countcomment.NewCpp()
		// }else if  strings.HasSuffix(path, os.Args[2]){
		// eFile = countcomment.NewPython()
	} else if strings.HasSuffix(path, ".go") {
		eFile = countcomment.NewGolang()
	} else {
		return
	}

	f, err := os.Open(path)
	if err != nil {
		printError(err)
		return
	}
	defer f.Close()

	// state of one line
	isCommentLine, isCommentBlock := false, false
	isString, isBreak := false, false

	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF { // end of a file
				break
			} else {
				printError(err)
			}
		}

		// analyse content of one line
		numCommentLineRet, numCommentBlockRet, isCommentLineRet, isCommentBlockRet, isStringRet, isBreakRet := eFile.ScanLine(strings.TrimSpace(line), isCommentLine, isCommentBlock, isString, isBreak)
		// save state
		isCommentLine, isCommentBlock, isString, isBreak = isCommentLineRet, isCommentBlockRet, isStringRet, isBreakRet
		// count numbers
		numLines++
		numCommentLine += numCommentLineRet
		numCommentBlock += numCommentBlockRet

		// debug
		// fmt.Printf("%d  %d  %d\n", numLines, numCommentLine, numCommentBlock)
	}

	return
}

func printError(err error) {
	fmt.Println(err.Error())
}
