package internal

type ProductSrvConfig struct {
	SrvName string   `mapstruct:"srvName" json:"srvName"`
	Host    string   `mapstructure:"host" json:"host"`
	Port    int      `mapstructure:"port" json:"port"`
	Tags    []string `mapstructure:"tags" json:"tags"`
}

type ProductWebConfig struct {
	SrvName string   `mapstruct:"srvName" json:"srvName"`
	Host    string   `mapstructure:"host" json:"host"`
	Port    int      `mapstructure:"port" json:"port"`
	Tags    []string `mapstructure:"tags" json:"tags"`
}

type AppConfig struct {
	DBConfig         DBConfig         `mapstructure:"db" json:"db"`
	RedisConfig      RedisConfig      `mapstructure:"redis" json:"redis"`
	ConsulConfig     ConsulConfig     `mapstructure:"consul" json:"consul"`
	AccountSrvConfig ProductSrvConfig `mapstructure:"account_srv" json:"account_srv"`
	AccountWebConfig ProductWebConfig `mapstructure:"account_web" json:"account_web"`
	JWTConfig        JWTConfig        `mapstructure:"jwt" json:"jwt"`
}
