package lrucache

import "net/http"

type CacheStorer interface {
	Set(key string, value Payload)
	Get(key string) (Payload, bool)
	Clear(key string)
}

type Payload struct {
	Status int
	Header http.Header
	Data   string
}
