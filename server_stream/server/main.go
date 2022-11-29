package main

import (
	"encoding/json"
	"net"
	"time"

	documents "github.com/chyiyaqing/grpc_calls/server_stream/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type server struct {
	documents.UnimplementedDocumentsServer
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
	glog.Sugar().Info("Starting server...")

	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		glog.Sugar().Fatalf("Unable to listen on port 3000: %v", err)
	}

	s := grpc.NewServer()
	documents.RegisterDocumentsServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		glog.Sugar().Fatalf("Failed to serve: %v", err)
	}
}

// container struct
type container struct {
	documents []*documents.Document
}

// GetDocuments function
func (*server) GetDocuments(req *documents.EmptyReq, stream documents.Documents_GetDocumentsServer) error {
	glog.Info("GetDocuments function")

	// Initialize the container struct and call the initDocuments function to get dummy data to send on the stream response message.
	docs := container{}.initDocuments()

	// Iterate over the documents
	for _, v := range docs {
		// Run some validation on each object
		if v.Size > 250 {
			// Create the response object
			res := &documents.GetDocumentsRes{
				Document: v,
			}

			// Use the stream object to send the response stream message
			stream.Send(res)

			// Sleep for a little bit...
			time.Sleep(500 * time.Microsecond)
		}
	}
	return nil
}

// initDocuments function
func (c container) initDocuments() []*documents.Document {
	c.documents = append(c.documents, c.getDocument("Doc One", "nat", 345))
	c.documents = append(c.documents, c.getDocument("Doc Two", "zip", 245))
	c.documents = append(c.documents, c.getDocument("Doc Three", "nat", 445))
	c.documents = append(c.documents, c.getDocument("Doc Four", "pid", 545))
	c.documents = append(c.documents, c.getDocument("Doc Five", "nat", 145))
	return c.documents
}

// getDocument function
func (c container) getDocument(name, documentType string, size int64) *documents.Document {
	return &documents.Document{
		Name:         name,
		DocumentType: documentType,
		Size:         size,
	}
}
