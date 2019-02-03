package scripting

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/robertkrimen/otto"
)

var _ Engine = &JavaScriptEngine{}

// JavaScriptEngine stores scripting engine state
type JavaScriptEngine struct {
	vm *otto.Otto
}

// New instantiates a new scripting engine
func New() (engine *JavaScriptEngine) {
	engine = &JavaScriptEngine{
		vm: otto.New(),
	}

	return
}

// LoadScripts implements Engine
func (engine *JavaScriptEngine) LoadScripts(dirname string) (err error) {
	err = filepath.Walk(dirname, func(path string, fileInfo os.FileInfo, err error) error {
		if fileInfo.IsDir() {
			return nil
		}

		if !strings.HasSuffix(fileInfo.Name(), ".js") {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return errors.Wrap(err, path)
		}
		_, err = engine.vm.Run(file)
		if err != nil {
			return errors.Wrap(err, "failed to run script")
		}

		return nil
	})

	return err
}

// OnMessageSend implements Engine
func (engine *JavaScriptEngine) OnMessageSend(oldText string) (newText string) {
	jsValue, jsError := engine.vm.Run(fmt.Sprintf(`onMessage("%s")`, oldText))
	if jsError != nil {
		//TODO Return error?
		return oldText
	}
	return jsValue.String()
}