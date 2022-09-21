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
		BaseID:      uint(in.BaseId),
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
	return &pb.Product{Id: uint64(m.ID), Name: m.Name, Url: m.Url, Description: m.Description, Price: float32(m.Price), BaseId: uint64(m.BaseID)}, err
}

func (p *ProductService) GetAll(ctx context.Context, in *pb.RequestGetAll) (*pb.ResponseGetAll, error) {
	ms, err := controller.GetAllProducts()
	if err != nil {
		return &pb.ResponseGetAll{}, err
	}
	products := make([]*pb.Product, len(ms))
	for i, m := range ms {
		products[i] = &pb.Product{Id: uint64(m.ID), Name: m.Name, Url: m.Url, Description: m.Description, Price: float32(m.Price), BaseId: uint64(m.BaseID)}
	}
	return &pb.ResponseGetAll{Products: products}, nil
}

func (ps *ProductService) Update(ctx context.Context, in *pb.RequestUpdate) (*pb.Response, error) {
	p := in.Product
	new := model.Product{Name: p.Name, Description: p.Description, Price: float64(p.Price), Url: p.Url, BaseID: uint(p.BaseId)}
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

func (p *ProductService) GetInBatch(ctx context.Context, in *pb.RequestGetInBatch) (*pb.ResponseGetAll, error) {
	ms, err := controller.GetProductsInBatch(in.Ids)
	if err != nil {
		return nil, err
	}
	if ms == nil {
		return nil, nil
	}
	products := make([]*pb.Product, len(ms))
	for i, m := range ms {
		products[i] = &pb.Product{Id: uint64(m.ID), Name: m.Name, Url: m.Url, Description: m.Description, Price: float32(m.Price), BaseId: uint64(m.BaseID)}
	}
	return &pb.ResponseGetAll{Products: products}, nil
}

func (p *ProductService) GetByBase(ctx context.Context, in *pb.RequestGetByBase) (*pb.ResponseGetAll, error) {
	ms, err := controller.GetProductsByBase(in.Base)
	if err != nil {
		return nil, err
	}
	if ms == nil {
		return nil, nil
	}
	products := make([]*pb.Product, len(ms))
	for i, m := range ms {
		products[i] = &pb.Product{Id: uint64(m.ID), Name: m.Name, Url: m.Url, Description: m.Description, Price: float32(m.Price), BaseId: uint64(m.BaseID)}
	}
	return &pb.ResponseGetAll{Products: products}, nil
}
