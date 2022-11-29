package main

import (
	"context"
	"encoding/json"
	"io"

	documents "github.com/chyiyaqing/grpc_calls/server_stream/proto"
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
	glog.Info("Client...")

	opts := grpc.WithInsecure()
	con, err := grpc.Dial("localhost:3000", opts)
	if err != nil {
		glog.Sugar().Fatalf("Error connecting: %v", err)
	}
	defer con.Close()

	c := documents.NewDocumentsClient(con)

	fetchDocuments(c)
}

// featchDocuments function
func fetchDocuments(c documents.DocumentsClient) {
	// Initialize request message
	req := &documents.EmptyReq{}

	// Get the stream and err
	stream, err := c.GetDocuments(context.Background(), req)
	if err != nil {
		glog.Sugar().Fatalf("Error on GetDocuments: %v", err)
	}

	for {
		// Start receiving streaming messages
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			glog.Sugar().Fatalf("error when receiving server response stream: %v", err)
		}
		glog.Sugar().Infof("Response from GetDocuments: %v", res.GetDocument())
	}
}
