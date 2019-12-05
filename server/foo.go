package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aptpod/opentracing-go-sandbox/lib"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type Foo struct {
}

func (h *Foo) Serve() error {
	// tracer の初期化
	tracer, closer, err := lib.CreateTracer("foo-service")
	if err != nil {
		return err
	}
	defer closer.Close()
	log.Println("start foo")
	return http.ListenAndServe(fmt.Sprintf(":%d", portMapping["foo"]), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("foo")
		// SpanのExtract
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))

		// リクエストからのSpanContextを引き継いで新しいSpanの開始
		span := tracer.StartSpan("get_foo", ext.RPCServerOption(spanCtx))
		defer span.Finish()
		// タグ付け
		ext.HTTPMethod.Set(span, r.Method)
		ext.HTTPUrl.Set(span, r.URL.String())

		<-time.After(time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "name" : "foo!" }`))

		// タグ付け
		ext.HTTPStatusCode.Set(span, http.StatusOK)
	}))
}
