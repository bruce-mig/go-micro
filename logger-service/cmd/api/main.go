package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/bruce-mig/go-micro/log-service/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRPCPort = "50001"
)

var client *mongo.Client

type Config struct {
	Session *scs.SessionManager
	Models  data.Models
}

func main() {
	// Connect to Mongo and get a client.
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection to Mongo when application exits
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// Set up application configuration with session and our Models type,
	// which allows us to interact with Mongo.
	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	app := Config{
		Session: session,
		Models:  data.New(client),
	}

	// Start webserver in its own GoRoutine
	go app.serve()

	// Start the gRPC server in its own GoRoutine
	go app.gRPCListen()

	// register the RPC server
	err = rpc.Register(new(RPCServer))
	if err != nil {
		return
	}
	app.rpcListen()

}

// serve starts the web server.
func (app *Config) serve() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	fmt.Println("--------------------------------------")
	fmt.Println("Starting logging web service on port", webPort)
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func (app *Config) rpcListen() error {
	log.Println("Starting RPC server on port,", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}
	defer listen.Close()

	// this loop executes forever, waiting for connections
	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		log.Println("Working...")
		go rpc.ServeConn(rpcConn)
	}
}

// Connect opens a connection to the Mongo database and returns a client.
func connectToMongo() (*mongo.Client, error) {
	// create connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// Connect to the MongoDB and return Client instance
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("mongo.Connect() ERROR:", err)
		return nil, err
	}

	log.Println("Connected to mongo!")

	return c, nil
}
