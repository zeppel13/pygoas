#+Date: 20.05.2015
#+Author: Sebastian Kind
#+EMail: mail@sebastiankind.de


Mein Compiler (Noch ohne Namen) übersetzt sehr einfachen Pythoncode
in Assembler bzw. lässt die Ausgaben in native Linuxprogramme
kompilieren und linken.

* ***** Eine kleine Anleitung ***** *

| Systemvorausetzungen: | Lösung (ins Terminal)             |
|-----------------------+-----------------------------------|
| Linux                 | uname && uname -p                 |
| x86_64                |                                   |
| go (Compiler)         | sudo apt-get install golang       |
| nasm                  | sudo apt-get install nasm         |
| ld (Gnu Linker)       | sollte i.d.R. vorinstalliert sein |
|-----------------------+-----------------------------------|

nasm und ld sollten im Anschluss unter 

	/usr/bin/nasm
	/usr/bin/ld 

liegen, sonst wird der Pythoncompiler nicht übersetzen und
umgangssprachliche Fehlermeldungen von sich geben.


* ***** Kompilieren ***** *


Selbst kompilieren mit Go geht ganz einfach:

	cd ./compiler
	go build

Da go Platformunabhängig ist, kann mein Compiler auch auf anderen
System als Linux x86_64 benutzt werden, jedoch schränken sich die
möglichen Kommandos auf

	./compiler -s ausgabe.asm ./input.py
	./compiler --stdout ./input.py

ein, da hierbei nicht auf nasm oder ld zurückgegriffen werden muss.

* ***** Compiler 101 ***** *

Im Anschluss kann man mit dem Tool ausprobieren, bis man Bugs findet. In der
Regel sieht ein Befehl so aus. Das Ergebnis ist hier eine ausführbare 
Datei. Der Sprachumfang ist immer noch sehr begrenzt.

	./compiler -o ausgabe.elf ./input.py
	

	./ausgabe.elf

Führt die das Fertige Programm aus.


So lässt sich die Arbeit des Compilers ohne kompilieren in eine Datei
umleiten

	./compiler --stdout ./input.py > ./ausgabe.asm

	Im Anschluss kann man mit dem Editor der Wahl die Datei
	(ausgabe.asm) öffnen (Emacs, Vim, Nano, Gedit, Notepad++,
	etc. sollten Syntaxhighlighting für Assembler haben).

oder auch einfach ins Terminal geben mit

	./compiler --stdout ./input.py | less

Es wird auffalen, dass der Compiler den Assemblercode zur besseren
Lesbarkeit selbständig kommentiert.

	
* ***** Flags ***** *

Das ist Liste aller Flags, die der Compiler bei Falscheingabe ausgibt:

Usage of ./compiler:
  -c	This runs the nasm compiler (default true)
  -debug
    	Print debug information e.g. labelStack, cryptic numbers and other stuff. This option renders the output useless.
  -l	This runs the ld linker (default true)
  -o string
    	Name of output file (default "prog")
  -s	Compile only to NASM assembly (default true)
  -stdout
    	Compile only to NASM assembly and print to stdout


Ende der Einleitung
