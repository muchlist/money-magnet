package mid

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/muchlist/moneymagnet/pkg/observ/mmetric"
	"github.com/muchlist/moneymagnet/pkg/web"
)

func EndpoitnCounter(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// create new response writer to support cache response body
		// use the current response writer if it's already an instance of responsewriter.ResponseWriter
		w1, ok := w.(*web.ResponseWriter)
		if !ok {
			w1 = web.NewResponseWritter(w)
		}
		t1 := time.Now()

		h.ServeHTTP(w1, r)

		// cache data
		statusCode := w1.StatusCode()
		path := fmt.Sprintf("%s_%s", r.Method, r.URL.Path)

		mmetric.AddEndpointHitCounter(context.Background(), statusCode, path)
		mmetric.AddLatencyPerPath(context.Background(), time.Since(t1).Microseconds(), path)

	})
}
