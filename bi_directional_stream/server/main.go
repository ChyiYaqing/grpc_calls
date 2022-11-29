package main

import (
	"encoding/json"
	"io"
	"net"

	users "github.com/chyiyaqing/grpc_calls/bi_directional_stream/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type server struct {
	users.UnimplementedUsersServer
}

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
	glog.Sugar().Info("Starting the server...")

	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		glog.Sugar().Fatalf("Unable to listen on port 3000: %v", err)
	}

	s := grpc.NewServer()
	users.RegisterUsersServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		glog.Sugar().Fatalf("Failed to serve: %v", err)
	}
}

// CreateUser function
func (*server) CreateUser(stream users.Users_CreateUserServer) error {
	glog.Sugar().Infof("CreateUser Function")

	for {
		// Receive the request and possible error from the stream object
		req, err := stream.Recv()

		// If there are no more requests, we return
		if err == io.EOF {
			return nil
		}

		// Handle error from the stream object
		if err != nil {
			glog.Sugar().Fatalf("Error when reading client request stream: %v", err)
		}

		// Get name, last name and user id, form the request
		name, lastName, userID := req.GetName(), req.GetLastName(), req.GetId()
		glog.Sugar().Infof("Request: name: %v, last_name: %v, id: %v", name, lastName, userID)

		// Initialize the errors and success variables
		errors := []string{}
		success := true

		// Run some validations
		if len(name) <= 3 {
			errors = append(errors, "Name is too short")
			success = false
		}

		if lastName == "Phill" {
			errors = append(errors, "Last Name already taken")
			success = false
		}

		// Build and send response to the client
		res := stream.Send(&users.CreateUserRes{
			UserId:  userID,
			Success: success,
			Errors:  errors,
		})

		// Handle any possible error, when sending the response
		if res != nil {
			glog.Sugar().Fatalf("Error when response was sent to the client: %v", res)
		}
	}
}
