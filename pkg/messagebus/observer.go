package messagebus

type Observer interface {
	OnEvent(interface{})
}