package converter

import (
	"inventory-service/grpc/inventorypb"
	"inventory-service/internal/model"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToProto(p model.Part) *inventorypb.Part {

	var dimensions *inventorypb.Dimensions
	if p.Dimensions != nil {
		dimensions = &inventorypb.Dimensions{
			Length: p.Dimensions.Length,
			Width:  p.Dimensions.Width,
			Height: p.Dimensions.Height,
			Weight: p.Dimensions.Weight,
		}
	}

	var manuf *inventorypb.Manufacter
	if p.Manufacter != nil {
		manuf = &inventorypb.Manufacter{
			Name:    p.Manufacter.Name,
			Country: p.Manufacter.Country,
			Website: p.Manufacter.Website,
		}
	}

	return &inventorypb.Part{
		Uuid:          p.UUID,
		Name:          p.Name,
		Description:   p.Description,
		Category:      inventorypb.Category(p.Category),
		Price:         p.Price,
		StockQuantity: p.StockQuantity,
		Dimensions:    dimensions,
		Manufacter:    manuf,
		Tags:          p.Tags,
		CreatedAt:     timestamppb.New(p.CreatedAt),
		UpdatedAt:     timestamppb.New(p.UpdatedAt),
	}
}
