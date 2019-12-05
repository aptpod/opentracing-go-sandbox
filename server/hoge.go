package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aptpod/opentracing-go-sandbox/lib"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	tracelog "github.com/opentracing/opentracing-go/log"
)

type Hoge struct {
}

func (h *Hoge) Serve() error {
	// tracer の初期化
	tracer, closer, err := lib.CreateTracer("hoge-service")
	if err != nil {
		return err
	}
	defer closer.Close()
	cli := &ServiceClient{
		tracer: tracer,
	}
	log.Println("start hoge")
	// :18081でServe
	return http.ListenAndServe(fmt.Sprintf(":%d", portMapping["hoge"]), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("hoge")
		// SpanのExtract
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))

		// リクエストからのSpanContextを引き継いで新しいSpanの開始
		span := tracer.StartSpan("get_hoge", ext.RPCServerOption(spanCtx))
		defer span.Finish()

		// タグ付け
		ext.HTTPMethod.Set(span, r.Method)
		ext.HTTPUrl.Set(span, r.URL.String())
		span.SetTag("hoge_tag", "hoge_tag_value")

		// Key-Value形式で簡易指定
		span.LogKV(
			"hoge.log.key1", "hoge-log",
			"hoge.log.key2", "hoge-log2",
		)
		// 型情報ありの指定
		span.LogFields(
			tracelog.String("hoge.logfields.string", "hoge-log"),
			tracelog.Bool("hoge.logfields.bool", true),
		)

		// ctxにセット
		ctx := r.Context()
		ctx = opentracing.ContextWithSpan(ctx, span)

		// APIコール
		_, _ = cli.Call(ctx, "foo")
		_, _ = cli.Call(ctx, "fuga")

		// ダミー処理
		<-time.After(time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "name": "hoge!" }`))

		// タグ付け
		ext.HTTPStatusCode.Set(span, http.StatusOK)
	}))
}
