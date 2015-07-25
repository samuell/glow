# LIBRARY DEPRECATED - Replaced by: [scipipe.org](http://scipipe.org)

## Glow - Simplistic library of (streaming) Go(lang) workflow components for scientific and bioinformatics workflows

This is a work in progress, exploring how far we can get in writing in a "flow-based programming" style, completely without any frameworks,just using pure Go channels and go-routines, and just using a design pattern.

For an in depth explanation of the ideas behind the library, see [this post on GopherAcademy.org](http://blog.gopheracademy.com/composable-pipelines-pattern).

### Example usage

This is an example program, utilising the StdInReader, BaseComplementer, and the Printer component, to do base complement processing of fasta file content that it piped to the program.

### Example component, STDIN reader

First let's just have a look at how a component looks. Every component has one or more "in" and "outports",
consisting of struct-fields of type channel (of some type that you choose. []byte arrays in this case). 
Then it has a run method that initializes a go-routine, and reads on the inports, and writes on the outports,
as it processes incoming "data packets". Finally, as you can see in the OutChan() method, it provides a 
convenience method for each outgoing chan field, that initializes the channel and returns it, which can later
be used for easier wiring of the network:

````go
package glow

import (
	"bufio"
	"os"
)

type StdInReader struct {
	Out chan []byte
}

func (self *StdInReader) OutChan() chan []byte {
	self.Out = make(chan []byte, 16)
	return self.Out
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
	BUFSIZE = 128 // Set a buffer size to use for channels
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

#### Using OutChan() convenience methods

... we can save some keystrokes by using the functions defined for each outport,
which initializes the outport with a channel, and returns the channel.

````go
package main

import (
	"fmt"
	"github.com/samuell/glow"
)

func main() {
	// Create channels / connections
	fileReader := new(glow.FileReader)
	baseComplementer := new(glow.BaseComplementer)
	printer := new(glow.Printer)

	// Connect components (THIS IS WHERE THE NETWORK IS DEFINED!)
	baseComplementer.In = fileReader.OutChan()
	printer.In = baseComplementer.OutChan()

	// Initialize / set up go-routines
	fileReader.Init()
	baseComplementer.Init()
	printer.Init()

	// The InFilePath channel has to be created manually
	fileReader.InFilePath = make(chan string)
	fileReader.InFilePath <- "test.fa"

	// Loop over the last channel, to drive the execution
	cnt := 0
	for i := range printer.DrivingBeltChan() {
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
