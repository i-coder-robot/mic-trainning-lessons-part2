package model

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type BaseModel struct {
	ID        int32 `gorm:"primary_key;type:int"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Category struct {
	BaseModel
	Name             string `grom:"type:varchar(32);not null"`
	ParentCategoryID int32
	ParentCategory   *Category
	SubCategory      []*Category `gorm:"foreignKey:ParentCategoryID;references:ID"`
	Level            int32       `gorm:"type:int;default:1"`
}

type Brand struct {
	BaseModel
	Name string `grom:"type:varchar(32);not null"`
	Logo string `gorm:"type:varchar(256);not null;default:''"`
}

type Advertise struct {
	BaseModel
	Index int32  `gorm:"type:int;not null;default:1"`
	Image string `gorm:"type:varchar(256);not null"`
	Url   string `gorm:"type:varchar(256);not null"`
	Sort  int32  `gorm:"type:int;not null;default:1"`
}

type Product struct {
	BaseModel
	CategoryID int32 `gorm:"type:int;not null"`
	Category   Category

	BrandID int32 `gorm:"type:int;not null"`
	Brand   Brand

	Selling    bool `gorm:"default:false;not null"`
	IsShipFree bool `gorm:"default:false;not null"`
	IsPop      bool `gorm:"default:false;not null"`
	IsNew      bool `gorm:"default:false;not null"`

	Name       string  `gorm:"type:varchar(64);not null"`
	SN         string  `gorm:"type:varchar(64);not null"`
	FavNum     int32   `gorm:"type:int;not null;default:0"`
	SoldNum    int32   `gorm:"type:int;not null;default:0"`
	Price      float32 `gorm:"not null"`
	RealPrice  float32 `gorm:"not null"`
	ShortDesc  string  `gorm:"type:varchar(256);not null"`
	Images     MyList  `gorm:"type:varchar(1024);not null"`
	DescImages MyList  `gorm:"type:varchar(1024);not null"`
	CoverImage string  `gorm:"type:varchar(256);not null"`
}

type ProductCategoryBrand struct {
	BaseModel
	CategoryID int32 `gorm:"type:int;index:idx_category_brand;unique"` //联合唯一索引.2个索引是一样的，那么就是联合唯一索引了
	Category   Category
	BrandID    int32 `gorm:"type:int;index:idx_category_brand;unique"`
	Brand      Brand
}

type MyList []string

func (myList MyList) Value() (driver.Value, error) {
	return json.Marshal(myList)
}
func (myList MyList) Scan(v interface{}) error {
	return json.Unmarshal(v.([]byte), myList)
}

func (p *Product) AfterCreate(tx *gorm.DB) (err error) {
	esProduct := ESProduct{
		ID:         p.ID,
		CategoryID: p.CategoryID,
		BrandID:    p.BrandID,
		Selling:    p.Selling,
		ShipFree:   p.IsShipFree,
		IsNew:      p.IsNew,
		IsPop:      p.IsPop,
		Name:       p.Name,
		SoldNum:    p.SoldNum,
		FavNum:     p.FavNum,
		Price:      p.Price,
		RealPrice:  p.RealPrice,
		ShortDesc:  p.ShortDesc,
	}

	_, err = ESClient.Index().Index(GetIndex()).BodyJson(esProduct).Id(strconv.Itoa(int(p.ID))).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}
