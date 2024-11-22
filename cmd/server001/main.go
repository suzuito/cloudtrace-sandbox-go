package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func main() {
	ctx := context.Background()

	// OpenTelemetryライブラリを用いた計装の初期化

	// Exporterは計装で必要となるデータを送る
	// Exporter自体はインターフェースであり、送り先毎に実装が異なる
	// デバッグ用途で使えるstdouttrace.Tracerというものもある
	// exporter, err := autoexport.NewSpanExporter(ctx)
	// exporter, err := stdouttrace.New(
	// 	stdouttrace.WithPrettyPrint(),
	// )
	exporter, err := texporter.New()
	if err != nil {
		panic(err)
	}

	// TraceProvider設定
	// ここでexporterやサンプリング方法を指定する
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		// サンプリング方法には様々ある
		trace.WithSampler(
			trace.AlwaysSample(),
		),
	)
	otel.SetTracerProvider(tp)
	defer tp.Shutdown(ctx)

	tracer := otel.Tracer("github.com/suzuito/cloudtrace-sandbox-go")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /hoge", handleGetHoge(tracer))

	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

func handleGetHoge(tracer oteltrace.Tracer) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		var span1 oteltrace.Span
		ctx, span1 = tracer.Start(ctx, fmt.Sprintf("%s %s", req.Method, req.URL.Path))
		defer span1.End()

		wg := sync.WaitGroup{}
		for range 50 {
			wg.Add(1)
			go func() {
				defer wg.Done()

				var span2 oteltrace.Span
				_, span2 = tracer.Start(ctx, "for loop")
				defer span2.End()

				span2.AddEvent("ev1")
				span2.SetStatus(codes.Error, "foo1")
				span2.AddLink(oteltrace.LinkFromContext(ctx))
				span2.SetAttributes(attribute.KeyValue{Key: "k1", Value: attribute.StringValue("v1")})
				span2.RecordError(fmt.Errorf("this is err1"))

				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			}()
		}
		wg.Wait()

		fmt.Println("hoge")
	}
}
