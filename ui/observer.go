package ui

type UIMsgObserver struct {
}

const (
	StartState                string = "start"
	StopState                        = "stop"
	AddEventStreamState              = "addEventStream"
	FinishAddEventStreamState        = "finishAddEventStream"
)

func (ob *UIMsgObserver) OnEvent(event interface{}) {
	state := event.(string)

	if state == StartState {
		appIns.isStarted.Store(true)
		appIns.isAddEventStream.Store(false)
	} else if state == StopState || state == FinishAddEventStreamState {
		appIns.isStarted.Store(false)
		appIns.isAddEventStream.Store(false)
	} else if state == AddEventStreamState {
		appIns.isStarted.Store(false)
		appIns.isAddEventStream.Store(true)

	}
}
