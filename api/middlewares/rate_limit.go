package middlewares

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

type RateLimitOpts struct {
	Period time.Duration
	Limit  int64
}

type RateLimit struct {
	Limiter        *limiter.Limiter
	OnError        stdlib.ErrorHandler
	OnLimitReached stdlib.LimitReachedHandler
}

func NewRateLimit(opts RateLimitOpts) *RateLimit {
	store := memory.NewStore()

	limiter := limiter.New(store, limiter.Rate{
		Period: opts.Period,
		Limit:  opts.Limit,
	})

	return &RateLimit{
		Limiter:        limiter,
		OnError:        stdlib.DefaultErrorHandler,
		OnLimitReached: stdlib.DefaultLimitReachedHandler,
	}
}

func (middleware *RateLimit) Handle(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := middleware.Limiter.GetIPKey(r)

		context, err := middleware.Limiter.Get(r.Context(), key)
		if err != nil {
			middleware.OnError(w, r, err)
			return
		}

		w.Header().Add("X-RateLimit-Limit", strconv.FormatInt(context.Limit, 10))
		w.Header().Add("X-RateLimit-Remaining", strconv.FormatInt(context.Remaining, 10))
		w.Header().Add("X-RateLimit-Reset", strconv.FormatInt(context.Reset, 10))

		if context.Reached {
			middleware.OnLimitReached(w, r)
			return
		}

		next.ServeHTTP(w, r)
	}
}
