# pygoas
Ein Schüler baut einen Compiler

Pygoas übersetzt eine sehr einfache an Python angelehnte Sprache in NASM-Linux-AMD64 Assembler und kompiliert diese mit NASM und LD zu einer ausführbarem ELF64-Programm, das auf Linux läuft. Die Idee mich an einem kleinen Compiler zu versuchen, kam mir, als ich versuchte einen Programmierbaren Taschenrechner mit meinem dünnen Assemblerwissen zu verbinden. Das Projekt zeigt, dass man mit genug Kreativität, Verstand und Koffein auch als Schüler einen Compiler eigener Art „entwickeln“ kann.  Dieses Programm habe ich hauptsächlich während der Schul- und Freizeit und mitten in der Nacht zwischen Neujahr und Mai 2016 in Go geschrieben.

Ursprünglich war das Projekt auf https://sebastiankind.de/pygoas.git gehostet, wo es auch weiterhin abrufbar bleibt. Zwecks besserer Einsicht lade ich meinen Code auf GitHub hoch.



```
cd pygoas
go build
```

Außerdem ist es relevant den Linker ld (Bestandteil von GCC) und den Assemblercompiler NASM installiert zu haben. Dieser befindet sich in aller Regel in den Paketquellen der üblichen Linuxdistributionen.

Im Verzeichnis Beispiele befinden sich Programme, die sich mit dem Compiler übersetzen lassen. Der Sprachumfang ist recht rudimentär, hatte jedoch auch nie den Anspruch sich mit großen Sprachen oder Projekten aus Universitären Umfeld zu messen. Er erlaubt keinen verschachtelten Ausdrücke, jedoch ist es im wesentlichen möglich:

- Funktionen zu definieren und auszuführen
- Ein 10 Stufen tiefer Funktionsaufrufstack exisitert
- Variabeln und Stringkonstanten zu definieren
- For-Schleifen und If-Blöcke verwenden im Speziellen Integervariabeln mit ihren Literalen zu vergleichen
- Text und Zahlen zu printen (e.g. print("hi"))
- inline Assemblercode zu verwenden, der vom Compiler direkt in die Ausgabe eingefügt wird
- Kommentarte die bis zum Ende der Zeile alle Symbole ignorieren

Außerdem befindet sich ein HTTP-Server im Compilercode, der angedacht war bei Presentationen das compilierte Programm auf allen Workstations eines Computerraums, CIP-Pools, LAN, etc anzuzeigen. 


```
NL = 10 # Neue Zeile
TB = 9 # Tabulator

for (3):
    print ("Außen", NL)
    for (4):
        print (TB, "Innen", NL)
        #endfor
    #endfor
print (NL)
```

Mit lieben Grüßen

Sebastian Kind


## Anhang

Im Folgen füge ich hier die alte readme.txt ein die auch im Unterverzeichnis Doku zu finden ist. 

#+Date: 20.05.2015
#+Author: Sebastian Kind
#+EMail: mail Ð sebastiankind.de


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

```
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
```

Ende der Einleitung




