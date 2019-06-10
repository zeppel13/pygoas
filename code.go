/* code.go
 * This file is part of pygoas
 *
 * This is probably going to be a simple Mini-Python to NASM-assembly compiler
 * written by Sebastian Kind 2nd Janury -- 20th May 2016
 */

package main

import (
	"fmt"
	"strconv"
)

// programCode is the essential struct (comparable to a class in
// python) which contains the logic necessary to compile code for
// example variables, labels, indents and espacially the assembly
// source code, which was written by the compiler.

// A Struct is similar to a class in Python. Variables and other
// Valueholding are declared inside of a struct.

// Go -- Python
// struct -- class
// slice -- list
// map -- dictionary

type programCode struct {
	intMap     map[string]int64
	stringMap  map[string]string
	strCounter int64 // 64 because pygoas is a 64bit compiler :P

	commentFlag    bool
	labelFlag      bool
	labelCounter   int64
	forLoopFlag    bool
	loopVarCounter int64

	funcSlice   []string
	lastLabel   []string
	indentLevel int

	code     string
	funcCode []string
}

// constructor of programCode

// Go wants contructors look the following function.

func newProgramCode() programCode {
	var code programCode
	code.intMap = make(map[string]int64)
	code.stringMap = make(map[string]string)
	code.strCounter = 0
	code.funcSlice = make([]string, 100) // TODO: Make this more dynamic e.g.: zero-length slice with dynamic append()
	code.lastLabel = make([]string, 100) // TODO: Make this more dynamic e.g.: zero-length slice with dynamic append()
	code.labelCounter = 0

	code.indentLevel = 0
	code.labelFlag = false
	code.forLoopFlag = false
	code.code = ""
	code.funcCode = make([]string, 100) // TODO: Make this more dynamic e.g.: zero-length slice with dynamic append()

	return code
}

// This method appends the raw nasm assmebly code to the output
// program.

func (pc *programCode) appendCode(code string) {
	pc.funcCode[pc.indentLevel] += code // Append Code to processed indent level

}

// This method adds new Variables to the compiler logic and to the
// output program. Variables are unsigned 64 Bit wide integers
// 0...2**64 (python-power)

func (pc *programCode) addVar(name string, val int64) {
	// the following line checks if an element exists inside the map element.
	if _, ok := pc.intMap[name]; ok {
		pc.setVar(name, val)
	} else {
		pc.intMap[name] = val
	}

}

// This method sets a Variable to a known value while the compiled binary is running.
func (pc *programCode) setVar(name string, val int64) {
	code := ""
	strVal := strconv.FormatInt(val, 10)
	code += "\tmov rax, " + strVal + "\t;set " + name + " to " + strVal + "\n"
	code += "\tmov [" + name + "], rax \n"
	pc.appendCode(code)
}

// The builtin print code will be created by the compiler everytime
// print is called inside the (mini)python program. It accepts somehow
// variadic parameters.

// usage of print: print ("text", variable, variable, "text", ...)

/*
;print:
	mov rax, 1		;syscall: write print value
	mov rdi, 1 ; stdout is the 'output file'
	mov rsi, value ; ptr
	mov rdx, 1 ; len
	syscall
	;ret
*/

// createPrint only prints one letter/char/byte at the time. The
// stringlength in |rdx| has the value 1.
// FIXME: Is this funtion in use?

func (pc *programCode) createPrint(s string) {
	print := "\tmov rax, 1\t;print " + s + "\n\tmov rdi, 1\n\tmov rdx, 1\n\tmov rsi, " + s + "\n\tsyscall\n"
	pc.appendCode(print)
}

// createPrintString is able to print a whole string.
func (pc *programCode) createPrintString(sname string) {
	len := (int64)(len(pc.stringMap[sname])) - 2

	// FIXME: WTF int64. Why not use int and strconv.Atoi(var string)
	// and stringcon.Itoa(var int)

	strlen := strconv.FormatInt(len, 10)
	code := "\tmov rax, 1\t;print String" + sname + "\n\tmov rdi, 1\n\tmov rdx, " + strlen + "\n\tmov rsi, " + sname + "\n\tsyscall\n"
	pc.appendCode(code)
}

/*
mov eax, [var1]
add eax, [var2]
mov [var3], eax
*/

// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+//

// Mathsnippets The following code templates are appended to the
// output program to do ugly numbercrunshing work The following
// functions' parameters are names of variables inside the output
// program

// createAdd("AnzahlMurmeln", "MurmelSack3000", "AnzahlMurmeln") ?

// Addition
func (pc *programCode) createAdd(a, b, sum string) {
	code := "\n\t\t\t; Add " + b + " to " +
		a + " and save sum in " + sum + "\n"

	if _, err := strconv.Atoi(a); err == nil {
		code += "\tmov rax, " + a + "\n"
	} else {
		code += "\tmov rax, [" + a + "]\n"
	}
	if _, err := strconv.Atoi(b); err == nil {
		code += "\tadd rax, " + b + "\n"
	} else {
		code += "\tadd rax, [" + b + "]\n"
	}
	code += "\tmov [" + sum + "], rax\n"

	pc.appendCode(code)
}

// Subtraction
func (pc *programCode) createSub(m, s, dif string) {
	code := "\n\t\t\t; Substract " + s + " from " +
		m + " and save difference in " + dif + "\n"

	if _, err := strconv.Atoi(m); err == nil {
		code += "\tmov rax, " + m + "\n"
	} else {
		code += "\tmov rax, [" + m + "]\n"
	}

	if _, err := strconv.Atoi(s); err == nil {
		code += "\tsub rax, " + s + "\n"
	} else {
		code += "\tsub rax, [" + s + "]\n"
	}
	code += "\tmov [" + dif + "], rax\n"
	pc.appendCode(code)
}

// Multiplication
func (pc *programCode) createMul(a, b, prod string) {
	code := "\n\t\t\t; Multiply " + a + " with " +
		b + " and store product in " + prod + "\n"

	if _, err := strconv.Atoi(a); err == nil {
		code += "\tmov rax, " + a + "\n"
	} else {
		code += "\tmov rax, [" + a + "]\n"
	}
	if _, err := strconv.Atoi(b); err == nil {
		code += "\timul rax, " + b + "\n"
	} else {
		code += "\timul rax, [" + b + "]\n"
	}
	code += "\tmov [" + prod + "], rax\n"
	pc.appendCode(code)
}

// Division

/*
	mov rax, [divisor]		;divides rax by rbx remainder is stored in rdx quotient is stored in rax
	mov rbx, [dividend]
	div rbx
	mov [q], rax ;; quotient
	mov [r], rdx ;; remainder
*/

// Make shure to not divide by zero. It'll cause a floting point error
// and program will crash. This feature is still buggy.

func (pc *programCode) createDiv(divisor, dividend, quotient, remainder string) {
	divcode := "\n\t\t\t; Divide " + divisor + " by " +
		dividend + " and safe quotient in " + quotient + "\n"

	divcode += "\t\t\t; Safe remainder in " + remainder + "\n"
	if _, err := strconv.Atoi(divisor); err == nil {
		divcode += "\tmov rax, " + divisor + "\n"
	} else {
		divcode += "\tmov rax, [" + divisor + "]\n"
	}

	if _, err := strconv.Atoi(dividend); err == nil {
		divcode += "\tmov rbx, " + dividend + "\n"
	} else {
		divcode += "\tmov rbx, [" + dividend + "]\n"
	}

	divcode += "\tdiv rbx\n"

	if quotient != "" {
		divcode += "\tmov [" + quotient + "], rax\n"
	}
	if remainder != "" {
		divcode += "\tmov [" + remainder + "], rdx\n"
	}
	pc.appendCode(divcode)
}

// End of Math
//*/*/*/*/*/*/*/*/*/*/*/*/*/*/*/*/*/*/*/*/*/*/*/*/*/*/*/*/*/*/*/*/*/*/*/*/*/*//*//

// createParams the code to copy values into argument registers
// as defined in the *amd64 System V calling convention*

// Marginal Note: MiniPython functions can'take any arguments still not a supported feature*
// FIXME: Is this in use?

// As long as argSlice delivers a value, it'll be thrown in one of the
// following registers as ordered in registerSlice. On this way six
// 64bit numbers may be passed over to the called function, which can
// easily read from the registers.

func (pc *programCode) createParams(argSlice []string) {
	code := ""
	registerSlice := []string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"} // SysV ABI calling register for parameters
	for i := 0; i < len(argSlice) && i < 6; i++ {
		if _, err := strconv.Atoi(argSlice[i]); err == nil {
			code += "\tmov " + registerSlice[i] + argSlice[i] + "\n"
		} else {
			code += "\tmov " + registerSlice[i] + "[" + argSlice[i] + "]\n"
		}
	}
	pc.appendCode(code)
}

// createCall allows the label passed as string argument to be
// called.
//"call" executes a function inside assembly. It cleanes used
// registers before and after the function did its job. -> amd64 Sys V
// abi
// FIXME: 'jmp vs. call'?

func (pc *programCode) createCall(name string) {
	code := ""
	code += "\n\tcall " + name + "\t; call label " + name + "\n"
	pc.appendCode(code)
}

// crateLabel marks a label inside the assembly source code. It also
// increments the indentLevel counter, in order to write the following
// code block into a separated buffer. Labels can be called or jumped
// to. createLabel accepts the label's name as a argument

func (pc *programCode) createLabel(name string) {
	code := ""
	code += "\n" + name + ":\n"
	pc.funcSlice = append(pc.funcSlice, name)
	pc.indentLevel += 1 // dive deeper -> next buffer.
	// Please have a look to FIXME: Where can I find what?
	pc.appendCode(code)

}

// createReturn leaves the innermost function/indent buffer by
// decrementing the pc.indentLevel. It is the last function appending
// information to a asm-Code block.

func (pc *programCode) createReturn() {
	code := "\tret\n"

	pc.appendCode(code)
	pc.indentLevel-- //  get back -> buffer before
}

// createJump allows the final program to jump to a label. This is
// used for functions. FIXME: Rest(for, if)?

func (pc *programCode) createJump(label string) {
	code := ""
	code += "\tjmp " + label + "\t; jmp to " + label + "\n"
	pc.appendCode(code)
}

// createJumpBackLabel writes a label to the main code to which the
// final program can jump back after a functions, if-clause or
// for-loop was finished

// Interesting function call:
// pc.pushLastLabel(label) places the label (func, if, for, etc) onto
// a stack memory, in order to remind the program where it should jump
// next

func (pc *programCode) createJumpBackLabel(category string) {
	code := ""
	strlabelCounter := strconv.FormatInt(pc.labelCounter, 10)
	label := category + strlabelCounter
	pc.pushLastLabel(label)
	code += "\t" + label + ":\t; return point\n"
	pc.appendCode(code)

}

func (pc *programCode) createJumpBack() {
	code := ""
	label := pc.popLastLabel()
	code += "\tjmp " + label + "\t; return to last place\n"
	pc.appendCode(code)
	pc.indentLevel--
}

// createResetLoopVar appends a code snippet to pc.code which resets a
// loopVarN to a given value.
// Is this funtions necessary? Why not use programCode.SetVar(int64, string)?
func (pc *programCode) createResetLoopVar(name string, val int) {
	valStr := strconv.Itoa(val)
	code := ""
	code += "\tmov rax, " + valStr + "\t;reset LoopVar to" + valStr + "\n"
	code += "\t mov [" + name + "], rax;\t done\n"
	pc.appendCode(code)
}

// The compiler has a stack to manage nested functions, conditions and
// loops. It is still a so called Brechstangen-Methode due to the
// inflexibility of Go's slices compared to Python's lists. Slices
// refer to an underlying array of something. They are basically a pointer
// to the real chunk of date used, which has some dynamic aspects.

// pc.LastLabel[n] represents the postion of a label in the hierarchy
// of a running program.

// A generic funtion to let the stack grow and shrink is indispensable
// for a MiniPython program which consists of a lot branching like
// conditions, loops, functions. The sad trueth is that a limited
// brechstangen-code sets the borders of a MiniPython system.
// Branching should work good enough with eight stack layers.

func (pc *programCode) pushLastLabel(name string) {

	// errors happend often enough to place some debug logic here. The
	// really ugly and terminal filling printed debug messages should
	// mainly show the changes made to the stack.

	if debug == 2 {
		fmt.Println("Lastlabel stack before push")
		for i, v := range pc.lastLabel { // iterate over the stack'n'print it.
			fmt.Println("Number", i, ":", v)
		}
	}

	// FIXME: Fix this!
	// #Brechstangen Methode
	pc.lastLabel[8] = pc.lastLabel[7]
	pc.lastLabel[7] = pc.lastLabel[6]
	pc.lastLabel[6] = pc.lastLabel[5]
	pc.lastLabel[5] = pc.lastLabel[4]
	pc.lastLabel[4] = pc.lastLabel[3]
	pc.lastLabel[3] = pc.lastLabel[2]
	pc.lastLabel[2] = pc.lastLabel[1]
	pc.lastLabel[1] = pc.lastLabel[0]
	pc.lastLabel[0] = name

	if debug == 2 {
		fmt.Println("Lastlabel stack after push:")
		for i, v := range pc.lastLabel {
			fmt.Println("Number", i, ":", v)
		}
	}

}

// popLastLabel() pops a lable from the stack. The label is returned as a string.
func (pc *programCode) popLastLabel() string {

	// These debug messags show how the stack was changed. See
	// pushLastLabel(name string) for more information

	if debug == 2 {
		fmt.Println("Lastlabel stack before pop:")
		for i, v := range pc.lastLabel {
			fmt.Println("Number", i, ":", v)
		}
	}

	// Popping labels off the stack just works fine. No one fears a
	// Brechstangen-Methode to appear here anytime soon.

	label := ""
	if len(pc.lastLabel) != 0 {
		label = pc.lastLabel[0]
	}

	if len(pc.lastLabel)-1 > 1 {
		pc.lastLabel = pc.lastLabel[1 : len(pc.lastLabel)-1]
	}

	// These debug messags show how the stack was changed
	if debug == 2 {
		fmt.Println("Lastlabel stack after pop:")
		for i, v := range pc.lastLabel {
			fmt.Println("Number", i, ":", v)
		}
	}
	return label
}

// FIXME: DONE
// <s> The BAUSTELLE! : Solved on Monday July 27th
// For loops are working but still strange to use. The loopvariable
// can('t) be accessed by their predefined name and appended counter
// number e.g. loopVar0, loopVar1, loopVar3 counting is still
// necessarry.  Todo: Change loopVar32 to something more general like </s>
//
// for loops just work fine

// This is the code snipped checking the condition inside an assembly loop.

func (pc *programCode) createForCheck(loopVar string) {
	code := "\n\tmov rax, [" + loopVar + "] \t; for-loop\n"
	code += "\tdec rax\n\tmov [" + loopVar + "], rax\n"
	forJmpBackLabel := pc.popLastLabel()
	code += "\tcmp rax, 0\n\tjle " + forJmpBackLabel + "\t; if zero close loop\n\t \n" // Fixed this line

	pc.appendCode(code)
}

// createCmp(a, b string) initialises a comparison of two values in
// the assembly code. The funtion tries to read a variable identifier,
// but it'll interpret the token as a numeric value if this is
// possible.

/*
;; Assembly for n00bs.
mov rax, [a]
mov rbx, [b]
cmp rax

;; check camparison with conditional jump (|j**|) after |cmp|.
*/

// necessary for conditionss in assembly
// a, b are variable identifier; or may be numbers stored in strings

func (pc *programCode) createCmp(a, b string) {
	code := "\t\t\t; compare " + a + " with " + b + "\n"
	if _, err := strconv.Atoi(a); err == nil {
		code += "\tmov rax, " + a + "\n"
	} else {
		code += "\tmov rax, [" + a + "]\n"
	}

	if _, err := strconv.Atoi(b); err == nil {
		code += "\tmov rbx, " + b + "\n"
	} else {
		code += "\tmov rbx, [" + b + "]\n"
	}
	code += "\tcmp rax, rbx\n"
	pc.appendCode(code)

}

// The following methods are responible for the boolean logic of the
// compiled program. They create a jump to the If-Satementsbody code,
// when their condition is true. No Expression handling is
// involved. Everything is hard coded.

func (pc *programCode) isEqual(label string) {
	code := ""
	code += "\tmov rax, 1\t; check equality\n"
	if label != "" {
		code += "\tje " + label + "\t; if so jump to " + label + "\n"
	}

	// Why if var != "" { ... }. What was the reason for this? What
	// error should it prevent?

	pc.appendCode(code)

}

func (pc *programCode) isGreater(label string) {
	code := ""
	code += "\tmov rax, 1\t; check equality\n"
	if label != "" {
		code += "\tjg " + label + "\t; if so jump to " + label + "\n"
	}

	pc.appendCode(code)

	// Ideas of handling boolean expressions
	// jge call true ??
	// mov rax, 1 // true
	// mov rax, 0 // false
}
func (pc *programCode) isSmaller(label string) {
	code := ""
	code += "\tmov rax, 1\t; check equality\n"
	if label != "" {
		code += "\tjl " + label + "\t; if so jump to " + label + "\n"
	}

	pc.appendCode(code)

	// mov rax, 1 // true
	// mov rax, 0 // false
}

func (pc *programCode) isGreaterEqual(label string) {
	code := ""
	code += "\tmov rax, 1\t; check equality\n"
	if label != "" {
		code += "\tjge " + label + "\t; if so jump to " + label + "\n"
	}

	pc.appendCode(code)

	// mov rax, 1 // true
	// mov rax, 0 // false
}

func (pc *programCode) isSmallerEqual(label string) {
	code := ""
	code += "\tmov rax, 1\t; check equality\n"
	if label != "" {
		code += "\tjle " + label + "\t; if so jump to " + label + "\n"
	}

	pc.appendCode(code)
	// mov rax, 1 // true
	// mov rax, 0 // false
}

// --------- END OF COMPARISON OPERATOR CODE ------------

// createStart() writes an assembly template to the code. An official header, a start label

func (pc *programCode) createStart() {
	start := "section .text\nglobal _start\n_start:\n"
	pc.code += start
}

// createExit(val srting) writes a exit statement to the code. The
// return status is transmitted as argument string.

// assembly for n00bs:
// mov rax, 60 ;; close me
// mov rdi, 0 ;; no error
// syscall ;; linux please!

func (pc *programCode) createExit(val string) {
	code := ""
	code += "\tmov rax, 60\t; exit program\n\tmov rdi, " + val + "\n\tsyscall\n"
	pc.funcCode[0] += code
	// Appends this code snippet to the first
	// level of indentation e.g. main-function
}

// createAllFunctions() adds all the functions' buffers to the final
// sourced code at the end of the compiling process.

func (pc *programCode) createAllFunctions() {
	for _, v := range pc.funcCode {
		pc.code += v // Finally chunk everything to a string!
	}
}

/*
mov eax, var1Value
mov [var1], eax
*/

// The BSS-Segmet allows to reserve n-byte sized space in the main memory.
// Those space chunks are tagged with a name which represents their address in the source code.

// [someVar] -- get the value of the memory chunk behind the someName like *someVar
// someVar	-- gets the address of the memory chunk like &someVar

// initBssVars() fills the reserverd space with variables' values.
func (pc *programCode) initBssVars() {
	code := "\t\t\t;; fill .bss variables with their values\n"
	for k, v := range pc.intMap {
		s := strconv.FormatInt((int64)(v), 10)
		code += "\tmov rax, " + s + "\n\tmov [" + k + "], rax\n"
	}
	pc.code += code
}

// initVars gives every memory chunk, its value it is supposed to
// hold. The final compiles MiniPython program runs the prodruced code
// of initVar(...) at the beginning of its main function.

func (pc *programCode) initVar(s, k string) {
	pc.code += "\tmov rax, " + k + "\t; newVar \n\tmov [" + s + "], rax\n" // start with these instructions
}

// Inline assembly can be written with the asm method its arguments
// are directly appended to the active function buffer. THIS IS REALLY
// HELPFULL to urgently change some internal code without
// rewriting/rereading the whole compiler code especially when bugs
// appear.

func (pc *programCode) asm(code string) {
	pc.appendCode("\t" + code + "\n")
}

// createBss() creates a code snipped which reserves space for the
// variables.  This creates the .bss Segment at the end of the
// assembly output code. The BSS-segment contains mutable memory
// space. Which is ideal for storing all my variables there

func (pc *programCode) createBss() {
	bssString := "\nsection .bss\n"
	for v := range pc.intMap {
		bssString += "\t" + v + ": resb 8" + "\n"
	}
	pc.code += bssString // appends this Assembly code to the end after creating all the functions
}

// createData() creates a code snippet which creates all the string
// constants

func (pc *programCode) createData() {
	dataString := "\nsection .data\n"
	for k, v := range pc.stringMap {
		dataString += "\t" + k + ": db " + v + "\n"
	}
	pc.code += dataString
}
