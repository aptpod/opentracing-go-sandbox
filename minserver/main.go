package main

import (
	"log"
	"net/http"
	"time"

	"github.com/aptpod/opentracing-go-sandbox/lib"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func main() {
	// tracer の初期化
	tracer, closer, err := lib.CreateTracer("hoge-service")
	if err != nil {
		panic(err)
	}
	defer closer.Close()
	log.Println("start hoge")
	http.ListenAndServe(":18080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("hoge")
		// SpanのExtract
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))

		// リクエストからのSpanContextを引き継いで新しいSpanの開始
		span := tracer.StartSpan("get_hoge", ext.RPCServerOption(spanCtx))
		defer span.Finish()
		// ダミー処理
		<-time.After(time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "name": "hoge!" }`))
	}))
}
