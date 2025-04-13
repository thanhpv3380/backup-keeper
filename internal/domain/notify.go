package domain

type Notifier interface {
	Notify(msg string) error
}
