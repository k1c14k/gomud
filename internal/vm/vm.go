package vm

import "log"

type Command interface {
	Handle(vm *VirtualMachine)
}

type StopCommand struct{}

type MethodCallCommand struct {
	object    Object
	method    string
	arguments []Value
	context   map[string]Object
}

func NewMethodCallCommand(object Object, method string, arguments []Value, context map[string]Object) *MethodCallCommand {
	return &MethodCallCommand{
		object:    object,
		method:    method,
		arguments: arguments,
		context:   context,
	}
}

func (c *MethodCallCommand) Handle(vm *VirtualMachine) {
	vm.execute(c.object, c.method, c.arguments, c.context)
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
		vm.classes[name] = NewClass(name)
	}

	return vm.classes[name]
}

func (vm *VirtualMachine) execute(object Object, method string, arguments []Value, ctxObjs map[string]Object) {
	var ctxObjValue = make(map[string]ObjectValue)
	for k, v := range ctxObjs {
		ctxObjValue[k] = ObjectValue{&v}
	}
	ctx := NewExecutionContext(object.GetClass().GetStringPool(), ctxObjValue)
	ef := NewExecutionFrame(ctx)

	for _, arg := range arguments {
		ef.valueStack.push(arg)
	}

	m := object.GetClass().GetMethod(method)
	m.Execute(ef)
}

func GetCommandChannel() chan Command {
	return GetVirtualMachine().commandChannel
}
