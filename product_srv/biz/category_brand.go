package biz

import (
	"context"
	"errors"
	"github.com/i-coder-robot/mic-trainning-lessons-part2/custom_error"
	"github.com/i-coder-robot/mic-trainning-lessons-part2/internal"
	"github.com/i-coder-robot/mic-trainning-lessons-part2/model"
	"github.com/i-coder-robot/mic-trainning-lessons-part2/proto/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (p ProductServer) CreateCategoryBrand(ctx context.Context, req *pb.CategoryBrandReq) (*pb.CategoryBrandRes, error) {
	var res pb.CategoryBrandRes
	var item model.ProductCategoryBrand
	var category model.Category
	var brand model.Brand
	//分类判断
	r := internal.DB.First(&category, req.CategoryId)
	if r.RowsAffected < 1 {
		return nil, errors.New(custom_error.CategoryNotExists)
	}
	//品牌判断
	r = internal.DB.First(&brand, req.BrandId)
	if r.RowsAffected < 1 {
		return nil, errors.New(custom_error.BrandNotExits)
	}
	//是否已经存在关系
	item.CategoryID = req.CategoryId
	item.BrandID = req.BrandId
	internal.DB.Save(&item)
	res.Id = item.ID
	return &res, nil
}

func (p ProductServer) DeleteCategoryBrand(ctx context.Context, req *pb.CategoryBrandReq) (*emptypb.Empty, error) {
	r := internal.DB.Delete(&model.ProductCategoryBrand{}, req.Id)
	if r.RowsAffected < 1 {
		return nil, errors.New(custom_error.DelProductCategoryBrandFailed)
	}
	return &emptypb.Empty{}, nil
}

func (p ProductServer) UpdateCategoryBrand(ctx context.Context, req *pb.CategoryBrandReq) (*emptypb.Empty, error) {
	//分类判断
	//品牌判断

	var productCategoryBrand model.ProductCategoryBrand
	r := internal.DB.First(&productCategoryBrand, req.Id)
	if r.RowsAffected < 1 {
		//输出错误
		return nil, errors.New(custom_error.ProductCategoryBrandNotFound)
	}
	productCategoryBrand.CategoryID = req.CategoryId
	productCategoryBrand.BrandID = req.BrandId
	internal.DB.Save(&productCategoryBrand)
	return &emptypb.Empty{}, nil
}

func (p ProductServer) CategoryBrandList(ctx context.Context, req *pb.PagingReq) (*pb.CategoryBrandListRes, error) {
	//各种逻辑判断
	var items []model.ProductCategoryBrand
	var resList []*pb.CategoryBrandRes
	var count int64
	var res pb.CategoryBrandListRes
	internal.DB.Model(&model.ProductCategoryBrand{}).Count(&count)
	internal.DB.Preload("Category").Preload("Brand").Scopes(internal.MyPaging(int(req.PageNo), int(req.PageSize))).Find(&items)
	for _, item := range items {
		pcb := ConvertProductCategoryBrand2Pb(item)
		resList = append(resList, pcb)
	}
	res.Total = int32(count)
	res.ItemList = resList
	return &res, nil
}

func (p ProductServer) GetCategoryBrandList(ctx context.Context, req *pb.CategoryItemReq) (*pb.BrandRes, error) {
	var res pb.BrandRes
	var category model.Category
	var itemList []model.ProductCategoryBrand
	var itemListRes []*pb.BrandItemRes

	r := internal.DB.First(&category, req.Id)
	if r.RowsAffected == 0 {
		return nil, errors.New(custom_error.CategoryNotExists)
	}
	r = internal.DB.Preload("Brand").Where(&model.ProductCategoryBrand{CategoryID: req.ParentCategoryId}).Find(&itemList)
	if r.RowsAffected > 0 {
		res.Total = int32(r.RowsAffected)
	}
	for _, item := range itemList {
		itemListRes = append(itemListRes, &pb.BrandItemRes{
			Id:   item.Brand.ID,
			Name: item.Brand.Name,
			Logo: item.Brand.Logo,
		})
	}
	res.ItemList = itemListRes
	return &res, nil
}

func ConvertProductCategoryBrand2Pb(pcb model.ProductCategoryBrand) *pb.CategoryBrandRes {
	cb := pb.CategoryBrandRes{
		Id: pcb.ID,
		Brand: &pb.BrandItemRes{
			Id:   pcb.Brand.ID,
			Name: pcb.Brand.Name,
			Logo: pcb.Brand.Logo,
		},
		Category: &pb.CategoryItemRes{
			Id:               pcb.Category.ID,
			Name:             pcb.Category.Name,
			ParentCategoryId: pcb.Category.ParentCategoryID,
			Level:            pcb.Category.Level,
		},
	}
	return &cb
}
