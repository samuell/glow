## Glow - Simplistic library of (streaming) Go(lang) workflow components for scientific and bioinformatics workflows

This is a work in progress.

### Example usage

This is an example program, utilising the StdInReader, BaseComplementer, and the Printer component, to do base complement processing of fasta file content that it piped to the program.

````go
package main

import (
	"fmt"
	"github.com/samuell/glow"
)

const (
	BUFSIZE = 2048 // Set a buffer size to use for channels
)

func main() {
	// Create channels / connections
	stdIn_BaseCompl := make(chan []byte, BUFSIZE)
	baseCompl_Printer := make(chan []byte, BUFSIZE)
	drivingBelt := make(chan int, 0)

	// Create components, connecting the channels
	glow.NewStdInReader(stdIn_BaseCompl)
	glow.NewBaseComplementer(stdIn_BaseCompl, baseCompl_Printer)
	glow.NewPrinter(baseCompl_Printer, drivingBelt)

	// Loop over the last channel, to drive the execution
	cnt := 0
	for i := range drivingBelt {
		cnt += i
	}
	fmt.Println("Processed ", cnt, " lines.")
}
````

Compile and run like this:
````bash
go build basecomplement.go
cat SomeFastaFile.fa | ./basecomplement > SomeFastaFile_Basecomplemented.fa
````

### Related projects / See also
- Blow - https://github.com/samuell/blow
- BioGo - https://code.google.com/p/biogo
