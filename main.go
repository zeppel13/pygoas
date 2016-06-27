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
	VERSION       = "bugful20.05.2016.1"
	YEAR          = "2016"
	versionnumber = 0.1
	AUTHOR        = "Sebastian Kind"
	EMAIL         = "mail@sebastiankind.de"
	LICNECE       = "Yolo / NO WARRANTY; USE THIS PROGRAM ON YOUR OWN RISK"
)

// debug levels
// 0 - Show no Debug information
// 1 - still reserved
// 2 - Show LabelStack
// 3 - show other fancy stuff // Is this still useful?
var (
	debug = 0 // Helps, when you are alone
)

func main() {

	/*
		IDEA: -32bit: A 32 bit modus with eax instead rax
	*/
	var code programCode
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

	if len(tail) == 0 {
		fmt.Println("usage: compiler -FLAGS filename.py")
		os.Exit(1)
	} else if len(tail) >= 1 {

		files := make([]*os.File, 0)

		// all input file pointers are stored in files (slice). This
		// is done by the following for-loop

		for _, filename := range tail {
			file, err := os.Open(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
			files = append(files, file)
			defer file.Close()
		}

		// Translation: TOKENS ---¯¯_*MaGiC*_¯¯---> ASSEMBLY
		// For more information consuilt parser.go
		tokenList := make([]string, 1)

		// this loop translates the file slice to a long slice of tokens
		for _, file := range files {
			for _, fileTokens := range tokenize(file) {
				tokenList = append(tokenList, fileTokens)
			}

			if debug == 2 {
				// useful if nothing is compiled
				for _, v := range tokenList {
					fmt.Println(v)
				}
			}
		}

		// remember the programCode “Object“ code from the beginning of this file?
		code = translate(tokenList)

		outputString := fmt.Sprintf(code.code)
		if *versionBoolPtr {
			fmt.Println("This is", NAME, "written by", AUTHOR, "in", YEAR)
			fmt.Println("Build", VERSION)
			fmt.Println("Do you have an idea? Send a Mail to ", EMAIL)
		}

		// write Code to os.Stdout
		if *stdoutBoolPtr == true {
			*linkBoolPtr, *compileBoolPtr = false, false
			fmt.Printf("%v", outputString)
		}

		if *httpBoolPtr == true { // bzw. *httpBoolptr
			// first clojure of my life <3
			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				printHttpOutput(w, r, outputString)
			})

			if *portIntPtr != 8080 {
				portStr := ":" + strconv.Itoa(*portIntPtr)
				fmt.Println("Serving on Port:", portStr)

				log.Fatal(http.ListenAndServe(portStr, nil))
			} else {
				fmt.Println("Serving on Port 8080")
				log.Fatal(http.ListenAndServe(":8080", nil))
			}
		}

		// write Code to assembly file
		if *assemblyBoolPtr == true {
			outputFile, err := os.Create("./" + fileName + ".asm")

			if err != nil {
				fmt.Fprintf(os.Stderr, "Couldn't write "+fileName+
					".asm file %v\n", err)
				os.Exit(1)
			}
			outputFile.Write(([]byte)(outputString))
		}

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

func printHttpOutput(w http.ResponseWriter, r *http.Request, outputCode string) {
	fmt.Fprintf(w, "%s", outputCode)
}
