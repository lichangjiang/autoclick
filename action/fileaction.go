package action

import (
	"autoclick/constant"
	"autoclick/model"
	"autoclick/pkg/messagebus"
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func init() {
	ob := &jsonFileObserver{}
	messagebus.RegisterObserver(constant.JsonFileObserverName, ob)
	fmt.Println("jsonFileObserver init")
}

type jsonFileObserver struct {
}

func (ob *jsonFileObserver) OnEvent(ev interface{}) {
	msg, ok := ev.(model.JsonMsg)

	if !ok {
		fmt.Printf("jsonFileObserver Event fail cast to JsonMsg")
		return
	}

	if msg.IsWrite {
		if msg.EventMap != nil && msg.StartEventNames != nil {
			eg := &model.EventGroup{
				StartEvents: msg.StartEventNames,
				EndEvents:   []string{},
				Events:      []model.Event{},
			}

			for _, v := range msg.EventMap {
				eg.Events = append(eg.Events, *v)
			}
			dstr := time.Now().Format(time.RFC3339)
			if msg.NeedCopy {
				_, err := copy("event.json", "pre_event_"+dstr+".json")
				if err != nil {
					fmt.Printf("fail to copy event.json to pre_event.json:%+v\n", err)
				}
			}

			err := writeEventJson(eg)
			if err != nil {
				fmt.Printf("fail to write event.json:%+v\n", err)
			}
		} else {
			fmt.Println("jsonFileObserver got error data and fail to write json file")
		}
	}
	if msg.IsReadJson {
		eg, err := readEventJson()
		if err != nil {
			fmt.Printf("fail to read event.json %+v\n", err)
			return
		}

		if eg == nil {
			return
		}
		eventMap := map[string]model.Event{}
		for _, ev := range eg.Events {
			if ev.ImageFile != "" {
				img, err := readImage(filepath.Join("snapshot", ev.ImageFile))
				if err != nil {
					fmt.Printf("%+v\n", err)
				} else {
					ev.Image = img
				}
			}
			eventMap[ev.Name] = ev
		}

		for _, name := range eg.StartEvents {
			eventStream := &model.EventStream{
				Events: []*model.Event{},
			}
		LOOP:
			cev, ok := eventMap[name]
			if ok {
				eventStream.Events = append(eventStream.Events, &cev)
				if cev.NextEvent != "" {
					name = cev.NextEvent
					goto LOOP
				}
			} else {
				fmt.Printf("event name:%s\n not exist", name)
			}
			msg := model.EventStreamMsg{
				Msg:   constant.FileAddEventStream,
				Value: eventStream,
			}
			messagebus.SendMsg(constant.GlobalEventObserverName, msg)
		}
	}
	if msg.IsDir {
		dirPath := "snapshot"
		if dirPath != "" {
			dirNotExistCreate(dirPath)
		}
	}
	if msg.IsImage {
		fileName := msg.ImageFileName
		img := msg.Image
		if img != nil && fileName != "" {
			fullPath := filepath.Join("snapshot", fileName)
			file, err := os.Create(fullPath)
			defer file.Close()
			if err != nil {
				fmt.Printf("cannot create %s ,error:%+v\n", fullPath, file)
				return
			}
			err = png.Encode(file, img)
			if err != nil {
				fmt.Printf("cannot encode image %s,error:%+v\n", fullPath, file)
			}
		}
	}
}

func readEventJson() (*model.EventGroup, error) {
	if Exists("event.json") {
		eg := &model.EventGroup{}
		b, err := ioutil.ReadFile("event.json")
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(b, eg)
		if err != nil {
			return nil, err
		}
		return eg, nil
	}
	return nil, nil
}

func writeEventJson(eg *model.EventGroup) error {
	b, err := json.Marshal(eg)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("event.json", b, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func dirNotExistCreate(path string) error {
	if !Exists(path) {
		return os.Mkdir(path, os.ModePerm)
	}
	return nil
}

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func readImage(path string) (image.Image, error) {
	if Exists(path) {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		img, _, err := image.Decode(f)
		if err != nil {
			err = fmt.Errorf("%s decode error:%+v", path, err)
		}
		return img, err
	} else {
		return nil, fmt.Errorf("%s not exist", path)
	}
}

func copy(src, dst string) (int64, error) {
	if !Exists(src) {
		return 0, nil
	}
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	var destination *os.File
	if Exists(dst) {
		destination, err = os.Open(dst)
	} else {
		destination, err = os.Create(dst)
	}
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
