section .text
global _start
_start:
	mov rax, 10	; newVar 
	mov [NL], rax
	mov rax, 32	; newVar 
	mov [WS], rax
	mov rax, 90	; newVar 
	mov [a], rax
	mov rax, 4	; newVar 
	mov [loopVar1], rax
	mov rax, 5	; newVar 
	mov [loopVar2], rax
	jmp forLabel1	; jmp to forLabel1
	forReturn1:	; return point
	mov rax, 4	;reset LoopVar to4
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
	 
	mov rax, 1	;print StringstrConst0
	mov rdi, 1
	mov rdx, 6
	mov rsi, strConst0
	syscall
	mov rax, 1	;print NL
	mov rdi, 1
	mov rdx, 1
	mov rsi, NL
	syscall
	jmp forLabel2	; jmp to forLabel2
	forReturn2:	; return point
	mov rax, 5	;reset LoopVar to5
	 mov [loopVar2], rax;	 done
	jmp forLabel1	; return to last place

forLabel2:

	mov rax, [loopVar2] 	; for-loop
	dec rax
	mov [loopVar2], rax
	cmp rax, 0
	jle forReturn2	; if zero close loop
	 
	mov rax, 1	;print StringstrConst1
	mov rdi, 1
	mov rdx, 11
	mov rsi, strConst1
	syscall
	mov rax, 1	;print NL
	mov rdi, 1
	mov rdx, 1
	mov rsi, NL
	syscall
	jmp forLabel2	; return to last place

section .bss
	NL: resb 8
	WS: resb 8
	a: resb 8
	loopVar1: resb 8
	loopVar2: resb 8

section .data
	strConst0: db "Außen"
	strConst1: db "----> Innen"
