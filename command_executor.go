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
		// Set up the command to run
		cmdParts := strings.Fields(<-self.CommandIn)
		cmdHead := cmdParts[0]
		cmdTail := cmdParts[1:len(cmdParts)]
		cmd := exec.Command(cmdHead, cmdTail...)

		// Connect a buffer to the stdout of the command
		var out bytes.Buffer
		cmd.Stdout = &out
		scan := bufio.NewScanner(&out)

		// Start executing the command
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}

		// Write the command output to the output channel LinesOut
		for scan.Scan() {
			self.LinesOut <- append([]byte(nil), scan.Bytes()...)
		}
		close(self.LinesOut)
	}()
}
