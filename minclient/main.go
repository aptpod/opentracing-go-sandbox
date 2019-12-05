package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aptpod/opentracing-go-sandbox/lib"
	"github.com/opentracing/opentracing-go"
)

func main() {
	// Tracerの初期化
	closer, err := lib.InitGlobalTracer("client")
	if err != nil {
		panic(err)
	}
	defer closer.Close()

	// Tracerの取得
	tracer := opentracing.GlobalTracer()

	// Spanの開始
	// このSpanを引き回す
	span := tracer.StartSpan("get_hoge")

	// Spanの終了
	defer span.Finish()

	// Hogeの呼び出し！
	res, err := getHoge(span)
	if err != nil {
		panic(err)
	}
	log.Println(res)
}

func getHoge(span opentracing.Span) (string, error) {
	// リクエストの生成
	// 今回はHogeのPath
	url := "http://localhost:18080"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	// ※ Tracerを使ってSpanの情報をInject
	if err := span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	); err != nil {
		return "", err
	}

	// 通常のHTTPアクセス
	resp, err := Do(req)
	if err != nil {
		return "", err
	}

	return string(resp), nil
}

func Do(req *http.Request) ([]byte, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("error")
	}

	return ioutil.ReadAll(resp.Body)
}
