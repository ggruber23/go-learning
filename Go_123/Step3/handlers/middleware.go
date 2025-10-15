package handlers

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
)

var traceIDGenerator int64

type traceId struct{}

func TraceIDMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		next := atomic.AddInt64(&traceIDGenerator, 1)
		ctx = context.WithValue(ctx, traceId{}, next) // add TraceID
		fmt.Println("Assigned TraceID: ", next)
		req = req.WithContext(ctx)
		handler.ServeHTTP(rw, req)
	})
}
