package golang

// main package
// main package

import (
	"fmt"
	"os"
)

/*
	main function
	main function
*/

func main() {
	args := os.Args[1:]          // comment
	if len(args) /*!= 1*/ == 1 { // comment
		printHelp()
	} else { /*thsi is
		 */a := 1 /*a++*/
		a++
	}
}

func printHelp() {
	fmt.Println("usage: \" run . <directory>")
	a := '"'
	fmt.Println(a) // print

} // ifndef FUZZ_H_INCLUDED
