package main

import (
	"context"
	"encoding/json"
	"io"
	"time"

	users "github.com/chyiyaqing/grpc_calls/bi_directional_stream/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var glog *zap.Logger

func init() {
	rawJSON := []byte(`{
		"level": "debug",
		"encoding": "json",
		"outputPaths": ["stdout"],
		"errorOutputPaths": ["stdout"],
		"encoderConfig": {
			"messageKey": "message",
			"levelKey": "level",
			"levelEncoder": "lowercase"
		}
	}`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}

	glog, _ = cfg.Build()
	// if glog, err = cfg.Build(); err != nil {
	// 	panic(err)
	// }
}

func main() {
	glog.Sugar().Info("Client Stream...")

	opts := grpc.WithInsecure()
	con, err := grpc.Dial("localhost:3000", opts)
	if err != nil {
		glog.Sugar().Fatalf("Error connecting: %v", err)
	}

	defer con.Close()

	c := users.NewUsersClient(con)
	bulkUsers(c)
}

// container struct
type container struct {
	users []*users.User
}

// builkUsers function
func bulkUsers(c users.UsersClient) {
	// Get the stream and a possible error from the CreateUser function
	stream, err := c.CreateUser(context.Background())
	if err != nil {
		glog.Sugar().Fatalf("Error when getting stream object: %v", err)
		return
	}

	// Initialize the container struct and call the initUsers function to get user objects to send on the request message.
	requests := container{}.initUsers()

	// Create a new channel
	waitResponse := make(chan struct{})

	// Use a go routine to send request messages to the server
	go func() {
		// Iterate over the requets slice
		for _, req := range requests {
			// Send request message
			stream.Send(req)

			// Sleep for a little bit...
			time.Sleep(500 * time.Millisecond)
		}

		// Close stream
		stream.CloseSend()
	}()

	// Use a go routine to receive response messages from the server
	go func() {
		for {
			// Get response and possible error message from the stream
			res, err := stream.Recv()

			// Break for loop if there are no more response messages
			if err == io.EOF {
				break
			}

			// Handle a possible error
			if err != nil {
				glog.Sugar().Fatalf("Error when receiving response: %v", err)
			}

			// Log the response
			glog.Sugar().Infof("Server Response: %v", res)
		}

		// Close channel
		close(waitResponse)
	}()
	<-waitResponse
}

// initUsers function
func (c container) initUsers() []*users.User {
	c.users = append(c.users, c.getUser("1", "Carl", "Phill", 23))
	c.users = append(c.users, c.getUser("2", "Marisol", "Richardson", 29))
	c.users = append(c.users, c.getUser("3", "Mia", "Phill", 27))
	c.users = append(c.users, c.getUser("4", "Tomas", "Smith", 25))
	c.users = append(c.users, c.getUser("5", "Zian", "Heat", 28))
	return c.users
}

// getUser function
func (c container) getUser(id, name, lastName string, age int32) *users.User {
	return &users.User{
		Id:       id,
		Name:     name,
		LastName: lastName,
		Age:      age,
	}
}
