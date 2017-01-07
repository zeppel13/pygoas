/* main.go
 * This file is part of pygoas
 * This (pygoas is probably going to be a simple Mini-Python to NASM-assembly compiler
 * written by Sebastian Kind 2nd Janury -- 20th May 2016
 */

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

const (
	NAME          = "pygoas"
	VERSION       = "bugful23.10.2016"
	YEAR          = "2016"
	versionnumber = "0.1"
	AUTHOR        = "Sebastian Kind"
	EMAIL         = "mail@sebastiankind.de"
	LICENCE       = "Yolo / NO WARRANTY; USE THIS PROGRAM AT YOUR OWN RISK" //#R?
)

// debug levels
// 0 - Show no Debug information
// 1 - still reserved
// 2 - Show LabelStack
// 3 - show other fancy stuff // Is this still useful?

var (
	debug = 0 // Helps, when you are alone

	// debug can't be a const due to the power of commandline flags to
	// change its value
)

func main() {

	// IDEA: -32bit: A 32 bit mode with eax instead rax. --> Nope
	// IDEA: 32bitArm(v6) mode? --> Nope

	var code programCode

	// code contains the final asm output stored in separate buffers
	// for every indentation level. It has lots of methods to change
	// its contents.

	// $0 -arguments
	compileBoolPtr := flag.Bool("c", true, "This runs the nasm compiler")
	linkBoolPtr := flag.Bool("l", true, "This runs the ld linker")
	assemblyBoolPtr := flag.Bool("s", true, "Compile only to NASM assembly")
	stdoutBoolPtr := flag.Bool("stdout", false, "Compile only to NASM assembly and print to stdout")
	versionBoolPtr := flag.Bool("v", false, "Show Versionnumber, Author, „Licence“")
	outFileStringPtr := flag.String("o", "prog", "Name of output file")

	// Silly stuff like a http-server :D
	httpBoolPtr := flag.Bool("http", false, "Enables http-server for code viewing")
	portIntPtr := flag.Int("port", 8080, "Used port for http-server")

	// turn this into an integer!
	debugBoolPtr := flag.Bool("debug", false, "Print debug information e.g. labelStack, cryptic numbers and other stuff. This option renders the output useless.")

	// get commandline arguments
	flag.Parse()
	fileName := *outFileStringPtr

	if *debugBoolPtr {
		debug = 2
	}

	// get input files
	tail := flag.Args()
	//fmt.Println("len", len(tail), "\n", tail)

	// The compiler may read many files at a time
	if len(tail) == 0 {
		fmt.Println("usage: compiler -FLAGS filename.py...")
		os.Exit(1)
	} else if len(tail) >= 1 {

		// #HowToUse Go?
		// reserve space for a slice, similar to Python's lists

		files := make([]*os.File, 0)

		// all file pointers are stored in files (slice). This
		// is done by the following for-loop

		for _, filename := range tail {
			file, err := os.Open(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1) // Error
			}
			files = append(files, file)
			defer file.Close() // Go's defer expressoin() releases an action at the end of a function
		}

		// Now tell me, how does this compiler work?
		// Translation: TOKENS ---¯¯_*MaGiC*_¯¯---> ASSEMBLY
		// For more information consult the source file parser.go

		tokenList := make([]string, 1)

		// this loop translates the file slice to a long slice of tokens
		// Files --> Tokens

		for _, file := range files {
			for _, fileTokens := range tokenize(file) {
				tokenList = append(tokenList, fileTokens)
			}

			if debug == 2 {
				// This is useful if compilation fails
				// most errors occure in a corrupt tokenList
				for _, v := range tokenList {
					fmt.Println(v)
				}
			}
		}

		// Do you remember the programCode “Object“ code from the beginning of code.go?
		// This is what it takes to convert tokens to some interpretation of assembly

		code = translate(tokenList)

		// outputstring represents the final assembly program
		outputString := fmt.Sprintf(code.code)

		// Now handle the diverse output formats of different flag
		// formats and their meanings.

		// print version legal blabla stuff
		// -v

		if *versionBoolPtr {
			fmt.Println("This is", NAME, "written by", AUTHOR, "in", YEAR)
			fmt.Println("Build", VERSION)
			fmt.Println("Do you have an idea? Send a Mail to ", EMAIL)
		}

		// write output to os.Stdout

		// this a easy way to check the results or find errors in the
		// assembly logic

		// --stdout
		if *stdoutBoolPtr == true {
			*linkBoolPtr, *compileBoolPtr = false, false
			fmt.Printf("%v", outputString)
		}

		// write the output to a http-server. This is maybe a useless
		// feature and somehow against the UNIX-philosophy, but I like
		// my decision. :-P Nonetheless all the net/http packages are
		// baked into my compiler, what increases the file size. To
		// not sacrify uploading speed over slow #FirstWorldProblems
		// Upload links, is surely possible to just ommit anything
		// what looks like http, networking or other neat l'internet
		// golang stuff. Don't forget the handler printHttpOutput(...)
		// at the end of this file.

		// read the code in your Browser with:
		// http://hostname_or_just_localhost:8080/anything

		//--http
		if *httpBoolPtr == true { // bzw. *httpBoolptr
			// first clojure of my life <3

			// Do I need a second call to printHttpOutout(...)
			// here. Can't I just do it inside the anonymous function?
			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				printHttpOutput(w, r, outputString)
			})

			// This sets the port
			if *portIntPtr != 8080 {
				portStr := ":" + strconv.Itoa(*portIntPtr)
				fmt.Println("Serving on Port:", portStr)

				log.Fatal(http.ListenAndServe(portStr, nil))
			} else {
				fmt.Println("Serving on Port 8080")
				log.Fatal(http.ListenAndServe(":8080", nil))
			}
		}

		// write Code to assembly file is good for debugging if NSAM wont
		// do its job. Then make changes with vim directly in the
		// assembly file compile with NASM by Hand using these commands:
		// 		nasm -felf64 -o outfile.o ./infile.asm
		//		ld -o program ./outfile.o
		//  Also consider the --stdout flag of pygoas

		// -s
		if *assemblyBoolPtr == true {
			outputFile, err := os.Create("./" + fileName + ".asm")

			if err != nil {
				fmt.Fprintf(os.Stderr, "Couldn't write "+fileName+
					".asm file %v\n", err)
				os.Exit(1)
			}
			outputFile.Write(([]byte)(outputString))
		}

		// Compile and link the output of pygoas with NASM and LD
		// invoke NASM to compile file
		if *compileBoolPtr == true {
			// run nasm: nasm -f elf64 ./prog.asm
			nasmPath, err := exec.LookPath("nasm")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Sorry, I couldn't find an installation of NASM. This is what I got: %v\n", err)
				os.Exit(1)
			}
			nasmArgs := []string{"-f", "elf64", "./" + fileName + ".asm"}
			err = exec.Command(nasmPath, nasmArgs...).Run()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Dude, something went wrong starting NASM. This is what I got: %v\n", err)
			}
		}

		// invoke LD to link objectfile made by NASM
		if *linkBoolPtr == true {
			// run ld: ld ./prog.o
			ldPath, err := exec.LookPath("ld")
			if err != nil {
				fmt.Fprintf(os.Stderr, "I couldn't find an installation of LD. This is what I got: %v\n", err)
				os.Exit(1)
			}

			ldArgs := []string{"-o", "./" + fileName, "./" + fileName + ".o"}
			err = exec.Command(ldPath, ldArgs...).Run()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Something went wrong starting LD. This is what I got: %v\n", err)
				os.Exit(1)
			}
		}
	}

}

// Magic HTTP-Server inside a compiler
func printHttpOutput(w http.ResponseWriter, r *http.Request, outputCode string) {
	fmt.Fprintf(w, "%s", outputCode)
}
