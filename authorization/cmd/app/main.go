package main

import (
	"authorization/internal/app/core"
	db "authorization/internal/app/database"
	grpc_service "authorization/internal/app/grpc"
	http_service "authorization/internal/app/http"
	"authorization/internal/app/pb"
	"authorization/internal/app/sessions"
	"authorization/internal/app/user"
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	port = ":9090"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDB(ctx)
	if err != nil {
		return
	}
	defer database.GetPool(ctx).Close()

	usersRepo := user.NewUsersRepo(user.NewRepository(database))
	sessionsRepo := sessions.NewSessionRepo(sessions.NewRepository(database))

	c := core.NewCore(usersRepo, sessionsRepo)
	r := gin.Default()
	s := http_service.NewHTTPService(r, c)

	lsn, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	go s.StartHTTPService(port)

	server := grpc.NewServer()
	pb.RegisterAuthenticationServiceServer(server, grpc_service.NewAuthenticationService(c))
	log.Printf("starting server on %s", lsn.Addr().String())
	if err := server.Serve(lsn); err != nil {
		log.Fatal(err)
	}
}
