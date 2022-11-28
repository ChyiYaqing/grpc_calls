package main

import (
	"context"
	"encoding/json"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	users "github.com/chyiyaqing/grpc_calls/unary/proto"
)

type server struct {
	users.UnimplementedUsersServer
}

func main() {
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

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Info("Starting server...")

	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		logger.Sugar().Fatalf("Unable to listen on port 3000: %v", err)
	}

	s := grpc.NewServer()
	users.RegisterUsersServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		logger.Sugar().Fatalf("Fails to serve: %v", err)
	}
}

// GetUsers function
func (*server) GetUsers(ctx context.Context, req *users.GetUsersReq) (*users.GetUsersRes, error) {
	status := req.GetStatus()
	userList := getUserList()
	usersFiltered := []*users.User{}
	switch status {
	case users.UserStatus_USER_STATUS_ACTIVE:
		usersFiltered = filterBy("active", userList)
	case users.UserStatus_USER_STATUS_BLOCKED:
		usersFiltered = filterBy("blocked", userList)
	case users.UserStatus_USER_STATUS_SUSPENDED:
		usersFiltered = filterBy("suspended", userList)
	default:
		usersFiltered = userList
	}

	res := users.GetUsersRes{
		Users: usersFiltered,
	}
	return &res, nil
}

// getUserList function
func getUserList() []*users.User {
	userObj := []*users.User{}
	userObj = append(userObj, &users.User{Name: "John", LastName: "Phill", Age: 34, Email: "john@gmail.com", Status: "active"})
	userObj = append(userObj, &users.User{Name: "Carl", LastName: "Meertz", Age: 23, Email: "carl@gmail.com", Status: "active"})
	userObj = append(userObj, &users.User{Name: "Sysan", LastName: "Zeanz", Age: 30, Email: "susan@gmail.com", Status: "blocked"})
	userObj = append(userObj, &users.User{Name: "Marylen", LastName: "Inc", Age: 29, Email: "marylen@gmail.com", Status: "blocked"})
	userObj = append(userObj, &users.User{Name: "Peet", LastName: "Green", Age: 25, Email: "peet@gmail.com", Status: "ignored"})
	userObj = append(userObj, &users.User{Name: "Maty", LastName: "Jackson", Age: 28, Email: "maty@gmail.com", Status: "suspended"})
	return userObj
}

// FilterBy function
func filterBy(status string, userList []*users.User) []*users.User {
	usersFiltered := []*users.User{}
	for _, v := range userList {
		if (v.Status == "blocked" || v.Status == "ignored") && status == "blocked" {
			usersFiltered = append(usersFiltered, v)
		} else if v.Status == status {
			usersFiltered = append(usersFiltered, v)
		}
	}
	return usersFiltered
}
