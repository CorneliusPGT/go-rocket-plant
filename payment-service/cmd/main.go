package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"payment-service/grpc/paymentpb"
	"syscall"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50052

type PaymentService struct {
	paymentpb.UnimplementedPaymentServiceServer
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}

	s := grpc.NewServer()
	service := &PaymentService{}
	paymentpb.RegisterPaymentServiceServer(s, service)
	reflection.Register(s)

	go func() {
		log.Printf("gRPC listening on %d: ", grpcPort)
		err := s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("Server stopped")
}

func (s *PaymentService) PayOrder(ctx context.Context, req *paymentpb.PayOrderRequest) (*paymentpb.PayOrderResponse, error) {
	transaction_uuid := uuid.NewString()
	log.Printf("Заказ %s успешно оплачен с помощью %s пользователем %s\n transaction_uuid: %s", req.OrderUuid, req.PaymentMethod, req.UserUuid, transaction_uuid)
	return &paymentpb.PayOrderResponse{
		TransactionUuid: transaction_uuid,
	}, nil
}
