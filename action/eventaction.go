package action

import (
	"autoclick/constant"
	"autoclick/controller"
	"autoclick/model"
	"autoclick/pkg/imageutil"
	"autoclick/pkg/messagebus"
	"fmt"
	"strconv"
	"time"
)

type EventStreamObserver struct {
	currentEventStream *model.EventStream
}

func init() {
	fmt.Println("EventStreamObserver init")
	ob := &EventStreamObserver{}
	messagebus.RegisterObserver(constant.EventStreamObserverName, ob)

	gob := &globalEventObserver{
		[]string{},
		map[string]*model.Event{},
		false,
	}
	messagebus.RegisterObserver(constant.GlobalEventObserverName, gob)
}

func (ob *EventStreamObserver) OnEvent(ev interface{}) {
	msg, ok := ev.(model.EventStreamMsg)

	if !ok {
		fmt.Println("EventStreamObserver Event fail cast to EventStreamMsg")
		return
	}

	if msg.Msg == constant.AddEventMsg {
		fmt.Printf("EventStreamObserver receive AddEventMsg:%+v\n", &msg)
		axis, ok := msg.Value.(model.Axis)
		if !ok {
			fmt.Println("AddEventMsg Value fail cast to Axis")
			return
		}

		if ob.currentEventStream == nil {
			fmt.Println("CurrentEventStream is nil")
			return
		}
		n := time.Now().Nanosecond()
		event := &model.Event{
			Name:      "event" + strconv.Itoa(n),
			Axis:      axis,
			Action:    "click",
			NextEvent: "",
			Revert:    false,
		}

		eventLen := len(ob.currentEventStream.Events)
		if eventLen > 0 {
			pre := ob.currentEventStream.Events[eventLen-1]
			pre.NextEvent = event.Name
		}
		ob.currentEventStream.Events = append(ob.currentEventStream.Events, event)
	} else if msg.Msg == constant.NewStreamMsg || msg.Msg == "reset" {
		fmt.Printf("EventStreamObserver %s\n", msg.Msg)
		cs := &model.EventStream{
			Events: []*model.Event{},
		}
		ob.currentEventStream = cs
	} else if msg.Msg == constant.EndStreamMsg {
		cs := ob.currentEventStream
		ob.currentEventStream = nil
		if cs != nil && len(cs.Events) > 0 {
			newMsg := model.EventStreamMsg{
				Msg:   constant.AddEventStream,
				Value: cs,
			}
			messagebus.SendMsg(constant.GlobalEventObserverName, newMsg)
		}
	}
}

type globalEventObserver struct {
	startEventNames []string
	eventMap        map[string]*model.Event
	isChanged       bool
}

func (ob *globalEventObserver) OnEvent(ev interface{}) {
	msg, ok := ev.(model.EventStreamMsg)

	if !ok {
		fmt.Println("globalEventObserver Event fail cast to EventStreamMsg")
		return
	}

	if msg.Msg == constant.AddEventStream ||
		msg.Msg == constant.FileAddEventStream {
		es, ok := msg.Value.(*model.EventStream)
		if !ok {
			fmt.Printf("EventStreamMsg Value fail cast to *model.EventStream")
			return
		}

		for i, event := range es.Events {
			if i == 0 {
				ob.startEventNames = append(ob.startEventNames, event.Name)
			}
			ob.eventMap[event.Name] = event
		}
		if msg.Msg == constant.AddEventStream {
			ob.isChanged = true
		}
	} else if msg.Msg == constant.ResetEventStream {
		/*for _, v := range ob.eventMap {
			messagebus.SendMsg(constant.JsonFileObserverName, model.JsonMsg{
				IsDeleteImage: true,
				ImageFileName: v.ImageFile,
			})
		}*/
		ob.startEventNames = []string{}
		ob.eventMap = map[string]*model.Event{}
		ob.isChanged = true
	} else if msg.Msg == "start" {
		if len(ob.startEventNames) == 0 ||
			len(ob.eventMap) == 0 {
			controller.OnStartBtnClick()
			return
		}
		var err error
		ob.eventMap, err = imageutil.CreateImageForEvents(ob.eventMap)
		if err != nil {
			fmt.Printf("fail to createImage for events:%+v\n", err)
			controller.OnStartBtnClick()
			return
		}

		for _, ev := range ob.eventMap {
			if ev.ImageFile == "" {
				ev.ImageFile = ev.Name + ".png"
			}
			messagebus.SendMsg(constant.JsonFileObserverName, model.JsonMsg{
				IsImage:       true,
				Image:         ev.Image,
				ImageFileName: ev.ImageFile,
			})
		}

		jsonMsg := model.JsonMsg{
			IsWrite:  true,
			NeedCopy: ob.isChanged,
			EventMsg: model.EventMsg{
				StartEventNames: ob.startEventNames,
				EventMap:        ob.eventMap,
			},
		}

		messagebus.SendMsg(constant.JsonFileObserverName, jsonMsg)

		//启动监控点击
		workMsg := model.WorkMsg{
			Msg: constant.StartState,
			EventMsg: model.EventMsg{
				StartEventNames: ob.startEventNames,
				EventMap:        ob.eventMap,
			},
		}
		messagebus.SendMsg(constant.WorkObserverName, workMsg)
	}
}
