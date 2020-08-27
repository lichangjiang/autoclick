package action

import (
	"autoclick/constant"
	"autoclick/model"
	"autoclick/pkg/imageutil"
	"autoclick/pkg/messagebus"
	"fmt"
	"time"
)

func init() {
	fmt.Println("workObserver init")
	ob := &workObserver{}
	messagebus.RegisterObserver(constant.WorkObserverName, ob)
}

type workObserver struct {
	isStarted bool
}

func (ob *workObserver) OnEvent(ev interface{}) {
	event, ok := ev.(model.WorkMsg)
	if !ok {
		fmt.Println("WorkObserver fail cast Event to WorkMsg")
		return
	}

	if event.Msg == constant.StartState && !ob.isStarted {
		em := event.EventMap
		sen := event.StartEventNames

		if em == nil || sen == nil {
			fmt.Println("WorkObserver can not start beacuse eventMap or startEventNames is nil")
			return
		}

		startEventNames := map[string]bool{}
		eventMap := map[string]model.Event{}

		for _, n := range sen {
			startEventNames[n] = true
		}

		for k, v := range em {
			eventMap[k] = *v
		}

		eventMap, err := imageutil.CreateImageForEvents(eventMap)
		if err != nil {
			fmt.Printf("fail to createImage for events:%+v\n", err)
			return
		}
		ob.isStarted = true

		fmt.Println("WorkObserver isStarted change to true")
		go func() {
			for {
				if ob.isStarted {
					fmt.Println("image check goroutine start")
					err := imageutil.StartImageCheck(eventMap, startEventNames)

					if err != nil {
						fmt.Printf("checkImage error:%+v\n", err)
					}
					time.Sleep(3 * time.Second)
				} else {
					fmt.Println("image check goroutine stop")
					return
				}
			}
		}()
	} else if event.Msg == constant.StopState {
		ob.isStarted = false
		fmt.Println("WorkObserver isStarted change to false")
	}
}
