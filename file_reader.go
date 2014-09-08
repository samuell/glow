package glow

import (
	"bufio"
	"log"
	"os"
)

func NewFileReader(inFilePathChan chan string, outChan chan []byte) *FileReader {
	fileReader := new(FileReader)
	fileReader.InFilePath = inFilePathChan
	fileReader.Out = outChan
	fileReader.Init()
	return fileReader
}

type FileReader struct {
	InFilePath chan string
	Out        chan []byte
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
