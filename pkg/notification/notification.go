package notification

type Notificator interface {
	Name() string
	Notify(title string, message string) error
}
