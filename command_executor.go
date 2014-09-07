package glow

import (
	"bufio"
	"bytes"
	"log"
	"os/exec"
	"strings"
)

func NewCommandExecutor(CommandIn chan string, LinesOut chan []byte) *CommandExecutor {
	commandExecutor := new(CommandExecutor)
	commandExecutor.CommandIn = CommandIn
	commandExecutor.LinesOut = LinesOut
	commandExecutor.Init()
	return commandExecutor
}

type CommandExecutor struct {
	CommandIn chan string
	LinesOut  chan []byte
}

func (self *CommandExecutor) Init() {
	go func() {
		// Set up things
		cmdParts := strings.Fields(<-self.CommandIn)
		cmdHead := cmdParts[0]
		cmdTail := cmdParts[1:len(cmdParts)]

		cmd := exec.Command(cmdHead, cmdTail...)
		var out bytes.Buffer
		cmd.Stdout = &out

		// Execute the actual command
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		scan := bufio.NewScanner(&out)
		for scan.Scan() {
			self.LinesOut <- append([]byte(nil), scan.Bytes()...)
		}
		close(self.LinesOut)
	}()
}
