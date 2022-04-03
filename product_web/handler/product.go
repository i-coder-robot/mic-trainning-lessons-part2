package handler

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/i-coder-robot/mic-trainning-lessons-part2/custom_error"
	"github.com/i-coder-robot/mic-trainning-lessons-part2/internal"
	"github.com/i-coder-robot/mic-trainning-lessons-part2/product_web/req"
	"github.com/i-coder-robot/mic-trainning-lessons-part2/proto/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net/http"
	"strconv"
)

var productClient pb.ProductServiceClient

func init() {
	//addr := fmt.Sprintf("%s:%d", internal.AppConf.ProductSrvConfig.Host, internal.AppConf.ProductSrvConfig.Port)
	addr := fmt.Sprintf("%s:%d", internal.AppConf.ConsulConfig.Host, internal.AppConf.ConsulConfig.Port)
	dialAddr := fmt.Sprintf("consul://%s/product_srv?wait=14s&tag=happy", addr)
	conn, err := grpc.Dial(
		dialAddr,
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal(err)
	}
	productClient = pb.NewProductServiceClient(conn)
}

func ProductListHandler(c *gin.Context) {
	var condition pb.ProductConditionReq
	//c.ShouldBindJSON(&condition)
	//list?pageNo=1&pageSize=2

	minPriceStr := c.DefaultQuery("minPrice", "0")
	minPrice, err := strconv.Atoi(minPriceStr)
	if err != nil {
		zap.S().Error("minPrice error")
		c.JSON(http.StatusOK, gin.H{"msg": custom_error.ParamError})
		return
	}
	maxPriceStr := c.DefaultQuery("maxPrice", "0")
	maxPrice, err := strconv.Atoi(maxPriceStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"msg": custom_error.ParamError})
		return
	}
	condition.MinPrice = int32(minPrice)
	condition.MaxPrice = int32(maxPrice)

	categoryIdStr := c.DefaultQuery("categoryId", "0")
	categoryId, err := strconv.Atoi(categoryIdStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"msg": custom_error.ParamError})
		return
	}

	condition.CategoryId = int32(categoryId)

	brandIdStr := c.DefaultQuery("brandId", "0")
	brandId, err := strconv.Atoi(brandIdStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"msg": custom_error.ParamError})
		return
	}

	condition.BrandId = int32(brandId)

	isHot := c.DefaultQuery("isPop", "0")
	if isHot == "1" {
		condition.IsPop = true
	}

	isNew := c.DefaultQuery("isNew", "0")
	if isNew == "1" {
		condition.IsNew = true
	}

	pageNoStr := c.DefaultQuery("pageNo", "1")
	pageNo, err := strconv.Atoi(pageNoStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"msg": custom_error.ParamError})
		return
	}

	condition.PageNo = int32(pageNo)

	pageSizeStr := c.DefaultQuery("pageSize", "10")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"msg": custom_error.ParamError})
		return
	}
	condition.PageSize = int32(pageSize)

	keyword := c.DefaultQuery("keyword", "")
	condition.KeyWord = keyword

	r, err := productClient.ProductList(context.Background(), &condition)

	if err != nil {
		zap.S().Error(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": "产品列表查询失败",
			//默认值
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":   "",
		"total": r.Total,
		"data":  r.ItemList,
	})

}

func DetailHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		zap.S().Error(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": "参数错误",
		})
		return
	}
	res, err := productClient.GetProductDetail(context.Background(), &pb.ProductItemReq{Id: int32(id)})
	if err != nil {
		zap.S().Error(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": "获取详情失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "",
		"data": res,
	})
}

func AddHandler(c *gin.Context) {
	var productReq req.ProductReq
	err := c.ShouldBindJSON(&productReq)
	if err != nil {
		zap.S().Error(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": "参数解析错误",
		})
		return
	}
	r := ConvertProductReq2Pb(productReq)
	res, err := productClient.CreateProduct(context.Background(), r)
	if err != nil {
		zap.S().Error(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": "添加产品失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "",
		"data": res,
	})
}

func DelHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		zap.S().Error(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": "参数错误",
		})
		return
	}
	_, err = productClient.DeleteProduct(context.Background(), &pb.ProductDelItem{Id: int32(id)})
	if err != nil {
		zap.S().Error(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": "删除产品失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "",
	})
}

func UpdateHandler(c *gin.Context) {
	var productReq req.ProductReq
	err := c.ShouldBindJSON(&productReq)
	if err != nil {
		zap.S().Error(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": "参数解析错误",
		})
		return
	}
	r := ConvertProductReq2Pb(productReq)

	_, err = productClient.UpdateProduct(context.Background(), r)
	if err != nil {
		zap.S().Error(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": "更新产品失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "",
	})
}

func ConvertProductReq2Pb(productReq req.ProductReq) *pb.CreateProductItem {
	item := pb.CreateProductItem{
		Name:        productReq.Name,
		Sn:          productReq.SN,
		Price:       productReq.Price,
		RealPrice:   productReq.RealPrice,
		ShortDesc:   productReq.ShortDesc,
		ProductDesc: productReq.Desc,
		Images:      productReq.Images,
		DescImages:  productReq.DescImages,
		CoverImage:  productReq.CoverImage,
		IsNew:       productReq.IsNew,
		IsPop:       productReq.IsPop,
		Selling:     productReq.Selling,
		CategoryId:  productReq.CategoryId,
		BrandId:     productReq.BrandId,
		FavNum:      productReq.FavNum,
		SoldNum:     productReq.SoldNum,
	}
	if productReq.Id > 0 {
		item.Id = productReq.Id
	}
	return &item
}
