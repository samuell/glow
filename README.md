## Glow - Simplistic library of (streaming) Go(lang) workflow components for scientific and bioinformatics workflows

This is a work in progress, exploring how far we can get in writing in a "flow-based programming" style, completely without any frameworks,just using pure Go channels and go-routines, and just using a design pattern.

### Example usage

This is an example program, utilising the StdInReader, BaseComplementer, and the Printer component, to do base complement processing of fasta file content that it piped to the program.

### Example component, STDIN reader

First let's just have a look at how a component looks. Every component has one or more "in" and "outports", consisting of struct-fields of type channel (of some type that you choose. []byte arrays in this case). Then it has a run method that initializes a go-routine, and reads on the inports, and writes on the outports, as it processes incoming "data packets":

````go
package glow

import (
	"bufio"
	"os"
)

func NewStdInReader(outChan chan []byte) *StdInReader {
	stdInReader := new(StdInReader)
	stdInReader.Out = outChan
	stdInReader.Init()
	return stdInReader
}

type StdInReader struct {
	Out chan []byte
}

func (self *StdInReader) Init() {
	go func() {
		scan := bufio.NewScanner(os.Stdin)
		for scan.Scan() {
			self.Out <- append([]byte(nil), scan.Bytes()...)
		}
		close(self.Out)
	}()
}
````

#### Connecting channels and components manually

And then, to connect such processes together, we just create a bunch of channels, a bunch of processes, and then stitch them together, and run it!

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

#### Using New... convenience functions

... we can save a lot of keystrokes, and make the code shorter, and maybe more readable:

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


Finally, to compile and run the program above, do like this:
````bash
go build basecomplement.go
cat SomeFastaFile.fa | ./basecomplement > SomeFastaFile_Basecomplemented.fa
````



### Related projects / See also
- Blow - https://github.com/samuell/blow
- BioGo - https://code.google.com/p/biogo
