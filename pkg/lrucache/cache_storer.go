package lrucache

type CacheStorer interface {
	Set(key string, value Payload)
	Get(key string) (Payload, bool)
	Clear(key string)
}

type Payload struct {
	Status int
	Data   map[string]any
}
