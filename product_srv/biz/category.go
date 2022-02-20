package biz

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/i-coder-robot/mic-trainning-lessons-part2/custom_error"
	"github.com/i-coder-robot/mic-trainning-lessons-part2/internal"
	"github.com/i-coder-robot/mic-trainning-lessons-part2/model"
	"github.com/i-coder-robot/mic-trainning-lessons-part2/proto/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (p ProductServer) CreateCategory(ctx context.Context, req *pb.CategoryItemReq) (*pb.CategoryItemRes, error) {
	category := model.Category{}
	//TODO 业务逻辑判断
	category.Name = req.Name
	category.Level = req.Level
	if category.Level > 0 {
		category.ParentCategoryID = req.ParentCategoryId
	}
	r := internal.DB.Save(&category)
	if r.RowsAffected < 1 {
		return nil, r.Error
	}
	res := ConvertCategoryModel2Pb(category)
	return res, nil
}

func (p ProductServer) DeleteCategory(ctx context.Context, req *pb.CategoryDelReq) (*emptypb.Empty, error) {
	internal.DB.Delete(&model.Category{}, req.Id)
	//TODO 逻辑判断
	//如果你删除的是一级分类，那么下面的2级，3级，也要删掉
	return &emptypb.Empty{}, nil
}

func (p ProductServer) UpdateCategory(ctx context.Context, req *pb.CategoryItemReq) (*emptypb.Empty, error) {
	var category model.Category
	r := internal.DB.Find(&category, req.Id)
	if r.RowsAffected < 1 {
		return nil, errors.New(custom_error.CategoryNotExists)
	}
	if req.Name != "" {
		category.Name = req.Name
	}
	if req.ParentCategoryId > 0 {
		category.ParentCategoryID = req.ParentCategoryId
	}
	if req.Level > 0 {
		category.Level = req.Level
	}
	internal.DB.Save(&category)
	return &emptypb.Empty{}, nil
}

func (p ProductServer) GetAllCategoryList(ctx context.Context, empty *emptypb.Empty) (*pb.CategoriesRes, error) {
	var categoryList []model.Category
	internal.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categoryList)
	var res pb.CategoriesRes
	b, err := json.Marshal(categoryList)
	if err != nil {
		return nil, errors.New(custom_error.MarshalCategoryFailed)
	}
	res.CategoryJsonFormat = string(b)
	return &res, nil
}

func (p ProductServer) GetSubCategory(ctx context.Context, req *pb.CategoriesReq) (*pb.SubCategoriesRes, error) {
	var category model.Category
	var res pb.SubCategoriesRes
	r := internal.DB.First(&category, req.Id)
	if r.RowsAffected < 1 {
		return nil, errors.New(custom_error.CategoryNotExists)
	}
	pre := "SubCategory"
	if category.Level == 1 {
		pre = "SubCategory.SubCategory"
	}
	var subCategoryList []model.Category
	internal.DB.Where(&model.Category{ParentCategoryID: req.Id}).Preload(pre).Find(&subCategoryList)

	b, err := json.Marshal(subCategoryList)
	if err != nil {
		return nil, errors.New(custom_error.MarshalCategoryFailed)
	}
	res.CategoryJsonFormat = string(b)
	return &res, nil
}

func ConvertCategoryModel2Pb(c model.Category) *pb.CategoryItemRes {
	item := &pb.CategoryItemRes{
		Id:    c.ID,
		Name:  c.Name,
		Level: c.Level,
	}
	if c.Level > 1 {
		item.ParentCategoryId = c.ParentCategoryID
	}
	return item
}
