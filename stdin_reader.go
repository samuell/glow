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
