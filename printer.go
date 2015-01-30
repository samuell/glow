package glow

import (
	"bufio"
	"fmt"
	"os"
)

type Printer struct {
	In          chan []byte
	DrivingBelt chan int
}

func (self *Printer) DrivingBeltChan() chan int {
	self.DrivingBelt = make(chan int)
	return self.DrivingBelt
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
