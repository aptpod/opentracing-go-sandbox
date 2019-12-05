package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aptpod/opentracing-sandbox/lib"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type Fuga struct {
}

func (h *Fuga) Serve() error {
	// tracer の初期化
	tracer, closer, err := lib.CreateTracer("fuga-service")
	if err != nil {
		return err
	}
	defer closer.Close()
	cli := &ServiceClient{
		tracer: tracer,
	}
	log.Println("start fuga")
	return http.ListenAndServe(fmt.Sprintf(":%d", portMapping["fuga"]), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("fuga")
		// SpanのExtract
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))

		// リクエストからのSpanContextを引き継いで新しいSpanの開始
		span := tracer.StartSpan("get_fuga", ext.RPCServerOption(spanCtx))
		defer span.Finish()
		// タグ付け
		ext.HTTPMethod.Set(span, r.Method)
		ext.HTTPUrl.Set(span, r.URL.String())

		// ctxにセット
		ctx := r.Context()
		ctx = opentracing.ContextWithSpan(ctx, span)

		// APIコール
		_, _ = cli.Call(ctx, "bar")
		_, _ = cli.Call(ctx, "baz")
		<-time.After(time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "name": fuga!" }`))

		// タグ付け
		ext.HTTPStatusCode.Set(span, http.StatusOK)
	}))
}
