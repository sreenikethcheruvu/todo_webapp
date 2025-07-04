package main

import (
    "log"
    "net"

    "google.golang.org/grpc"
    pb "todoapp-backend/src/pb"
    "todoapp-backend/src"
)

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()

    pb.RegisterTodoServiceServer(grpcServer, &src.TodoServiceServer{})

    log.Println("🚀 gRPC server listening on port 50051...")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
