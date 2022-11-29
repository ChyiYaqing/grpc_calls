package main

import (
	"encoding/json"
	"io"
	"net"

	books "github.com/chyiyaqing/grpc_calls/client_stream/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type server struct {
	books.UnimplementedBooksServer
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
	defer glog.Sync()

	glog.Info("Starting server...")

	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		glog.Sugar().Fatalf("Unable to listen on port 3000: %v", err)
	}

	s := grpc.NewServer()
	books.RegisterBooksServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		glog.Sugar().Fatalf("Failed to serve: %v", err)
	}
}

// ValidateBooks function
func (*server) ValidateBooks(stream books.Books_ValidateBooksServer) error {
	glog.Info("Validate Books Function")

	// Initialize the ValidationError message
	errors := []*books.ValidationError{}
	for {
		// Start receiving stream from the client
		req, err := stream.Recv()

		// Check if the stream has finished
		if err == io.EOF {
			// Close the connection and return the response to the client
			return stream.SendAndClose(&books.ValidationRes{Errors: errors})
		}

		// Handle any possible errors while streaming requests
		if err != nil {
			glog.Sugar().Fatalf("Error when reading client request stream: %v", err)
		}

		// Get the title, pages and year fields from the req
		title := req.GetBook().GetTitle()
		pages := req.GetBook().GetPages()
		year := req.GetBook().GetYear()

		// Run some validations
		if len(title) <= 5 && pages < 300 && year < 2015 {
			// Create ValidationError object
			e := &books.ValidationError{
				BookId: req.GetBook().GetId(),
				Errors: []string{
					"Title must be at least 5 characters",
					"The book should have at least a minimum of 300 pages",
					"The year should be greated than 2015",
				},
			}
			// Append a new error message
			errors = append(errors, e)
		}
	}
}
