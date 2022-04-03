package biz

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/i-coder-robot/mic-trainning-lessons-part2/custom_error"
	"github.com/i-coder-robot/mic-trainning-lessons-part2/internal"
	"github.com/i-coder-robot/mic-trainning-lessons-part2/model"
	"github.com/i-coder-robot/mic-trainning-lessons-part2/proto/pb"
	"github.com/olivere/elastic/v7"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProductServer struct {
}

type EsCategory struct {
	CategoryID int32
}

func (p ProductServer) CreateProduct(ctx context.Context, req *pb.CreateProductItem) (*pb.ProductItemRes, error) {
	var category model.Category
	var brand model.Brand
	var res *pb.ProductItemRes

	//TODO 业务逻辑判断更复杂
	r := internal.DB.First(&category, req.CategoryId)
	if r.RowsAffected < 1 {
		return nil, errors.New(custom_error.CategoryNotExists)
	}
	r = internal.DB.First(&brand, req.BrandId)
	if r.RowsAffected < 1 {
		return nil, errors.New(custom_error.BrandNotExits)
	}

	item := ConvertReq2Model(model.Product{}, req, category, brand)

	tx := internal.DB.Begin()
	result := internal.DB.Save(&item)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()
	res = ConvertProductModel2Pb(item)
	return res, nil
}

func (p ProductServer) DeleteProduct(ctx context.Context, req *pb.ProductDelItem) (*emptypb.Empty, error) {
	r := internal.DB.Delete(&model.Product{}, req.Id)
	if r.RowsAffected < 1 {
		return nil, errors.New(custom_error.DelProductFailed)
	}
	return &emptypb.Empty{}, nil
}

func (p ProductServer) UpdateProduct(ctx context.Context, req *pb.CreateProductItem) (*emptypb.Empty, error) {
	//TODO 业务逻辑判断
	var pro model.Product
	var c model.Category
	var b model.Brand

	r := internal.DB.First(&pro, req.Id)
	if r.RowsAffected < 1 {
		return nil, errors.New(custom_error.ProductNotExists)
	}
	r = internal.DB.First(&c, req.CategoryId)
	if r.RowsAffected < 1 {
		return nil, errors.New(custom_error.CategoryNotExists)
	}
	r = internal.DB.First(&b, req.BrandId)
	if r.RowsAffected < 1 {
		return nil, errors.New(custom_error.BrandNotExits)
	}
	pro = ConvertReq2Model(pro, req, c, b)
	internal.DB.Save(&pro)
	return &emptypb.Empty{}, nil
}

func (p ProductServer) ProductList(ctx context.Context, req *pb.ProductConditionReq) (*pb.ProductsRes, error) {
	//iDb := internal.DB.Model(model.Product{})
	//var procudtList []model.Product
	//var itemList []*pb.ProductItemRes
	//var res pb.ProductsRes
	//
	//if req.IsPop {
	//	iDb = iDb.Where("is_pop=?", req.IsPop)
	//}
	//
	//if req.IsNew {
	//	iDb = iDb.Where("is_new=?", req.IsNew)
	//}
	//
	//if req.BrandId > 0 {
	//	iDb = iDb.Where("brand_id=?", req.BrandId)
	//}
	//if req.KeyWord != "" {
	//	iDb = iDb.Where("key_word like ?", "%"+req.KeyWord+"%")
	//}
	//if req.MinPrice > 0 {
	//	iDb = iDb.Where("min_price>?", req.MinPrice)
	//}
	//if req.MaxPrice > 0 {
	//	iDb = iDb.Where("max_price=?", req.MaxPrice)
	//}
	//
	//if req.CategoryId > 0 {
	//	var category model.Category
	//	r := internal.DB.First(&category, req.CategoryId)
	//	if r.RowsAffected == 0 {
	//		return nil, errors.New(custom_error.CategoryNotExists)
	//	}
	//	var q string
	//	if category.Level == 1 {
	//		q = fmt.Sprintf("select id from category where parent_category_id in (select id from category WHERE parent_category_id=%d)", req.CategoryId)
	//	} else if category.Level == 2 {
	//		q = fmt.Sprintf("select id from category WHERE parent_category_id=%d", req.CategoryId)
	//	} else if category.Level == 3 {
	//		q = fmt.Sprintf("select id from category WHERE id=%d", req.CategoryId)
	//	}
	//	iDb = iDb.Where(fmt.Sprintf("category_id in %s", q))
	//}
	//var count int64
	//iDb.Count(&count)
	//fmt.Println(count)
	//
	//iDb.Joins("Category").Joins("Brand").Scopes(internal.MyPaging(int(req.PageNo), int(req.PageSize))).Find(&procudtList)
	//for _, item := range procudtList {
	//	res := ConvertProductModel2Pb(item)
	//	itemList = append(itemList, res)
	//}
	//res.ItemList = itemList
	//res.Total = int32(count)
	//return &res, nil

	//ES版本
	var res pb.ProductsRes
	q := elastic.NewBoolQuery()
	localDB := internal.DB.Model(model.Product{})
	if req.KeyWord != "" {
		q = q.Must(elastic.NewMultiMatchQuery(req.KeyWord, "name", "short_desc"))
	}
	if req.IsPop {
		q = q.Filter(elastic.NewTermQuery("is_pop", req.IsPop))
	}
	if req.IsNew {
		q = q.Filter(elastic.NewTermQuery("is_new", req.IsNew))
	}

	if req.MinPrice > 0 {
		q = q.Filter(elastic.NewRangeQuery("real_price").Gte(req.MinPrice))
	}
	if req.MaxPrice > 0 {
		q = q.Filter(elastic.NewRangeQuery("real_price").Lte(req.MaxPrice))
	}

	if req.BrandId > 0 {
		q = q.Filter(elastic.NewTermQuery("brand_id", req.BrandId))
	}

	var subQuery string
	categoryIdList := make([]interface{}, 0)
	if req.CategoryId > 0 {
		var category model.Category
		if result := internal.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "商品分类不存在")
		}

		if category.Level == 1 {
			subQuery = fmt.Sprintf("select id from category where parent_category_id in (select id from category WHERE parent_category_id=%d)", req.CategoryId)
		} else if category.Level == 2 {
			subQuery = fmt.Sprintf("select id from category WHERE parent_category_id=%d", req.CategoryId)
		} else if category.Level == 3 {
			subQuery = fmt.Sprintf("select id from category WHERE id=%d", req.CategoryId)
		}

		var EsCategoryList []EsCategory
		internal.DB.Model(model.Category{}).Raw(subQuery).Scan(&EsCategoryList)
		for _, item := range EsCategoryList {
			categoryIdList = append(categoryIdList, item.CategoryID)
		}
		q = q.Filter(elastic.NewTermsQuery("category_id", categoryIdList...))
	}

	if req.PageNo < 1 {
		req.PageNo = 1
	}

	switch {
	case req.PageSize > 100:
		req.PageSize = 100
	case req.PageSize < 1:
		req.PageSize = 10
	}
	result, err := internal.ESClient.Search().Index(model.GetIndex()).Query(q).
		From(int(req.PageNo)).Size(int(req.PageSize)).Do(context.Background())
	if err != nil {
		return nil, err
	}

	productIdList := make([]int32, 0)
	res.Total = int32(result.Hits.TotalHits.Value)
	for _, value := range result.Hits.Hits {
		esProduct := model.ESProduct{}
		_ = json.Unmarshal(value.Source, &esProduct)
		productIdList = append(productIdList, esProduct.ID)
	}

	var products []model.Product
	re := localDB.Preload("Category").Preload("Brand").Find(&products, productIdList)
	if re.Error != nil {
		return nil, re.Error
	}

	for _, item := range products {
		itemRes := ModelToResponse(item)
		res.ItemList = append(res.ItemList, itemRes)
	}

	return &res, nil
}

func ModelToResponse(product model.Product) *pb.ProductItemRes {
	return &pb.ProductItemRes{
		Id:         product.ID,
		CategoryId: product.CategoryID,
		Name:       product.Name,
		Sn:         product.SN,
		SoldNum:    product.SoldNum,
		FavNum:     product.FavNum,
		Price:      product.Price,
		RealPrice:  product.RealPrice,
		ShortDesc:  product.ShortDesc,
		CoverImage: product.CoverImage,
		IsNew:      product.IsNew,
		IsPop:      product.IsPop,
		Selling:    product.Selling,
		DescImages: product.DescImages,
		Images:     product.Images,
		Category: &pb.CategoryShortItemRes{
			Id:   product.Category.ID,
			Name: product.Category.Name,
		},
		Brand: &pb.BrandItemRes{
			Id:   product.Brand.ID,
			Name: product.Brand.Name,
			Logo: product.Brand.Logo,
		},
	}
}

func (p ProductServer) BatchGetProduct(ctx context.Context, req *pb.BatchProductIdReq) (*pb.ProductsRes, error) {
	var productList []model.Product
	var res pb.ProductsRes
	r := internal.DB.Find(&productList, req.Ids)
	res.Total = int32(r.RowsAffected)
	for _, item := range productList {
		pro := ConvertProductModel2Pb(item)
		res.ItemList = append(res.ItemList, pro)
	}
	return &res, nil
}

func (p ProductServer) GetProductDetail(ctx context.Context, req *pb.ProductItemReq) (*pb.ProductItemRes, error) {
	var pro model.Product
	var res *pb.ProductItemRes
	r := internal.DB.First(&pro, req.Id)
	if r.RowsAffected < 1 {
		return nil, errors.New(custom_error.ProductNotExists)
	}
	res = ConvertProductModel2Pb(pro)
	return res, nil
}

func ConvertReq2Model(p model.Product, req *pb.CreateProductItem, category model.Category, brand model.Brand) model.Product {

	if req.CategoryId > 0 {
		p.CategoryID = req.CategoryId
		p.Category = category
	}
	if req.BrandId > 0 {
		p.BrandID = req.BrandId
		p.Brand = brand
	}
	if req.Selling {
		p.Selling = true
	} else {
		p.Selling = false
	}
	if req.Selling {
		p.Selling = true
	} else {
		p.Selling = false
	}
	if req.IsShipFree {
		p.IsShipFree = true
	} else {
		p.IsShipFree = false
	}
	if req.IsPop {
		p.IsPop = true
	} else {
		p.IsPop = false
	}

	if req.IsNew {
		p.IsNew = true
	} else {
		p.IsNew = false
	}
	if req.Name != "" {
		p.Name = req.Name
	}
	if req.Sn != "" {
		p.SN = req.Sn
	}
	if req.FavNum > 0 {
		p.FavNum = req.FavNum
	}
	if req.SoldNum > 0 {
		p.SoldNum = req.SoldNum
	}
	if req.Price > 0 {
		p.Price = req.Price
	}
	if req.RealPrice > 0 {
		p.RealPrice = req.RealPrice
	}
	if req.ShortDesc != "" {
		p.ShortDesc = req.ShortDesc
	}
	if req.Images != nil {
		p.Images = req.Images
	}
	if req.DescImages != nil {
		p.DescImages = req.DescImages
	}
	if req.CoverImage != "" {
		p.CoverImage = req.CoverImage
	}
	if req.Id > 0 {
		p.ID = req.Id
	}
	return p
}

func ConvertProductModel2Pb(pro model.Product) *pb.ProductItemRes {
	p := &pb.ProductItemRes{
		Id:         pro.ID,
		CategoryId: pro.CategoryID,
		Name:       pro.Name,
		Sn:         pro.SN,
		SoldNum:    pro.SoldNum,
		FavNum:     pro.FavNum,
		Price:      pro.Price,
		RealPrice:  pro.RealPrice,
		ShortDesc:  pro.ShortDesc,
		Images:     pro.Images,
		DescImages: pro.DescImages,
		CoverImage: pro.CoverImage,
		IsNew:      pro.IsNew,
		IsPop:      pro.IsPop,
		Selling:    pro.Selling,
		Category: &pb.CategoryShortItemRes{
			Id:   pro.Category.ID,
			Name: pro.Category.Name,
		},
		Brand: &pb.BrandItemRes{
			Id:   pro.Brand.ID,
			Name: pro.Brand.Name,
			Logo: pro.Brand.Logo,
		},
	}
	return p
}
