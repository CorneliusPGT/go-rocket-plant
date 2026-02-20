package repo

import (
	"context"
	"inventory-service/grpc/inventorypb"
	"inventory-service/internal/converter"
	"inventory-service/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PartRepo interface {
	Get(ctx context.Context, uuid string) (*inventorypb.Part, error)
	List(ctx context.Context, filter *inventorypb.PartsFilter) ([]*inventorypb.Part, error)
}

type MongoRepo struct {
	col *mongo.Collection
}

func NewMongoRepo(col *mongo.Collection) *MongoRepo {
	return &MongoRepo{
		col: col,
	}
}

func (r *MongoRepo) Get(ctx context.Context, uuid string) (*inventorypb.Part, error) {
	var part model.Part
	err := r.col.FindOne(ctx, bson.M{"uuid": uuid}).Decode(&part)
	if err != nil {
		return nil, err
	}

	return converter.ToProto(part), nil
}

func (r *MongoRepo) List(ctx context.Context, filter *inventorypb.PartsFilter) ([]*inventorypb.Part, error) {

	filterBson := bson.M{}

	if filter != nil {
		if len(filter.Uuids) > 0 {
			filterBson["uuid"] = bson.M{"$in": filter.Uuids}
		}

		if len(filter.Names) > 0 {
			filterBson["name"] = bson.M{"$in": filter.Names}
		}
		if len(filter.Categories) > 0 {
			cats := make([]int32, len(filter.Categories))
			for i, c := range filter.Categories {
				cats[i] = int32(c)
			}
			filterBson["category"] = bson.M{"$in": cats}
		}

		if len(filter.ManufacturerCountries) > 0 {
			filterBson["manufacter.country"] = bson.M{"$in": filter.ManufacturerCountries}
		}

		if len(filter.Tags) > 0 {
			filterBson["tags"] = bson.M{"$in": filter.Tags}
		}

	}

	cur, err := r.col.Find(ctx, filterBson)
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	var parts []*inventorypb.Part
	for cur.Next(ctx) {
		var p model.Part
		if err := cur.Decode(&p); err != nil {
			return nil, err
		}

		parts = append(parts, converter.ToProto(p))
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	if len(parts) == 0 {
		return nil, err
	}
	return parts, nil
}
