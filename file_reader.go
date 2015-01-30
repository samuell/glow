package glow

import (
	"bufio"
	"log"
	"os"
)

type FileReader struct {
	InFilePath chan string
	Out        chan []byte
}

func (self *FileReader) OutChan() chan []byte {
	self.Out = make(chan []byte, 16)
	return self.Out
}

func (self *FileReader) Init() {
	go func() {
		file, err := os.Open(<-self.InFilePath)
		if err != nil {
			log.Fatal(err)
		}
		scan := bufio.NewScanner(file)
		for scan.Scan() {
			self.Out <- append([]byte(nil), scan.Bytes()...)
		}
		close(self.Out)
	}()
}
