package custom_error

const (
	BrandAlreadyExits = "品牌已存在"
	DelBrandFail      = "品牌删除失败"
	BrandNotExits     = "品牌不存在"

	ADNotExists = "广告不存在"

	CategoryNotExists     = "分类不存在"
	MarshalCategoryFailed = "序列化分类错误"

	DelProductCategoryBrandFailed = "删除分类品牌表失败"
	ProductCategoryBrandNotFound  = "品牌分类表找不到记录"

	DelProductFailed = "删除产品失败"
	ProductNotExists = "产品不存在"

	ParamError = "参数错误"
)
