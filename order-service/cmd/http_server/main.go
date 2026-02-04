package main

import (
	"context"
	"errors"
	"inventory-service/grpc/inventorypb"
	"log"
	"net/http"
	inventorygrpc "order-service/cmd/grpc/inventory"
	paymentgrpc "order-service/cmd/grpc/payment"
	"order-service/internal/handlers"
	api "order-service/internal/oapi"
	repository "order-service/internal/repository/order"
	"order-service/internal/service/order"
	"payment-service/grpc/paymentpb"

	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	const (
		httpPort          = ":8080"
		readHeaderTimeout = 5 * time.Second
		shutdownTimeout   = 10 * time.Second
	)
	repo := repository.NewRepository()

	payAddr := os.Getenv("PAYMENT_SERVICE_ADDR")
	if payAddr == "" {
		payAddr = "127.0.0.1:50052"
	}

	invAddr := os.Getenv("INVENTORY_SERVICE_ADDR")
	if invAddr == "" {
		invAddr = "127.0.0.1:50051"
	}
	invConn, err := grpc.NewClient(invAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer invConn.Close()

	payConn, err := grpc.NewClient(payAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer payConn.Close()

	invService := inventorygrpc.New(inventorypb.NewInventoryServiceClient(invConn))
	payService := paymentgrpc.New(paymentpb.NewPaymentServiceClient(payConn))
	orderService := order.NewService(repo, invService, payService)
	handler := &handlers.OrderHandler{
		Service: orderService,
	}
	server, err := api.NewServer(handler)
	if err != nil {
		log.Fatalf("не удалось создать сервер: %v", err)
	}

	httpServer := &http.Server{
		Addr:    httpPort,
		Handler: server,
	}
	go func() {
		log.Println("сервер запущен на порту " + httpPort)
		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("ошибка запуска сервера: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Завершение программы OrderService...")
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = httpServer.Shutdown(ctx)
	if err != nil {
		log.Printf("ошибка при остановке сервера: %v\n", err)
	}
	log.Println("Сервер остановлен")
}
