package biz

import (
	"context"
	"fmt"
	"github.com/i-coder-robot/mic-trainning-lessons-part2/proto/pb"
	"google.golang.org/protobuf/types/known/emptypb"
	"testing"
)

func TestProductServer_CreateAdvertise(t *testing.T) {
	for i := 0; i < 8; i++ {
		res, err := client.CreateAdvertise(context.Background(), &pb.AdvertiseReq{
			Index: int32(i),
			Image: fmt.Sprintf("image-%d", i),
			Url:   fmt.Sprintf("url-%d", i),
		})
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(res.Id)
	}
}

func TestProductServer_AdvertiseList(t *testing.T) {
	res, err := client.AdvertiseList(context.Background(), &emptypb.Empty{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestProductServer_DeleteAdvertise(t *testing.T) {
	res, err := client.DeleteAdvertise(context.Background(), &pb.AdvertiseReq{
		Id: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestProductServer_UpdateAdvertise(t *testing.T) {
	client.UpdateAdvertise(context.Background(), &pb.AdvertiseReq{
		Id:    6,
		Index: 666,
		Image: "666",
		Url:   "666",
	})
}
