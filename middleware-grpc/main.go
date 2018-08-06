package main

import (
	"net"

	"github.com/Sirupsen/logrus"
	"github.com/TykTechnologies/tyk-protobuf/bindings/go"
	"google.golang.org/grpc"
)

const (
	listenAddress          = ":9111"
	rabbitConnectionString = "amqp://guest:guest@rmq:5672/"
)

func main() {
	lis, err := net.Listen("tcp", listenAddress)
	fatalOnError(err, "failed to start tcp listener")

	logrus.Infof("starting grpc middleware on %s", listenAddress)
	s := grpc.NewServer()
	coprocess.RegisterDispatcherServer(s, &Dispatcher{})
	s.Serve(lis)

	//http.HandleFunc("/bundle.zip", func(w http.ResponseWriter, r *http.Request) {
	//	log.Println("received request for manifest")
	//	http.ServeFile(w, r, "bundle.zip")
	//})
	//
	//log.Printf("starting bundle manifest server on %v", ManifestAddress)
	//log.Fatal(http.ListenAndServe(ManifestAddress, nil))
}

func fatalOnError(err error, msg string) {
	if err != nil {
		logrus.WithError(err).Fatal(msg)
	}
}
