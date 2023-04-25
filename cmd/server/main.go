package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"notemaking/internal/httptransport"
	"notemaking/mongo"
	"notemaking/notes"
	"notemaking/users"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func initProvider(ctx context.Context) (func(context.Context) error, error) {

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("test-service"),
			semconv.ServiceNameKey.String("notemaking"),
			semconv.ServiceVersion("1.9"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, "127.0.0.1:4317",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tracerProvider.Shutdown, nil
}

func main() {

	client := mongo.ConnectDB()

	defer func() {
		err := client.Disconnect(context.TODO())
		if err != nil {
			panic(err)
		}
	}()

	coll := client.Database("notemaking")

	// kafka.CreateTopic()

	var port int

	flag.IntVar(&port, "port", 0, "Address ")
	flag.Parse()

	server := &http.Server{Handler: httptransport.NewHandler(users.NewInUserDatabase(coll), notes.NewInNoteDatabase(coll))}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutdown, err := initProvider(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	otel.Tracer("test-tracer")

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

		if err != nil {
			log.Panicf("can't create Listener : %v \n", err)
		}

		log.Printf(" starting http server on %q", lis.Addr())

		if err := server.Serve(lis); err != nil {
			log.Panicf("can't start server : %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	log.Printf("\nExit signal : %q", <-sig)
}
