package constant

const (
	StartBtnName string = "startBtn"
	AddESBtnName string = "addEventStreamBtn"
	AddEvBtnName string = "addEventBtn"
	ResetBtnName string = "resetBtn"

	AddEsBtnText  string = "add event stream"
	AddEvBtnText  string = "add event"
	StartBtnText  string = "start"
	ResetBtnText  string = "reset"
	FinishBtnText string = "finish"
	CommonText    string = "click what you want"

	ResetEventStreamBtnText string = "reset event stream"
	ResetEventBtxText       string = "reset event"
)

const (
	StartState                string = "start"
	StopState                        = "stop"
	AddEventStreamState              = "addEventStream"
	FinishAddEventStreamState        = "finishAddEventStream"

	UIStateObserverName = "UIStateObserver"
)

const (
	MouseEventObserverName string = "mouseEventObserver"
	HookObserverName              = "hookObserver"
)

const (
	EventStreamObserverName string = "EventStreamObserver"
	NewStreamMsg            string = "newStream"
	AddEventMsg             string = "addEvent"
	EndStreamMsg            string = "endEventStream"

	GlobalEventObserverName string = "globalEventObserver"
	AddEventStream          string = "addEventStream"
	FileAddEventStream      string = "fileAddEventStream"
	ResetEventStream        string = "resetEventStream"
)

const (
	JsonFileObserverName string = "JsonFileObserver"
)

const (
	WorkObserverName string = "workObserver"
)
