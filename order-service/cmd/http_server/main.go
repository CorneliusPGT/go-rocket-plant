package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	client "order-service/cmd/http_client"
	"order-service/internal/handlers"
	"order-service/internal/models"
	api "order-service/internal/oapi"
	"order-service/internal/service"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	const (
		httpPort          = ":8080"
		readHeaderTimeout = 5 * time.Second
		shutdownTimeout   = 10 * time.Second
	)
	paymentAddr := os.Getenv("PAYMENT_SERVICE_ADDR")
	if paymentAddr == "" {
		paymentAddr = "127.0.0.1:50052"
	}
	paymentClient, err := client.NewPaymentClient(paymentAddr)
	if err != nil {
		log.Fatalf("Не удалось создать gRPC-клиент PaymentService: %v", err)
	}
	defer paymentClient.Close()

	realPayment := service.NewPaymentServiceClient(paymentClient)

	inventAddr := os.Getenv("INVENTORY_SERVICE_ADDR")
	if inventAddr == "" {
		inventAddr = "127.0.0.1:50051"
	}
	inventoryClient, err := client.NewInventoryClient(inventAddr)
	if err != nil {
		log.Fatalf("Не удалось создать gRPC-клиент PaymentService: %v", err)
	}
	defer inventoryClient.Close()

	realInventory := service.NewInventoryServiceClient(inventoryClient)

	storage := models.NewOrderStorage()
	orderService := service.NewOrderService(storage, *realInventory, *realPayment)
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
