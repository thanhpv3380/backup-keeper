package domain

type Collector interface {
	Collect() (interface{}, error)
	Close() error
}
