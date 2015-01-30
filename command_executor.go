package glow

import (
	"bufio"
	"bytes"
	"log"
	"os/exec"
	"strings"
)

type CommandExecutor struct {
	CommandIn chan string
	LinesOut  chan []byte
}

func (self *CommandExecutor) LinesOutChan() chan []byte {
	self.LinesOut = make(chan []byte, 16)
	return self.LinesOut
}

func (self *CommandExecutor) Init() {
	go func() {
		// Read command from in-port
		commandParts := strings.Fields(<-self.CommandIn)
		executable := commandParts[0]
		arguments := commandParts[1:len(commandParts)]
		// Create command object
		cmd := exec.Command(executable, arguments...)

		// Connect a buffer to the stdout of the command, and
		// usd in scanner
		var out bytes.Buffer
		cmd.Stdout = &out
		scan := bufio.NewScanner(&out)

		// Start executing the command
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}

		// (Copy and) write the command output to the out-port
		for scan.Scan() {
			self.LinesOut <- append([]byte(nil), scan.Bytes()...)
		}
		close(self.LinesOut)
	}()
}
