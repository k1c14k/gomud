package vm

import "log"

type Command interface {
	Handle(vm *VirtualMachine)
}

type StopCommand struct{}

type MethodCallCommand struct {
	object          Object
	method          string
	arguments       []Value
	contextProvider ContextProvider
}

func NewMethodCallCommand(object Object, method string, arguments []Value, contextProvider ContextProvider) *MethodCallCommand {
	return &MethodCallCommand{
		object:          object,
		method:          method,
		arguments:       arguments,
		contextProvider: contextProvider,
	}
}

func (c *MethodCallCommand) Handle(vm *VirtualMachine) {
	vm.execute(c.object, c.method, c.arguments, c.contextProvider)
}

func (c *StopCommand) Handle(vm *VirtualMachine) {
	close(vm.commandChannel)
	log.Println("VM stopped")
}

type VirtualMachine struct {
	commandChannel chan Command
	classes        map[string]Class
}

var instance *VirtualMachine

func GetVirtualMachine() *VirtualMachine {
	if instance == nil {
		instance = &VirtualMachine{
			commandChannel: make(chan Command),
			classes:        make(map[string]Class),
		}
	}
	return instance
}

func (vm *VirtualMachine) Run() {
	log.Println("VM started")
	for {
		command := <-vm.commandChannel
		command.Handle(vm)
	}
}

func (vm *VirtualMachine) Stop() {
	vm.commandChannel <- &StopCommand{}
}

func (vm *VirtualMachine) getClass(name string) Class {
	if _, ok := vm.classes[name]; !ok {
		vm.classes[name] = *newClass(name)
	}

	return vm.classes[name]
}

func (vm *VirtualMachine) execute(object Object, method string, arguments []Value, contextProvider ContextProvider) {
	ef := NewExecutionFrame(contextProvider)

	calleeObjectValue := *NewObjectValue(&object)
	methodValue := NewStringValue(method)

	for _, arg := range arguments {
		ef.valueStack.push(arg)
	}

	ef.call(calleeObjectValue, methodValue)
}

func GetCommandChannel() chan Command {
	return GetVirtualMachine().commandChannel
}
