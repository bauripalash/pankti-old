package evaluator

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"go.cs.palashbauri.in/pankti/ast"
	"go.cs.palashbauri.in/pankti/object"
)

func evalShowStmt(
	args []object.Obj,
	printBuff *bytes.Buffer,
	isGui bool,
) object.Obj {

	output := []string{}

	for _, item := range args {
		output = append(output, item.Inspect())
		//buff.Write([]byte(item.Inspect()))
	}

	if isGui {
		oldStdout := os.Stdout
		r, w, err := os.Pipe()
		if err != nil {
			log.Fatalf(err.Error())
		}

		os.Stdout = w

		outC := make(chan string)

		go func() {
			var buf bytes.Buffer
			io.Copy(&buf, r)
			outC <- buf.String()
		}()

		fmt.Println(strings.Join(output, ""))
		w.Close()
		os.Stdout = oldStdout
		out := <-outC

		printBuff.Write([]byte(out))
	} else {
		fmt.Println(strings.Join(output, ""))
	}

	return NULL
}

func evalBlockStmt(
	block *ast.BlockStmt,
	env *object.EnvMap,
	eh *object.ErrorHelper,
	printBuff *bytes.Buffer,
	isGui bool,
) object.Obj {

	var res object.Obj

	for _, stmt := range block.Stmts {
		res = Eval(stmt, env, *eh, printBuff, isGui)

		//fmt.Println("E_BS=> " , res)

		if res != nil {
			rtype := res.Type()
			if rtype == object.RETURN_VAL_OBJ || rtype == object.ERR_OBJ {
				//fmt.Println("RET => " ,  res)
				return res
			}
		}
	}
	//fmt.Println("EBS 2=>" ,res)
	return res
}
