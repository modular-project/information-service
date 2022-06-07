package handler

import (
	"context"
	"information-service/controller"
	"information-service/model"

	pb "github.com/modular-project/protobuffers/information/product"
)

type ProductService struct {
	pb.UnimplementedProductServiceServer
}

func (p *ProductService) Create(c context.Context, in *pb.Product) (*pb.Response, error) {
	m := model.Product{
		Name:        in.Name,
		Price:       float64(in.Price),
		Description: in.Description,
		Url:         in.Url,
	}
	err := controller.CreateProduct(&m)
	if err != nil {
		return nil, err
	}
	return &pb.Response{Id: uint64(m.ID)}, nil
}

func (p *ProductService) Get(ctx context.Context, in *pb.RequestById) (*pb.Product, error) {
	m, err := controller.GetProductById(uint(in.Id))
	if err != nil {
		return &pb.Product{}, err
	}
	return &pb.Product{Name: m.Name, Url: m.Name, Description: m.Description, Price: float32(m.Price)}, err
}

func (p *ProductService) GetAll(ctx context.Context, in *pb.RequestGetAll) (*pb.ResponseGetAll, error) {
	ms, err := controller.GetAllProducts()
	if err != nil {
		return &pb.ResponseGetAll{}, err
	}
	products := make([]*pb.Product, len(ms))
	for i, m := range ms {
		products[i] = &pb.Product{Name: m.Name, Price: float32(m.Price), Url: m.Url, Description: m.Description}
	}
	return &pb.ResponseGetAll{Products: products}, nil
}

func (ps *ProductService) Update(ctx context.Context, in *pb.RequestUpdate) (*pb.Response, error) {
	p := in.Product
	new := model.Product{Name: p.Name, Description: p.Description, Price: float64(p.Price), Url: p.Url}
	err := controller.UpdateProductById(uint(in.Id), &new)
	if err != nil {
		return nil, err
	}
	return &pb.Response{Id: uint64(new.ID)}, err
}

func (p *ProductService) Delete(ctx context.Context, in *pb.RequestById) (*pb.ResponseDelete, error) {
	err := controller.DeleteProductById(uint(in.Id))
	if err != nil {
		return nil, err
	}
	return &pb.ResponseDelete{}, nil
}
