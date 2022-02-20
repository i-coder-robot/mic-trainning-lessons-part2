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

func (p ProductServer) AdvertiseList(ctx context.Context, empty *emptypb.Empty) (*pb.AdvertisesRes, error) {
	var adList []model.Advertise
	var adItemList []*pb.AdvertiseItemRes
	var advertiseRes pb.AdvertisesRes
	r := internal.DB.Find(&adList)
	for _, item := range adList {
		adItemList = append(adItemList, ConvertAdModel2Pb(item))
	}
	advertiseRes.Total = int32(r.RowsAffected)
	advertiseRes.ItemList = adItemList
	return &advertiseRes, nil

}

func (p ProductServer) CreateAdvertise(ctx context.Context, req *pb.AdvertiseReq) (*pb.AdvertiseItemRes, error) {
	var ad model.Advertise
	ad.Index = req.Index
	ad.Image = req.Image
	ad.Url = req.Url
	internal.DB.Save(&ad)
	return ConvertAdModel2Pb(ad), nil
}

func (p ProductServer) DeleteAdvertise(ctx context.Context, req *pb.AdvertiseReq) (*emptypb.Empty, error) {
	internal.DB.Delete(&model.Advertise{}, req.Id)
	//TODO 需要业务判断失败
	return &emptypb.Empty{}, nil
}

func (p ProductServer) UpdateAdvertise(ctx context.Context, req *pb.AdvertiseReq) (*emptypb.Empty, error) {
	var ad model.Advertise
	r := internal.DB.First(&ad, req.Id)
	if r.RowsAffected < 1 {
		return nil, errors.New(custom_error.ADNotExists)
	}
	if req.Index > 0 {
		ad.Index = req.Index
	}
	if req.Image != "" {
		ad.Image = req.Image
	}
	if req.Url != "" {
		ad.Url = req.Url
	}
	internal.DB.Save(&ad)
	return &emptypb.Empty{}, nil
}
func ConvertAdModel2Pb(item model.Advertise) *pb.AdvertiseItemRes {
	ad := &pb.AdvertiseItemRes{
		Index: item.Index,
		Image: item.Image,
		Url:   item.Url,
	}
	if item.ID > 0 {
		ad.Id = item.ID
	}
	return ad
}
