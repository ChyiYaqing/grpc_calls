package main

import (
	"context"
	"encoding/json"

	users "github.com/chyiyaqing/grpc_calls/unary/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

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

	logger.Info("Client...")

	opts := grpc.WithInsecure()
	con, err := grpc.Dial("localhost:3000", opts)
	if err != nil {
		logger.Sugar().Fatalf("Error connecting: %v \n", err)
	}
	defer con.Close()
	c := users.NewUsersClient(con)
	getUsers(c, logger)
}

// getUsers function
func getUsers(c users.UsersClient, logger *zap.Logger) {
	req := &users.GetUsersReq{
		Status: users.UserStatus_USER_STATUS_UNKNOWN,
	}

	res, err := c.GetUsers(context.Background(), req)
	if err != nil {
		logger.Sugar().Fatalf("Error on GetUsers rpc call: %v \n", err)
	}
	logger.Sugar().Infof("Response: %+v\n", res)
}
