package main

import (
	"os"
	"runtime/debug"
	"strings"

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
    is_noide := false
    bi , noerr := debug.ReadBuildInfo()
    if !noerr{
        return
    }
    for _,item := range bi.Settings{
        if item.Key == "-tags" && strings.Contains(item.Value , "noide"){
            is_noide = true
            break
        }
    }

	cmd.Execute(is_noide)
}
