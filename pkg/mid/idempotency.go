package mid

import (
	"fmt"
	"net/http"

	"github.com/muchlist/moneymagnet/pkg/lrucache"
	"github.com/muchlist/moneymagnet/pkg/web"
)

const idempotencyKey = "Idempotency"

type idempotent struct {
	cc lrucache.CacheStorer
}

func NewIdempotencyMiddleware(cache lrucache.CacheStorer) idempotent {
	return idempotent{cc: cache}
}

// IdempotentCheck is the instance middleware, it checks if the HTTP header has an idempotency key header
// and return the cached response when the same previous request already done processing.
// If the check failes to retrieve lock from the store, it will return `409 Conflict`.
func (o *idempotent) IdempotentCheck(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestKey := r.Header.Get(idempotencyKey)
		if requestKey == "" {
			h.ServeHTTP(w, r)
			return
		}

		// get cached data
		keyCache := fmt.Sprintf("%s-%s", r.URL.Path, requestKey)
		data, ok := o.cc.Get(keyCache)
		if ok {
			for k, v := range data.Header {
				for k1 := range v {
					w.Header().Add(k, v[k1])
				}
			}
			w.WriteHeader(data.Status)
			w.Write([]byte(data.Data))
			return
		}

		// create new response writer to support cache response body
		// use the current response writer if it's already an instance of responsewriter.ResponseWriter
		w1, ok := w.(*web.ResponseWriter)
		if !ok {
			w1 = web.NewResponseWritter(w)
		}
		h.ServeHTTP(w1, r)

		// cache data
		key := fmt.Sprintf("%s-%s", r.URL.Path, r.Header.Get(idempotencyKey))
		statusCode := w1.StatusCode()
		if statusCode >= http.StatusOK && statusCode < http.StatusBadRequest {
			o.cc.Set(key, lrucache.Payload{
				Status: statusCode,
				Header: w1.Header().Clone(),
				Data:   w1.Body(),
			})
		}
	})
}
