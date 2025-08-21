package main

import "devbox/internal/commands"

type ParserArgs struct {
	commands.SharedCmdArgs
	Verbose            bool
	InstallCmdFilePath string
	LogFilePath        string
}
