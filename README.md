## Glow - Simplistic library of (streaming) Go(lang) workflow components for scientific and bioinformatics workflows

This is a work in progress.

### Example usage

This is an example program, utilising the StdInReader, BaseComplementer, and the Printer component, to do base complement processing of fasta file content that it piped to the program.

#### Using New... convenience functions

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
	chan1 := make(chan []byte, BUFSIZE)
	chan2 := make(chan []byte, BUFSIZE)
	chan3 := make(chan int, 0)

	// Create components, connecting the channels
	glow.NewStdInReader(chan1)             // Here, chan1 is an output channel
	glow.NewBaseComplementer(chan1, chan2) // chan1 is input, chan2 is output
	glow.NewPrinter(chan2, chan3)          // chan2 is input, chan3 is output

	// Loop over the last channel, to drive the execution
	cnt := 0
	for i := range chan3 {
		cnt += i
	}
	fmt.Println("Processed ", cnt, " lines.")
}
````

#### Connecting channels and components manually

The above code is using convenient functions ("New...") for creating the components.
What is really happening under the hood, is the following, which is the same workflow
as above, but written without using the convenient functions:

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
	chan1 := make(chan []byte, BUFSIZE)
	chan2 := make(chan []byte, BUFSIZE)
	chan3 := make(chan int, 0)

	// Create components, connecting the channels
	stdInReader := new(glow.StdInReader)
	stdInReader.Out = chan1
	stdInReader.Init()

	baseCompler := new(glow.BaseComplementer)
	baseCompler.In = chan1
	baseCompler.Out = chan2
	baseCompler.Init()

	printer := new(glow.Printer)
	printer.In = chan2
	printer.DrivingBelt = chan3
	printer.Init()

	// Loop over the last channel, to drive the execution
	cnt := 0
	for i := range chan3 {
		cnt += i
	}
	fmt.Println("Processed ", cnt, " lines.")
}
````

Finally, to compile and run the program above, do like this:
````bash
go build basecomplement.go
cat SomeFastaFile.fa | ./basecomplement > SomeFastaFile_Basecomplemented.fa
````



### Related projects / See also
- Blow - https://github.com/samuell/blow
- BioGo - https://code.google.com/p/biogo
