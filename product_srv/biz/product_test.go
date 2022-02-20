package biz

import (
	"context"
	"fmt"
	"github.com/i-coder-robot/mic-trainning-lessons-part2/proto/pb"
	"testing"
)

func TestProductServer_CreateProduct(t *testing.T) {
	for i := 0; i < 8; i++ {
		res, err := client.CreateProduct(context.Background(), &pb.CreateProductItem{
			Name:        fmt.Sprintf("黄金牛排%d", i),
			Sn:          "123456789",
			CategoryId:  22,
			Price:       359.00,
			RealPrice:   199.00,
			ShortDesc:   "",
			ProductDesc: "",
			Images:      nil,
			DescImages:  nil,
			CoverImage:  "https://space.bilibili.com/375038855",
			IsNew:       true,
			IsPop:       true,
			Selling:     true,
			BrandId:     18,
			FavNum:      6666,
			SoldNum:     5432,
			IsShipFree:  false,
		})
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(res)
	}
}

func TestProductServer_UpdateProduct(t *testing.T) {
	client.UpdateProduct(context.Background(), &pb.CreateProductItem{
		Id:         9,
		CategoryId: 22,
		BrandId:    18,
		Name:       "战斧牛排66666",
	})
}

func TestProductServer_DeleteProduct(t *testing.T) {
	client.DeleteProduct(context.Background(),
		&pb.ProductDelItem{
			Id: 9,
		})
}

func TestProductServer_BatchGetProduct(t *testing.T) {
	ids := []int32{10, 11, 12}
	res, err := client.BatchGetProduct(context.Background(),
		&pb.BatchProductIdReq{
			Ids: ids,
		})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestProductServer_ProductList(t *testing.T) {
	res, err := client.ProductList(context.Background(), &pb.ProdcutConditionReq{
		PageNo:   2,
		PageSize: 5,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestProductServer_GetProductDetail(t *testing.T) {
	res, err := client.GetProductDetail(context.Background(), &pb.ProductItemReq{
		Id: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}
