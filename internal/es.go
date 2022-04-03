package internal

import (
	"context"
	"fmt"
	"github.com/i-coder-robot/mic-trainning-lessons-part2/model"
	"github.com/olivere/elastic/v7"
)

type ESConfig struct {
	Host string `mapstructure:host`
	Port int    `mapstructure:port`
}

func InitES() {
	host := fmt.Sprintf("http://%s:%d", AppConf.ESConfig.Host, AppConf.ESConfig.Port)
	var err error
	model.ESClient, err = elastic.NewClient(elastic.SetURL(host), elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}

	ok, err := model.ESClient.IndexExists(model.GetIndex()).Do(context.Background())
	if err != nil {
		panic(err)
	}
	if !ok {
		_, err = model.ESClient.CreateIndex(model.GetIndex()).BodyString(model.GetMapping()).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}
}
