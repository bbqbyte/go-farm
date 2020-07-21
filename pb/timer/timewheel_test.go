package timer

import (
	"fmt"
	"testing"
	"time"
)

var tw = NewTimeoutWheel(WithTickInterval(time.Second))

func TestWheel(t *testing.T) {
	_, err := tw.Schedule(time.Second, testw, nil)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		//tout.Stop()
	}
	fmt.Println("sleeping")
	time.Sleep(5*time.Second)
	fmt.Println("wakeup")
}

func testw(args interface{}) {
	fmt.Println("oooooo")
	tw.Schedule(time.Second, testw, nil)
}