// This file was generated by counterfeiter
package fakes

import (
	. "github.com/cloudfoundry/cli/cf/terminal"
	"sync"
)

type FakePrinter struct {
	PrintStub        func(a ...interface{}) (n int, err error)
	printMutex       sync.RWMutex
	printArgsForCall []struct {
		arg1 []interface{}
	}
	printReturns struct {
		result1 int
		result2 error
	}
	PrintfStub        func(format string, a ...interface{}) (n int, err error)
	printfMutex       sync.RWMutex
	printfArgsForCall []struct {
		arg1 string
		arg2 []interface{}
	}
	printfReturns struct {
		result1 int
		result2 error
	}
	PrintlnStub        func(a ...interface{}) (n int, err error)
	printlnMutex       sync.RWMutex
	printlnArgsForCall []struct {
		arg1 []interface{}
	}
	printlnReturns struct {
		result1 int
		result2 error
	}
}

func (fake *FakePrinter) Print(arg1 ...interface{}) (n int, err error) {
	fake.printMutex.Lock()
	defer fake.printMutex.Unlock()
	fake.printArgsForCall = append(fake.printArgsForCall, struct {
		arg1 []interface{}
	}{arg1})
	if fake.PrintStub != nil {
		return fake.PrintStub(arg1)
	} else {
		return fake.printReturns.result1, fake.printReturns.result2
	}
}

func (fake *FakePrinter) PrintCallCount() int {
	fake.printMutex.RLock()
	defer fake.printMutex.RUnlock()
	return len(fake.printArgsForCall)
}

func (fake *FakePrinter) PrintArgsForCall(i int) []interface{} {
	fake.printMutex.RLock()
	defer fake.printMutex.RUnlock()
	return fake.printArgsForCall[i].arg1
}

func (fake *FakePrinter) PrintReturns(result1 int, result2 error) {
	fake.printReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakePrinter) Printf(arg1 string, arg2 ...interface{}) (n int, err error) {
	fake.printfMutex.Lock()
	defer fake.printfMutex.Unlock()
	fake.printfArgsForCall = append(fake.printfArgsForCall, struct {
		arg1 string
		arg2 []interface{}
	}{arg1, arg2})
	if fake.PrintfStub != nil {
		return fake.PrintfStub(arg1, arg2...)
	} else {
		return fake.printfReturns.result1, fake.printfReturns.result2
	}
}

func (fake *FakePrinter) PrintfCallCount() int {
	fake.printfMutex.RLock()
	defer fake.printfMutex.RUnlock()
	return len(fake.printfArgsForCall)
}

func (fake *FakePrinter) PrintfArgsForCall(i int) (string, []interface{}) {
	fake.printfMutex.RLock()
	defer fake.printfMutex.RUnlock()
	return fake.printfArgsForCall[i].arg1, fake.printfArgsForCall[i].arg2
}

func (fake *FakePrinter) PrintfReturns(result1 int, result2 error) {
	fake.printfReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakePrinter) Println(arg1 ...interface{}) (n int, err error) {
	fake.printlnMutex.Lock()
	defer fake.printlnMutex.Unlock()
	fake.printlnArgsForCall = append(fake.printlnArgsForCall, struct {
		arg1 []interface{}
	}{arg1})
	if fake.PrintlnStub != nil {
		return fake.PrintlnStub(arg1)
	} else {
		return fake.printlnReturns.result1, fake.printlnReturns.result2
	}
}

func (fake *FakePrinter) PrintlnCallCount() int {
	fake.printlnMutex.RLock()
	defer fake.printlnMutex.RUnlock()
	return len(fake.printlnArgsForCall)
}

func (fake *FakePrinter) PrintlnArgsForCall(i int) []interface{} {
	fake.printlnMutex.RLock()
	defer fake.printlnMutex.RUnlock()
	return fake.printlnArgsForCall[i].arg1
}

func (fake *FakePrinter) PrintlnReturns(result1 int, result2 error) {
	fake.printlnReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

var _ Printer = new(FakePrinter)
