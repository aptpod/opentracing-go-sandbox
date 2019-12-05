package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type ServiceClient struct {
	tracer opentracing.Tracer
}

func (s *ServiceClient) Call(ctx context.Context, serviceName string) (int, map[string]interface{}) {
	// ctxにセットされているSpanから新しく子のSpanを開始する
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, fmt.Sprintf("call_%s_%s", http.MethodGet, serviceName))
	defer span.Finish()
	serviceURL, _ := url.Parse(fmt.Sprintf("http://localhost:%d", portMapping[serviceName]))
	req, _ := http.NewRequest(http.MethodGet, serviceURL.String(), nil)

	// Tag付け
	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, serviceURL.String())
	ext.HTTPMethod.Set(span, "GET")

	// Tracerを使ってSpanの情報をInject
	_ = span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)

	resp, _ := http.DefaultClient.Do(req)

	defer resp.Body.Close()
	var responseJson map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&responseJson)
	return resp.StatusCode, responseJson
}
