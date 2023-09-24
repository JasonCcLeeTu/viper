package main

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type config struct {
	Redis redis `mapstructure:"redis"`
}

type redis struct {
	ConnectString string `mapstructure:"connectionString"`
	Password      string `mapstructure:"password"`
	DB            string `mapstructure:"db"`
}

var forever = make(chan struct{})

func main() {

	viper.SetConfigName("local_old")             //設置讀取得config檔案名稱
	viper.AddConfigPath("../conf/")              //設置讀取得config檔案路徑
	viper.SetConfigType("json")                  //設置讀取得config檔案類型
	if err := viper.ReadInConfig(); err != nil { //讀取設置的config檔案
		log.Fatalf("viper ReadInConfig error:%s", err.Error())
	}

	viper.SetConfigName("local_new")              //設置讀取得config檔案名稱
	if err := viper.MergeInConfig(); err != nil { //將前面設置的兩個conmfig,併在一起,如果檔案內的field名稱相同則對應的value新的會覆蓋舊的(新的是後面設置檔案也就是local_new)
		log.Fatalf("viper MergeInConfig error:%s", err.Error())
	}

	config := config{}
	if err := viper.Unmarshal(&config); err != nil { //將設計的config的struct帶入viper unmarshal 賦值
		log.Fatalf("viper Unmarshal config error:%s ", err.Error())
	}

	log.Printf("config:\n%+v", config)

	viper.OnConfigChange( //添加function,讓想要做的動作寫到function中,當執行viper.WatchConfig()時,只要config檔案有異動就會觸發該function
		func(in fsnotify.Event) { //當config異動,就再 unmarshal config  更新一下
			viper.Unmarshal(&config)
			log.Printf("異動config:\n%+v", config)
		},
	)
	viper.WatchConfig() // 開啟監視config異動的功能,會有一個線程監視著

	<-forever
}
