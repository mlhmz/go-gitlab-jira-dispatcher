package dispatcher

type Listener interface {
	Accept(event *Event)
}
