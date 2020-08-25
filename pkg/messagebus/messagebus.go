package messagebus

import "fmt"

type messageBus struct {
	channelMap map[string]chan interface{}
}

var bus = messageBus{
	make(map[string]chan interface{}),
}

func RegisterObserver(chanName string, observer Observer) {
	channel := make(chan interface{})
	bus.channelMap[chanName] = channel

	go func() {
		for {
			select {
			case ev, ok := <-channel:
				if ok {
					observer.OnEvent(ev)
				} else {
					return
				}
			}
		}
	}()
}

func SendMsg(chanName string, event interface{}) error {
	if channel, ok := bus.channelMap[chanName]; ok {
		channel <- event
	} else {
		return fmt.Errorf("could not found channel %s", chanName)
	}

	return nil
}

func CloseAll() {
	for _, channel := range bus.channelMap {
		close(channel)
	}
}
