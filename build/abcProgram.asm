section .text
global _start
_start:
	mov rax, 10	; newVar 
	mov [NL], rax
	mov rax, 32	; newVar 
	mov [WS], rax
	mov rax, 90	; newVar 
	mov [a], rax
	mov rax, 91	; newVar 
	mov [loopVar1], rax
	mov rax, 91	; newVar 
	mov [a], rax
	jmp forLabel1	; jmp to forLabel1
	forReturn1:	; return point
	mov rax, 91	;reset LoopVar to91
	 mov [loopVar1], rax;	 done
	mov rax, 1	;print NL
	mov rdi, 1
	mov rdx, 1
	mov rsi, NL
	syscall
	mov rax, 60	; exit program
	mov rdi, 0
	syscall

forLabel1:

	mov rax, [loopVar1] 	; for-loop
	dec rax
	mov [loopVar1], rax
	cmp rax, 0
	jle forReturn1	; if zero close loop
	 

			; Substract loopVar1 from a and save difference in a
	mov rax, [a]
	sub rax, [loopVar1]
	mov [a], rax
			; compare a with 65
	mov rax, [a]
	mov rbx, 65
	cmp rax, rbx
	mov rax, 1	; check equality
	jge ifLabel2	; if so jump to ifLabel2
	ifReturn2:	; return point
	mov rax, 91	;set a to 91
	mov [a], rax 
	jmp forLabel1	; return to last place

ifLabel2:
	mov rax, 1	;print a
	mov rdi, 1
	mov rdx, 1
	mov rsi, a
	syscall
	jmp ifReturn2	; return to last place

section .bss
	a: resb 8
	loopVar1: resb 8
	NL: resb 8
	WS: resb 8

section .data
