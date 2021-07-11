package config

import (
	"errors"
	"fmt"
	"github.com/howdypl/xxl-go-job/infrastructure/logger"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"strings"
)


var cfgApplication *viper.Viper

// Setup 载入配置文件
func Setup(path string) error{

	viper.SetConfigFile(path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(fmt.Sprintf("Read config file fail: %s", err.Error()))
		return err
	}

	//Replace environment variables
	err = viper.ReadConfig(strings.NewReader(os.ExpandEnv(string(content))))
	if err != nil {
		log.Fatal(fmt.Sprintf("Parse config file fail: %s", err.Error()))
		return err
	}

	// 日志配置
	logger.Init()

	// 启动参数
	cfgApplication = viper.Sub("settings.application")
	if cfgApplication == nil {
		logger.Fatal("config not found settings.application")
		return errors.New("config not found settings.application")
	}
	ApplicationConfig = InitApplication(cfgApplication)


	return nil
}

func SetConfig(configPath string, key string, value interface{}) {
	viper.AddConfigPath(configPath)
	viper.Set(key, value)
	_ = viper.WriteConfig()
}
