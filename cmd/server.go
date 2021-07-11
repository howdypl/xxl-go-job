package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	config2 "github.com/howdypl/xxl-go-job/infrastructure/config"
	"github.com/howdypl/xxl-go-job/infrastructure/logger"
	"github.com/howdypl/xxl-go-job/interfaces/router"
	"github.com/howdypl/xxl-go-job/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	config   string
	port     string
	mode     string
	rootCmd = &cobra.Command{
		Use:               "xxl-go-job",
		Short:             "-v",
		SilenceUsage:      true,
		DisableAutoGenTag: true,
		Long:              `xxl-go-job`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires at least one arg")
			}
			return nil
		},
		PreRun: func(*cobra.Command, []string) {
			configInit()
		},
		Run: func(cmd *cobra.Command, args []string) {

			run()
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&config, "config", "c", "config/settings.yml", "Start server with provided configuration file")
	rootCmd.PersistentFlags().StringVarP(&port, "port", "p", "8002", "Tcp port server listening on")
	rootCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "dev", "server mode ; eg:dev,test,prod")
}

func configInit() {
	// 读取配置文件，并做应用的初始化
	err := config2.Setup(config)
	if err != nil {
		panic(err)
	}
}

func run() {
	// 启动服务端
	// 启动服务端
	if mode != "" {
		config2.SetConfig(config, "settings.application.mode", mode)
	}
	if viper.GetString("settings.application.mode") == string(utils.ModeProd) {
		gin.SetMode(gin.ReleaseMode)
	}

	r := router.InitRouter()


	if port != "" {
		config2.SetConfig(config, "settings.application.port", port)
	}

	srv := &http.Server{
		Addr:    config2.ApplicationConfig.Host + ":" + config2.ApplicationConfig.Port,
		Handler: r,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen: %s\n", err)
		}
	}()
	fmt.Printf("%s Server Run http://%s:%s/ \r\n",
		utils.GetCurrntTimeStr(),
		config2.ApplicationConfig.Host,
		config2.ApplicationConfig.Port)
	fmt.Printf("%s Swagger URL http://%s:%s/swagger/index.html \r\n",
		utils.GetCurrntTimeStr(),
		config2.ApplicationConfig.Host,
		config2.ApplicationConfig.Port)
	fmt.Printf("%s Enter Control + C Shutdown Server \r\n", utils.GetCurrntTimeStr())
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Printf("%s Shutdown Server ... \r\n", utils.GetCurrntTimeStr())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server Shutdown:", err)
	}
	logger.Info("Server exiting")
}

//Execute : apply commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
