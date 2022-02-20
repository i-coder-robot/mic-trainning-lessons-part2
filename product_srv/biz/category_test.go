package biz

import (
	"context"
	"fmt"
	"github.com/i-coder-robot/mic-trainning-lessons-part2/proto/pb"
	"google.golang.org/protobuf/types/known/emptypb"
	"testing"
)

func TestProductServer_CreateCategory(t *testing.T) {
	//第一级
	res, err := client.CreateCategory(context.Background(), &pb.CategoryItemReq{
		Name:             "鲜肉",
		ParentCategoryId: 10,
		Level:            1,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)

	//第二级
	res2, err := client.CreateCategory(context.Background(), &pb.CategoryItemReq{
		Name:             "牛肉",
		ParentCategoryId: res.Id,
		Level:            2,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res2)

	////第三级
	res3, err := client.CreateCategory(context.Background(), &pb.CategoryItemReq{
		Name:             "牛排",
		ParentCategoryId: res2.Id,
		Level:            3,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res3)
}

func TestProductServer_DeleteCategory(t *testing.T) {
	res, err := client.CreateCategory(context.Background(), &pb.CategoryItemReq{
		Name:             "测试",
		ParentCategoryId: 10,
		Level:            1,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
	client.DeleteCategory(context.Background(), &pb.CategoryDelReq{Id: res.Id})
}

func TestProductServer_UpdateCategory(t *testing.T) {
	res, err := client.UpdateCategory(context.Background(), &pb.CategoryItemReq{
		Id:   19,
		Name: "金枪鱼",
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestProductServer_GetAllCategoryList(t *testing.T) {
	res, err := client.GetAllCategoryList(context.Background(), &emptypb.Empty{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestProductServer_GetSubCategory(t *testing.T) {
	res, err := client.GetSubCategory(context.Background(), &pb.CategoriesReq{
		Id:    18,
		Level: 2,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}
