package vm

import (
	"fmt"
	"hellclientswitch/modules/app"
	"os"
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/herb-go/util"
)

type VM struct {
	API     VmAPI
	runtime *goja.Runtime
	lock    sync.Mutex
	ticker  *time.Ticker
}

func (vm *VM) initAPI() {
	vm.runtime.Set("Send", func(call goja.FunctionCall, r *goja.Runtime) goja.Value {
		vm.API.APISendMessage(call.Argument(0).String(), call.Argument(1).String())
		return nil
	})
	vm.runtime.Set("Broadcast", func(call goja.FunctionCall, r *goja.Runtime) goja.Value {
		vm.API.APIBroadcast(call.Argument(0).String())
		return nil
	})
	vm.runtime.Set("Print", func(call goja.FunctionCall, r *goja.Runtime) goja.Value {
		println(call.Argument(0).String())
		return nil
	})

}
func (vm *VM) Start() {
	vm.runtime = goja.New()
	if app.System.Script != "" {
		data, err := os.ReadFile(util.AppData(app.System.Script))
		if err != nil {
			panic(err)
		}
		_, err = vm.runtime.RunScript(app.System.Script, string(data))
		if err != nil {
			panic(err)
		}
		vm.initAPI()
		if app.System.TickerDurationInSeconds > 0 {
			vm.ticker = time.NewTicker(time.Duration(app.System.TickerDurationInSeconds) * time.Second)
			go func() {
				for _ = range vm.ticker.C {
					vm.OnTicker()
				}
			}()
		}
	}
}
func (vm *VM) Call(source string, args ...interface{}) goja.Value {
	vm.lock.Lock()
	defer vm.lock.Unlock()
	s, err := vm.runtime.RunString(source)
	if err != nil {
		util.LogError(err)
		return nil
	}
	fn, ok := goja.AssertFunction(s)
	if !ok {
		util.LogError(fmt.Errorf("js function %s not found", source))
		return nil
	}
	jargs := []goja.Value{}
	for _, v := range args {
		jargs = append(jargs, vm.runtime.ToValue(v))
	}
	var result goja.Value
	var scripterr error
	err = util.Catch(func() {
		result, scripterr = fn(goja.Undefined(), jargs...)
	})
	if scripterr != nil {
		util.LogError(scripterr)
		return nil
	}
	if err != nil {
		util.LogError(err)
		return nil
	}
	return result
}
func (vm *VM) Send(id string, msg string) {

}
func (vm *VM) OnMessage(id string, msg string) bool {
	if app.System.OnMessage != "" {
		result := vm.Call(app.System.OnMessage, id, msg)
		if result != nil {
			return result.ToBoolean()
		}
	}
	return true
}
func (vm *VM) OnTicker() {
	if app.System.OnTicker != "" {
		vm.Call(app.System.OnTicker, app.System.TickerDurationInSeconds)
	}
}

func Create(api VmAPI) *VM {
	return &VM{
		API: api,
	}
}
