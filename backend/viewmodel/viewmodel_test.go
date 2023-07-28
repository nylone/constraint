package viewmodel

import (
	"constraint/model"
	"sync"
	"testing"
	"time"
)

type client struct {
	in  chan<- Action
	out chan interface{}
}

var (
	count  int  = 0
	failed bool = false
	mutex  sync.Mutex
)

func checkCount(countExp int) {
	if count == countExp {
		count++
	} else {
		failed = true
	}
}

func TestAddClient(t *testing.T) {
	vm := NewViewmodel()
	c1 := client{
		out: make(chan interface{}),
	}
	c2 := client{
		out: make(chan interface{}),
	}
	go func() {
		mutex.Lock()
		defer mutex.Unlock()
		for v := range c1.out {
			switch v.(type) {
			case ModelUpdate:
				checkCount(4)
			case StartingInfo:
				checkCount(1)
			case NewClientInfo:
				checkCount(0)
			case GameClosed:
				checkCount(5)
				return
			case ChatMessage:
				checkCount(2)
			case ControllerResponse:
				checkCount(3)
			}
		}
	}()
	go func() {
		for range c2.out {
		}
	}()
	in, err := vm.AddClient("c1", c1.out)
	if err != nil {
		t.Fatal("viewmodel failed to add client correctly")
	}
	c1.in = in
	in, err = vm.AddClient("c2", c2.out)
	c2.in = in
	if err != nil {
		t.Fatal("viewmodel failed to add client correctly")
	}

	time.Sleep(1000000)

	c1.in <- Action{
		Id:  InputMsg,
		Msg: "hello, world!",
	}
	c1.in <- Action{
		Id:  InputAddPos,
		Pos: model.Pos{X: 0, Y: 0},
	}
	close(c1.in)

	mutex.Lock()
	defer mutex.Unlock()
	if failed {
		t.Fatal("viewmodel failed to send correct messages at count", count)
	}

	c3 := client{
		out: make(chan interface{}),
	}
	_, err = vm.AddClient("c3", c3.out)
	if err == nil {
		t.Fatal("viewmodel failed to reject client correctly")
	}

	vm = NewViewmodel()
	in, err = vm.AddClient("c3", c3.out)
	if err != nil {
		t.Fatal("viewmodel failed to add client correctly")
	}
	c1.in = in
	in, err = vm.AddClient("c3", c2.out)
	c2.in = in
	if err == nil {
		t.Fatal("viewmodel failed to reject client correctly")
	}

	time.Sleep(1000000)

	go func() {
		mutex.Lock()
		defer mutex.Unlock()
		for v := range c3.out {
			switch v.(type) {
			case ModelUpdate:
				failed = true
			case StartingInfo:
				checkCount(6 + 1)
			case NewClientInfo:
				failed = true
			case GameClosed:
				checkCount(6 + 3)
				return
			case ChatMessage:
				failed = true
			case ControllerResponse:
				checkCount(6 + 2)
			}
		}
	}()

	c1.in <- Action{
		Id:  InputAddPos,
		Pos: model.Pos{X: -1, Y: 0},
	}
	close(c1.in)
}
