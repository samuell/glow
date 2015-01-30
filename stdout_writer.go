package glow

import (
	"bufio"
	"fmt"
	"os"
)

type StdOutWriter struct {
	In chan []byte
}

func (self *StdOutWriter) Init() {
	go func() {
		w := bufio.NewWriter(os.Stdout)
		for line := range self.In {
			if len(line) > 0 {
				fmt.Fprintln(w, string(line))
			}
		}
		w.Flush()
	}()
}
