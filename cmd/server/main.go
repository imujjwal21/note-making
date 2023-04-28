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

	"notemaking/internal/httptransport"
	"notemaking/mongo"
	"notemaking/notes"
	"notemaking/users"

	"go.opentelemetry.io/otel"
)

type UriInformation struct {
	MongoUri string
	GrpcUri  string
}

func main() {

	ctx := context.Background()

	UriInformation, err := getAllUri()
	if err != nil {
		log.Printf("Uri information file reading error : %v", err)
	}

	client, err := mongo.ConnectDB(UriInformation.MongoUri)
	if err != nil {
		log.Printf("cannot connect to database : %v", err)
	}

	defer func() {
		err := client.Disconnect(ctx)
		if err != nil {
			panic(err)
		}
	}()

	coll := client.Database("notemaking")

	// kafka.CreateTopic()   sdktrace "go.opentelemetry.io/otel/sdk/trace"

	var port int

	flag.IntVar(&port, "port", 0, "Address ")
	flag.Parse()

	server := &http.Server{Handler: httptransport.NewHandler(users.NewInUserDatabase(coll), notes.NewInNoteDatabase(coll))}

	// start collector

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	shutdown, err := initProvider(ctx, UriInformation.GrpcUri)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	otel.Tracer("test-tracer")

	// end collector

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
