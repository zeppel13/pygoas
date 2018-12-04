# pygoas
Ein Schüler baut einen Compiler

Pygoas übersetzt eine sehr einfache an Python angelehnte Sprache in NASM-Linux-AMD64 Assembler und kompiliert diese mit NASM und LD zu einer ausführbarem ELF64 das auf Linux läuft. Dieses Programm habe ich hauptsächlich während der Schulzeit und mitten in der Nacht zwischen Neujahr und Mai 2016 in Go geschrieben. Die Idee mich an einem kleinen Compiler zu versuchen, kam mir, als ich versuchte einen Programmierbaren Taschenrechner mit meinem dünnen Assemblerwissen zu verbinden. Das Projekt zeigt, dass man mit genug Kreativität, Verstand und Koffein auch als Schüler einen Compiler eigener Art „entwickeln“ kann.

Ursprünglich war das Projekt auf https://sebastiankind.de/pygoas.git gehostet, wo es auch weiterhin abrufbar bleibt. Zwecks besserem Zugriff lade ich meinen Code auf GitHub hoch.



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






