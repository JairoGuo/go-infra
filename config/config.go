package config

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/jairoguo/go-infra/util/assert"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

const (
	env             = "CONFIG"
	defaultFilename = "config"
	defaultFiletype = YAML
)

type param struct {
	filename       string
	fileDir        string
	fileType       FileType
	enableMultiEnv bool
}

type Option func(*param)

func WithFilename(filename string) Option {
	return func(param *param) {
		param.filename = filename
	}
}

func WithFileType(filetype FileType) Option {
	return func(param *param) {
		param.fileType = filetype
	}
}

func WithMultiEnv(multiEnv bool) Option {
	return func(param *param) {
		param.enableMultiEnv = multiEnv
	}
}

// 优先级: 命令行 > 环境变量 > 默认值

func BindConfiguration(rawVal any, opts ...Option) *viper.Viper {

	param := &param{
		fileType: UNSET,
	}
	for _, opt := range opts {
		opt(param)
	}

	return handle(rawVal, *param)
}

func handle(rawVal any, param param) *viper.Viper {
	var config string

	if param.filename == "" {
		var configByParamC string
		var configByParamConfig string

		flag.StringVar(&configByParamC, "c", "", "choose config file.")
		flag.StringVar(&configByParamConfig, "config", "", "choose config file.")
		flag.Parse()

		if configByParamC != "" {
			config = configByParamC
		} else if configByParamConfig != "" {
			config = configByParamConfig
		} else {
			config = ""
		}

		if config == "" { // 判断命令行参数是否为空
			if configEnv := os.Getenv(env); configEnv == "" { // 判断 internal.Env 常量存储的环境变量是否为空
				config = defaultFilename
				param.fileType = defaultFiletype
			} else {
				config = isExtensionFile(config, &param)
				param.enableMultiEnv = false
			}
		} else {
			config = isExtensionFile(config, &param)
			param.enableMultiEnv = false
		}

	} else {
		config = isExtensionFile(param.filename, &param)
	}

	if param.fileType == UNSET {
		param.fileType = defaultFiletype
	}

	return toViper(rawVal, config, param)

}

func toViper(rawVal any, config string, param param) *viper.Viper {
	v := viper.New()

	if param.enableMultiEnv {
		multiEnvViper := viper.New()
		multiEnvViper.SetConfigFile(getFile(config, param))
		multiEnvViper.SetConfigType(param.fileType.String())
		err := multiEnvViper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
		envValue := multiEnvViper.Get("env")
		if envValue != nil {
			mode := ParseMode(envValue.(string))
			config = switchMode(config, mode, param.fileType)
		}

		multiEnvViper.OnConfigChange(func(e fsnotify.Event) {

			envValue := multiEnvViper.Get("env")
			if envValue != "" {
				mode := ParseMode(envValue.(string))
				config = switchMode(config, mode, param.fileType)
			}
			v.SetConfigFile(config)
			readConfig(v, rawVal)
		})

	}

	v.SetConfigFile(getFile(config, param))
	v.SetConfigType(param.fileType.String())
	readConfig(v, rawVal)

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		if err := v.Unmarshal(&rawVal); err != nil {
			fmt.Println(err)
		}
	})

	return v
}

func readConfig(v *viper.Viper, rawVal any) {
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	if err = v.Unmarshal(rawVal); err != nil {
		fmt.Println(err)
	}
}

func switchMode(config string, mode Mode, fileType FileType) string {

	if mode == DEFAULT {
		return config
	}
	return fmt.Sprintf("%v.%v", config, mode.String())

}

func isExtensionFile(config string, param *param) string {
	is, _ := assert.Is(param.filename, assert.FILENAME)
	if is {
		config = handleExtensionFile(param.filename, param)
	} else {
		if param.fileType == UNSET {
			param.fileType = defaultFiletype
		}
	}

	return config
}

func handleExtensionFile(filename string, param *param) string {
	fileName := filepath.Base(filename)
	dir := filepath.Dir(filename)
	fileExt := filepath.Ext(fileName)
	config := fileName[:len(fileName)-len(filepath.Ext(fileName))]

	param.fileType = ParseFileType(fileExt[1:])
	param.fileDir = dir

	return config

}

func getFile(config string, param param) string {

	return filepath.Join(param.fileDir, fmt.Sprintf("%v.%v", config, param.fileType.String()))

}
