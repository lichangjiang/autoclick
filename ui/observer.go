package ui

import (
	con "autoclick/constant"
	"autoclick/pkg/messagebus"
	"fmt"
)

type UIStateObserver struct {
}



func init() {
	fmt.Println("uiobserver init")

	ob := &UIStateObserver{}
	messagebus.RegisterObserver(con.UIStateObserverName, ob)
}

func (ob *UIStateObserver) OnEvent(event interface{}) {
	state := event.(string)

	if state == con.StartState {
		appIns.isStarted.Store(true)
		appIns.isAddEventStream.Store(false)
	} else if state == con.StopState || state == con.FinishAddEventStreamState {
		appIns.isStarted.Store(false)
		appIns.isAddEventStream.Store(false)
	} else if state == con.AddEventStreamState {
		appIns.isStarted.Store(false)
		appIns.isAddEventStream.Store(true)

	}
}
