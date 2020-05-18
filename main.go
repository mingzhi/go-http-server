// Package main contains a basic HTTP server using the net/http package.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	defaultName = "world"
)

var port = flag.String("port", ":8888", "http server port number. Default to :8888.")
var grpcAddress = flag.String("grpc", "127.0.0.1:50051", "gRPC helle world server address")
var client pb.GreeterClient

func main() {
	log.Println("go-http-server starting...")
	// Set up a connection to the server.
	conn, err := grpc.Dial(*grpcAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("did not connect: %v", err)
	}
	defer conn.Close()
	client = pb.NewGreeterClient(conn)
	log.Println("go-http-server connected to grpc server.")

	http.HandleFunc("/showHeaders", showHeaders)
	http.HandleFunc("/sayHello", sayHello)
	http.ListenAndServe(*port, nil)
}

func showHeaders(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func sayHello(w http.ResponseWriter, req *http.Request) {
	name := defaultName
	if len(req.URL.Query()["name"]) > 0 {
		name = req.URL.Query()["name"][0]
	}
	message, err := getMessage(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "%s", message)
}

func getMessage(name string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		return "", err
	}
	return r.GetMessage(), nil
}
