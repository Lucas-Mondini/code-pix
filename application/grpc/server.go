package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/Lucas-Mondini/code-pix/application/grpc/pb"
	"github.com/Lucas-Mondini/code-pix/application/usecase"
	"github.com/Lucas-Mondini/code-pix/infrastructure/repository"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartGrpcServer(database *gorm.DB, port int) {
	grpcServer := grpc.NewServer()

	reflection.Register(grpcServer)

	pixRepository := repository.PixKeyRepositoryDB{DB: database}
	pixUseCase := usecase.PixUseCase{PixKeyRepository: &pixRepository}
	pixGrpcService := NewPixGrpcService(pixUseCase)
	pb.RegisterPixServiceServer(grpcServer, pixGrpcService)

	address := fmt.Sprintf("0.0.0.0:%d", port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Cannot start GRPC server", err)
	}

	log.Printf("GRPC server has been started at port %d", port)
	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatal("Cannot start GRPC server", err)
	}
}
