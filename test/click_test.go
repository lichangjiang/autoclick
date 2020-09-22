package test

import (
	"testing"
	"time"

	"github.com/go-vgo/robotgo"
)

func TestMoveClick(t *testing.T) {
	for i := 0; i < 10; i++ {
		//robotgo.Move(100, 100)
		//robotgo.Click()
		robotgo.MoveClick(100, 100)
		time.Sleep(1 * time.Second)
		robotgo.MoveClick(200, 100)
		time.Sleep(1 * time.Second)
		robotgo.MoveClick(200, 200)
		time.Sleep(1 * time.Second)
	}

}
