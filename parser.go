/* parser.go
 * This file is part of pygoas
 *
 * This is probably going to be a simple Mini-Python to NASM-assembly compiler
 * written by Sebastian Kind 2nd Janury -- 20th May 2016
 */

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// tokenize (file *os.File) [] string reads a file return its content
// in a []string slice. Each element contains a word, a token. Special
// charakters like '\n', '\t' are replaced with more readable words.

// File should be red before tokenize(...) : []string
// ideally tokenize should get
func tokenize(file *os.File) []string {
	rawBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	rawCode := string(rawBytes)

	// real men indent with 4 spaces
	rawCode = strings.Replace(rawCode, "    ", " INDENT ", -1)
	rawCode = strings.Replace(rawCode, "\t", " INDENT ", -1)
	rawCode = strings.Replace(rawCode, "\n", " NEWLINE ", -1)
	rawCode += " EOF " // it's not the real EOF value

	tokens := strings.Fields(rawCode)
	debugFlag := false //Moment; What? FIXME // HI THIS IS JAN: „hOLLA, AMIGOS. DOSN!T WORK. MUCH NOOB“ AND PIZZA111 HANK I C U ON FRIDAY“
	if debugFlag == true {
		for _, v := range tokens {
			fmt.Println(v)
		}
	}

	return tokens
}

// translate(tokens []string) : programCode This functions accepts a
// []string slice, filled with valid tokens. They'll be translated to
// a programCode struct, which will be returned.
// programCode is a struct of code.go

// This function manages the main work of this program. It controlls a
// central programCode instance and filles it with content

func translate(tokens []string) programCode {
	code := newProgramCode()
	code.createStart() // assembly header must be created

	// slices of protected words and sign follow here:

	keywords := []string{
		"for", "while", "def", "if", "return", "#endif", "#endfor",
	}

	mathOperators := []string{
		"=",
		"+", "-", "*", "/", "%", "**",
		"+=", "-=", "*=", "/=", "%=",
	}

	special := []string{
		"NEWLINE",
		"INDENT",
		"EOF",
		",", ".", ";", ":",
		"(", ")",
		"#", "'''",
		"\"", "(\"", "\")",
		"^", "^^", ":)",
	}

	// still todo
	// probably requieres major change in the value passing structure
	// »Mehr Register wagen.« :)
	//
	// boolOperators := []string{"and", "or", "xor", "not"} the result
	// of an operator should be returned and stored into rax

	cmpOperators := []string{
		"==", "!=", "<", ">", "<=", ">="}
	builtInFunc := []string{
		"print", "return", "asm"}

ParserLoop:
	for i := 0; i < len(tokens); i++ {
		if debug == 2 {
			fmt.Println(code.indentLevel)
		}

		v := tokens[i]
		if v == "#" {
			code.commentFlag = true
			continue ParserLoop

		} else if v == "NEWLINE" && code.commentFlag {
			code.commentFlag = false
		} else if code.commentFlag {
			continue
		} else if v == "=" { // variable assignment with *easy math
			if stringInSlice(tokens[i+2], mathOperators) {
				code.addVar(tokens[i-1], 0)
				switch tokens[i+2] {
				case "+":
					code.createAdd(tokens[i+1], tokens[i+3], tokens[i-1])
				case "-":
					code.createSub(tokens[i+1], tokens[i+3], tokens[i-1])
				case "*":
					code.createMul(tokens[i+1], tokens[i+3], tokens[i-1])
				case "/":
					code.createDiv(tokens[i+1], tokens[i+3], tokens[i-1], "")
				case "%":
					code.createDiv(tokens[i+1], tokens[i+3], "", tokens[i-1])
				}

			} else if v, err := strconv.Atoi(tokens[i+1]); err == nil { // variable assignment
				code.addVar(tokens[i-1], (int64)(v))
				code.initVar(tokens[i-1], tokens[i+1])
				if debug == 3 {
					fmt.Fprintf(os.Stderr, "neue Variable\n")
				}
			}
			// token [i] should be '='
			// token:  a   =   "hi"
			//   [i]: -1   0   1
			if strings.HasPrefix(tokens[i+1], "\"") && strings.HasSuffix(tokens[i+1], "\"") {
				name := tokens[i-1]
				code.strCounter++
				code.stringMap[name] = tokens[i+1]
			}

		} else if strings.HasSuffix(v, "=") && v != "=" {
			// every abbreviation ends with '='
			// If you don't believe me: "+=", "-=", "*=", "%="
			switch v {
			case "+=":
				code.createAdd(tokens[i-1], tokens[i+1], tokens[i-1])
			case "-=":
				code.createSub(tokens[i-1], tokens[i+1], tokens[i-1])
			case "*=":
				code.createMul(tokens[i-1], tokens[i+1], tokens[i-1])
			case "/=":
				code.createDiv(tokens[i-1], tokens[i+1], tokens[i-1], "")
			case "%=":
				code.createDiv(tokens[i-1], tokens[i+1], "", tokens[i-1])

			}

		} else if v == "print" { // print is a builtin function, which requieres special magic
			argSlice := getArgSlice(i+1, tokens)
			for _, arg := range argSlice {
				if strings.HasPrefix(arg, "\"") && strings.HasSuffix(arg, "\"") {
					name := "strConst" + strconv.FormatInt(code.strCounter, 10)
					code.strCounter++
					code.stringMap[name] = arg
					code.createPrintString(name)
				} else {
					code.createPrint(arg)
				}

			}
		} else if v == "asm" { // Place nasm assembly code inlinely
			// This is the hardcore-brechstangen-way-of-life-methode:
			argString := ""
			for !strings.HasSuffix(tokens[i], "\")") {
				i++
				argString += tokens[i] + " "
			}
			argString = strings.TrimPrefix(argString, "(\"")
			argString = strings.TrimSuffix(argString, "\") ")

			code.asm(argString)
		} else if v == "def" { // "def" is the keyword to define "functions"
			code.labelFlag = true
			if debug == 3 { // debug fun
				fmt.Fprintf(os.Stderr, "def detected\n")
			}
			code.createLabel(tokens[i+1])
			// some random thoughts
			// + safe localvalues to .bss ? like funcFoo1Locals : resb 32 ; enough space for 4 64bit values
			// copy args into global variables (bad style -- but the only)
			// return max 2 vars: rax, rbx
		} else if v == "return" {
			// necessarry to end a function
			code.createReturn()

		} else if v == "if" {
			// Conditions work fine
			tokens[i+3] = strings.TrimSuffix(tokens[i+3], ":")
			code.labelCounter++
			strLabelCounter := strconv.FormatInt(code.labelCounter, 10)
			labelName := "ifLabel" + strLabelCounter

			code.createCmp(tokens[i+1], tokens[i+3])

			switch tokens[i+2] {
			case "==":
				code.isEqual(labelName)
			case "<":
				code.isSmaller(labelName)
			case ">":
				code.isGreater(labelName)
			case "<=":
				code.isSmallerEqual(labelName)
			case ">=":
				code.isGreaterEqual(labelName)
			case "!=":
				/* Do something */
			}

			code.createJumpBackLabel("ifReturn")
			code.createLabel(labelName)

		} else if v == "while" {
			// use a for-loop!

		} else if v == "for" {

			tokens[i+1] = strings.TrimRight(tokens[i+1], ":")
			tokens[i+1] = strings.TrimRight(tokens[i+1], ")")
			tokens[i+1] = strings.TrimLeft(tokens[i+1], "(")
			loopVarVal, _ := strconv.Atoi(tokens[i+1])
			loopVarVal++                                             // No clue why 'loopVarVal++'
			strLoopVar := strconv.FormatInt((int64)(loopVarVal), 10) // Why not: tokens[i+1] ??? :D

			code.labelCounter++
			strLabelCounter := strconv.FormatInt(code.labelCounter, 10)

			strLabelName := "forLabel" + strLabelCounter
			code.loopVarCounter++
			strCounter := strconv.FormatInt(code.loopVarCounter, 10)
			strLoopVarName := "loopVar" + strCounter

			code.createJump(strLabelName)

			code.pushLastLabel(strLabelName)
			code.createJumpBackLabel("forReturn")
			code.createResetLoopVar(strLoopVarName, loopVarVal)
			code.createLabel(strLabelName)
			code.addVar(strLoopVarName, (int64)(loopVarVal))
			code.initVar(strLoopVarName, strLoopVar)
			code.createForCheck(strLoopVarName)
		} else if v == "#endif" {
			// necessary to end a If-body
			code.createJumpBack()
			//code.labelFlag = false
		} else if v == "#endfor" {
			// necessarry to end a for-body
			code.createJumpBack()

		} else if i > 1 && i+1 < len(tokens) {
			// This lets the final progam to execute a function.
			// this allows to run functions inside (mini)Python

			// How dows this work?  If v doesn't fit into any
			// description and is followed by a opening parenthesis,
			// v will be a function call
			// if functionCall -> create Code to run a function

			_, err := strconv.ParseInt(v, 10, 64)
			if tokens[i-1] != "def" &&
				!stringInSlice(v, mathOperators) &&
				!stringInSlice(v, cmpOperators) &&
				!stringInSlice(v, builtInFunc) &&
				!stringInSlice(v, keywords) &&
				!stringInSlice(v, special) &&
				strings.HasPrefix(tokens[i+1], "(") &&
				!strings.HasPrefix(tokens[i+1], "\"") &&
				err != nil {
				code.createCall(v) //calls the function with the name stored in v
			}

		} else if v == "INDENT" {
			/* Python's are totally ignored */
			// this will come someday in the Futute(
		} else if v == "NEWLINE" {
			continue ParserLoop
			// When the line is over, the next command/line will be translated
		} else if v == "EOF" {
			// A fake EOF (not the real one), must be the last token
			// in []tokens, in order to finishe the program

			code.createExit("0") // Return value.
			// Read this value with
			// user@host ~% $? on your terminal
			// after the program has quit

			// closing work
			code.createAllFunctions()
			code.createBss()
			code.createData()
			return code
		}

	}
	return code // struct programCode
}

func stringInSlice(s string, list []string) bool {
	for _, v := range list {
		if s == v {
			return true
		}
	}
	return false
}

// this functions returns a list or arguments from an argument statement in the minipython code
// for example:
// print (a, 12, "hallo")
// return sclice of type []string holding {"Hallowelt", "3", "4", "hallo"}

// getArgSlice (...) makes a slice of strings from a given argument list

// getArgSlice(...) must know the tokenlist and the actual positon inside it
func getArgSlice(i int, tokens []string) []string {
	oldi := i
	argString := ""
	var argSlice []string
	for { // "for {...}" is go's "while true {...}" or "for (;;) {...}"
		if tokens[i] == "NEWLINE" {
			break
		}

		if (strings.HasSuffix(tokens[i], ")")) == false || (tokens[i] == ")") == false {
			if oldi == i {
				argString += tokens[i]
			} else {
				argString += " " + tokens[i]
			}
		}

		if strings.HasPrefix(tokens[i], "\"") || strings.HasSuffix(tokens[i], "\"") {
			// Currently doing nothing
			// Was this once useful?
		}
		i++
	}
	//this removes parenthesis
	argString = strings.TrimRight(argString, ")")
	argString = strings.TrimLeft(argString, "(")
	argSlice = strings.Split(argString, ", ")
	if debug == 3 {
		fmt.Fprintln(os.Stderr, argSlice)
	}
	return argSlice
}
