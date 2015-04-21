package notifier

var (
	AvailableNotifiers = make(map[string]func() interface{})
)

func RegisterNotifier(name string, factory func() interface{}) {
	AvailableNotifiers[name] = factory
}

type Plugin interface {
	Init([]string)
	NotifyUp(string)
	NotifyDown(string)
}
