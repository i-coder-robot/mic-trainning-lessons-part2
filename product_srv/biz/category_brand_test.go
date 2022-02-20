package biz

import (
	"context"
	"fmt"
	"github.com/i-coder-robot/mic-trainning-lessons-part2/proto/pb"
	"testing"
)

func TestProductServer_CreateCategoryBrand(t *testing.T) {
	res, err := client.CreateCategoryBrand(context.Background(), &pb.CategoryBrandReq{
		CategoryId: 22,
		BrandId:    17,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestProductServer_UpdateCategoryBrand(t *testing.T) {
	client.UpdateCategoryBrand(context.Background(), &pb.CategoryBrandReq{
		Id:         2,
		CategoryId: 22,
		BrandId:    19,
	})
}

func TestProductServer_CategoryBrandList(t *testing.T) {
	res, err := client.CategoryBrandList(context.Background(), &pb.PagingReq{
		PageNo:   1,
		PageSize: 5,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestProductServer_GetCategoryBrandList(t *testing.T) {
	res, err := client.GetCategoryBrandList(context.Background(), &pb.CategoryItemReq{
		Id:               19,
		ParentCategoryId: 22,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestProductServer_DeleteCategoryBrand(t *testing.T) {
	client.DeleteCategoryBrand(context.Background(), &pb.CategoryBrandReq{
		Id: 1,
	})
}
