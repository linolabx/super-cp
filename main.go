package main

import (
	"flag"
	"fmt"
	"os"

	"git.sxxfuture.net/taojiayi/super-cp/config"
	"git.sxxfuture.net/taojiayi/super-cp/core"
	"git.sxxfuture.net/taojiayi/super-cp/rules"
	"git.sxxfuture.net/taojiayi/super-cp/source"

	_ "git.sxxfuture.net/taojiayi/super-cp/targets/s3"
)

func main() {
	// 1. 解析命令行参数
	configPath := flag.String("config", ".super-cp.yml", "config file path")
	envName := flag.String("e", "default", "environment name")
	flag.Parse()

	// 2. 加载配置
	config, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Printf("load config file failed: %v\n", err)
		os.Exit(1)
	}

	ok, env := config.GetEnv(*envName)
	if !ok {
		fmt.Printf("environment %s not found\n", *envName)
		os.Exit(1)
	}

	// 3. 初始化 S3 客户端
	if err := core.Targets["@s3"].Init(env.Dist.DSN); err != nil {
		fmt.Printf("init s3 client failed: %v\n", err)
		os.Exit(1)
	}

	// 4. 扫描源文件 source -> target
	spFiles, err := source.Scan(env)
	if err != nil {
		fmt.Printf("scan files failed: %v\n", err)
		os.Exit(1)
	}

	// 5. 处理规则
	spFilesDeal, err := rules.ProcessRule(env.Rules, spFiles)
	if err != nil {
		fmt.Printf("process rules failed: %v\n", err)
		os.Exit(1)
	}

	// 6. 上传文件
	for _, spFile := range spFilesDeal {
		if err := core.Targets["@s3"].Upload(spFile); err != nil {
			fmt.Printf("upload files failed: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("upload files success")
}
