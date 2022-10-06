package main

import (
	"os"

	/*
		"vabna/evaluator"
		"vabna/lexer"
		"vabna/object"
		"vabna/parser"
	*/
	"pankti/cmd"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.ErrorLevel)
	log.SetFormatter(&log.TextFormatter{
		PadLevelText:  true,
		FullTimestamp: true,
	})

	log.SetOutput(os.Stdout)
}

func main() {
	cmd.Execute()
}
