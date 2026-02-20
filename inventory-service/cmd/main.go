package main

import (
	"context"
	"fmt"
	"inventory-service/grpc/handlers"
	"inventory-service/grpc/inventorypb"
	"inventory-service/internal/model"
	"inventory-service/internal/service"
	repo "inventory-service/repository"

	"time"

	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50051

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("mongo not reachable:", err)
	}

	db := client.Database("inventory")
	col := db.Collection("parts")
	repo := repo.NewMongoRepo(col)
	partService := service.NewPartService(repo)

	err = seedData(ctx, col)
	if err != nil {
		log.Println("did not seed:", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}

	s := grpc.NewServer()
	handler := handlers.NewInventoryHandler(partService)
	inventorypb.RegisterInventoryServiceServer(s, handler)
	reflection.Register(s)

	go func() {
		log.Printf("gRPC server listening on %d\n", grpcPort)
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

func seedData(ctx context.Context, col *mongo.Collection) error {
	count, err := col.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}
	if count > 0 {
		log.Println("Data already exists, seed skipped...")
		return nil
	}
	log.Println("seeding initial data...")

	parts := []interface{}{
		bson.M{
			"uuid":           "engine-1",
			"name":           "Main Engine",
			"description":    "Primary propulsion engine",
			"price":          1500000,
			"stock_quantity": 10,
			"category":       inventorypb.Category_CATEGORY_ENGINE,
			"dimensions": &inventorypb.Dimensions{
				Length: 4,
				Width:  2,
				Height: 2,
				Weight: 1500,
			},
			"manufacter": model.Manufacter{
				Name:    "SpaceY",
				Country: "USA",
				Website: "https://spacey.example",
			},
			"tags":       []string{"engine", "rocket"},
			"created_at": time.Now(),
			"updated_at": time.Now(),
		},
		bson.M{
			"uuid":           "wing-1",
			"name":           "Left Wing",
			"description":    "Aerodynamic wing",
			"price":          250000,
			"stock_quantity": 5,
			"category":       inventorypb.Category_CATEGORY_WING,
			"manufacter": model.Manufacter{
				Name:    "AeroWorks",
				Country: "Germany",
			},
			"tags":       []string{"wing"},
			"created_at": time.Now(),
			"updated_at": time.Now(),
		},
	}

	_, err = col.InsertMany(ctx, parts)
	return err
}
