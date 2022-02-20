package biz

import (
	"context"
	"fmt"
	"github.com/i-coder-robot/mic-trainning-lessons-part2/internal"
	"github.com/i-coder-robot/mic-trainning-lessons-part2/proto/pb"
	"google.golang.org/grpc"
	"testing"
)

var client pb.ProductServiceClient

func init() {
	addr := fmt.Sprintf("%s:%d", internal.AppConf.ProductSrvConfig.Host, internal.AppConf.ProductSrvConfig.Port)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client = pb.NewProductServiceClient(conn)
}

func TestProductServer_CreateBrand(t *testing.T) {
	brands := []string{
		"大希地", "恒都", "小牛凯西", "天莱香牛", "伊赛", "春禾秋牧", "元盛", "大庄园",
	}
	for _, item := range brands {
		res, err := client.CreateBrand(context.Background(), &pb.BrandItemReq{
			Name: item,
			Logo: "https://space.bilibili.com/375038855",
		})
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res.Id)
	}

}

func TestProductServer_BrandList(t *testing.T) {
	res, err := client.BrandList(context.Background(), &pb.BrandPagingReq{
		PageNo:   2,
		PageSize: 5,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestProductServer_DeleteBrand(t *testing.T) {
	res, err := client.DeleteBrand(context.Background(), &pb.BrandItemReq{
		Id: 24,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestProductServer_UpdateBrand(t *testing.T) {
	res, err := client.UpdateBrand(context.Background(), &pb.BrandItemReq{
		Id:   23,
		Name: "update测试一下",
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}
