package main

import (
	"context"
	"encoding/json"
	"time"

	books "github.com/chyiyaqing/grpc_calls/client_stream/proto"
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
	glog.Info("Client Stream...")

	opts := grpc.WithInsecure()
	con, err := grpc.Dial("localhost:3000", opts)
	if err != nil {
		glog.Sugar().Fatalf("Error connecting: %v \n", err)
	}
	defer con.Close()

	c := books.NewBooksClient(con)
	validateBooks(c)
}

// container struct
type container struct {
	books []*books.ValidationReq
}

// validateBooks function
func validateBooks(c books.BooksClient) {
	// Initialize the container struct and call the initBooks function
	// to get dummy data to send on the request message.
	req := container{}.initBooks()

	// Get the stream and err
	stream, err := c.ValidateBooks(context.Background())
	if err != nil {
		glog.Sugar().Fatalf("Error on ValidateBooks: %v", err)
	}

	// Iterate over the request message
	for _, v := range req {
		// Start making streaming requests by sending each book object inside the request message
		glog.Sugar().Info("Client streaming request: %v", v)
		stream.Send(v)
		time.Sleep(500 * time.Microsecond)
	}

	// Once the for loop finishes, the stream is closed and get the response and a potential error
	res, err := stream.CloseAndRecv()
	if err != nil {
		glog.Sugar().Fatalf("Error when closing the stream and receiving the response: %v", err)
	}

	// Print the response errors message
	glog.Sugar().Infof("Validation errors: %v \n", res.Errors)
}

// initBooks function
func (c container) initBooks() []*books.ValidationReq {
	c.books = append(c.books, c.getBook("1", "Book 1", "This is a really good book about history", "John Phill", "2.5", 395, 2010))
	c.books = append(c.books, c.getBook("2", "Book 2", "Improve your communication skills", "Carl Matz", "3.5", 425, 2008))
	c.books = append(c.books, c.getBook("3", "Book 3", "Movies and TV shows", "Carl Matz", "2.9", 225, 2015))
	c.books = append(c.books, c.getBook("3", "Bo", "Cookies", "Carl Matz", "2.2", 275, 2011))
	c.books = append(c.books, c.getBook("4", "Book 5", "Learn more about animals and nature", "John Phill", "2.7", 455, 2018))
	c.books = append(c.books, c.getBook("4", "Book 6 with long title", "Machine learning", "John Phill", "1.0", 375, 2016))
	c.books = append(c.books, c.getBook("5", "Book", "10 good ideas for decorating your house", "Maty Metzer", "3.2", 250, 2007))
	return c.books

}

// getBook function
func (c container) getBook(id, title, desc, author, edition string, pages, year int64) *books.ValidationReq {
	return &books.ValidationReq{
		Book: &books.Book{
			Id:          id,
			Title:       title,
			Description: desc,
			Pages:       pages,
			Author:      author,
			Year:        year,
			Edition:     edition,
		},
	}
}
