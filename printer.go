package glow

import (
	"bufio"
	"fmt"
	"os"
)

func NewPrinter(InChan chan []byte, DrivingBeltChan chan int) *Printer {
	printer := new(Printer)
	printer.In = InChan
	printer.DrivingBelt = DrivingBeltChan
	printer.Init()
	return printer
}

type Printer struct {
	In          chan []byte
	DrivingBelt chan int
}

func (self *Printer) Init() {
	go func() {
		w := bufio.NewWriter(os.Stdout)
		for line := range self.In {
			if len(line) > 0 {
				fmt.Fprintln(w, string(line))
			}
			self.DrivingBelt <- 1
		}
		w.Flush()
		close(self.DrivingBelt)
	}()
}
