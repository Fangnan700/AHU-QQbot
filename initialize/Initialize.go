package initialize

import (
	"Kira-qbot/model"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

var (
	initErr        error
	config         model.Config
	configDir      string
	configFileName string
	configFile     *os.File
	configData     []byte
	logDir         string
	logFileName    string
	logFile        *os.File
)

func init() {

	fmt.Println("[INFO] 正在检查配置文件")

	// 检查目录
	logDir = "log"
	configDir = "config"
	_, initErr = os.Stat(logDir)
	if os.IsNotExist(initErr) {
		initErr = os.Mkdir(logDir, 0777)
		if initErr != nil {
			fmt.Println("[ERROR] 创建日志文件夹失败")
		}
	}
	_, initErr = os.Stat(configDir)
	if os.IsNotExist(initErr) {
		initErr = os.Mkdir(configDir, 0777)
		if initErr != nil {
			fmt.Println("[ERROR] 创建配置文件夹失败")
		}
	}

	// 检查文件
	logFileName = "log/app.log"
	configFileName = "config/Config.yml"
	_, initErr = os.Stat(logFileName)
	if os.IsNotExist(initErr) {
		logFile, initErr = os.Create(logFileName)
		if initErr != nil {
			fmt.Println("[ERROR] 创建日志文件失败")
			os.Exit(-1)
		}
	}
	_, initErr = os.Stat(configFileName)
	if os.IsNotExist(initErr) {
		configFile, initErr = os.Create(configFileName)
		if initErr != nil {
			fmt.Println("[ERROR] 创建配置文件失败")
			os.Exit(-1)
		}
		config = model.Config{
			GptHost:    "",
			GptProxy:   "",
			GptModel:   "",
			GptKeys:    nil,
			RedisHost:  "",
			RedisPass:  "",
			CqHttpHost: "",
			CqHttpPath: "",
		}
		configData, initErr = yaml.Marshal(&config)
		_, initErr = configFile.Write(configData)
		if initErr != nil {
			fmt.Println("[ERROR] 写入配置文件失败")
			os.Exit(-1)
		}

		fmt.Println("[INFO] 配置文件创建完毕，请重新运行程序")
		os.Exit(0)
	}

	fmt.Println("[INFO] 初始化完毕")
}

func Init() {

}
