package domain

type Collector interface {
	Collect() ([]byte, error)
	Close() error
}
