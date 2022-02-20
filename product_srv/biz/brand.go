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

func (p ProductServer) CreateBrand(ctx context.Context, req *pb.BrandItemReq) (*pb.BrandItemRes, error) {
	var brand model.Brand
	r := internal.DB.Find("name=? and logo=?", req.Name, req.Logo)
	if r.RowsAffected > 0 {
		return nil, errors.New(custom_error.BrandAlreadyExits)
	}
	brand.Name = req.Name
	brand.Logo = req.Logo

	internal.DB.Save(&brand)
	return ConvertBrandModel2Pb(brand), nil
}

func (p ProductServer) BrandList(ctx context.Context, req *pb.BrandPagingReq) (*pb.BrandRes, error) {
	var brandList []model.Brand
	var brands []*pb.BrandItemRes
	var brandRes pb.BrandRes

	//
	//r:=internal.DB.Find(&brandList)
	//fmt.Println(r.RowsAffected)
	//for _,item:=range brandList{
	//	brands = append(brands,ConvertBrandModel2Pb(item))
	//}
	//brandRes.ItemList=brands
	//brandRes.Total = int32(r.RowsAffected)
	//
	//return &brandRes,nil

	//var count int64
	//skip:=(req.PageNo-1)*req.PageSize
	//r:=internal.DB.Model(&model.Brand{}).Count(&count).Offset(int(skip)).Limit(int(req.PageSize)).Find(&brandList)
	//if r.RowsAffected<1{
	//	//TODO,可以进行进一步判断，根据业务需求
	//}
	//brandRes.Total=int32(count)
	//for _,item:=range brandList{
	//	brands = append(brands,ConvertBrandModel2Pb(item))
	//}
	//brandRes.ItemList=brands
	//return &brandRes,nil

	_ = internal.DB.Scopes(internal.MyPaging(int(req.PageNo), int(req.PageSize))).Find(&brandList)
	for _, item := range brandList {
		brands = append(brands, ConvertBrandModel2Pb(item))
	}
	brandRes.ItemList = brands
	var count int64
	internal.DB.Model(&model.Brand{}).Count(&count)
	brandRes.Total = int32(count)
	return &brandRes, nil
}

func (p ProductServer) DeleteBrand(ctx context.Context, req *pb.BrandItemReq) (*emptypb.Empty, error) {
	r := internal.DB.Delete(&model.Brand{}, req.Id)
	if r.Error != nil {
		return nil, errors.New(custom_error.DelBrandFail)
	}
	return &emptypb.Empty{}, nil
}

func (p ProductServer) UpdateBrand(ctx context.Context, req *pb.BrandItemReq) (*emptypb.Empty, error) {
	var brand model.Brand
	r := internal.DB.First(&brand, req.Id)
	if r.RowsAffected < 1 {
		return nil, errors.New(custom_error.BrandNotExits)
	}
	if req.Name != "" {
		brand.Name = req.Name
	}
	if req.Logo != "" {
		brand.Logo = req.Logo
	}
	internal.DB.Save(&brand)
	return &emptypb.Empty{}, nil
}

func ConvertBrandModel2Pb(item model.Brand) *pb.BrandItemRes {
	brand := &pb.BrandItemRes{
		Name: item.Name,
		Logo: item.Logo,
	}
	if item.ID > 0 {
		brand.Id = item.ID
	}
	return brand
}
