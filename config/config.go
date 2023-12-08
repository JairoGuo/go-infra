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

type Param struct {
	Filename       string
	FileType       FileType
	EnableMultiEnv bool
}

type Option func(*Param)

func WithFilename(filename string) Option {
	return func(param *Param) {
		param.Filename = filename
	}
}

func WithFileType(filetype FileType) Option {
	return func(param *Param) {
		param.FileType = filetype
	}
}

func WithMultiEnv(multiEnv bool) Option {
	return func(param *Param) {
		param.EnableMultiEnv = multiEnv
	}
}

// 优先级: 命令行 > 环境变量 > 默认值

func BindConfiguration(rawVal any, opts ...Option) {

	param := &Param{
		FileType: UNSET,
	}
	for _, opt := range opts {
		opt(param)
	}

	handle(rawVal, *param)
}

func handle(rawVal any, param Param) *viper.Viper {
	var config string

	if param.Filename == "" {
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
				param.FileType = defaultFiletype
			} else {
				config = configEnv
				param.EnableMultiEnv = false
			}
		} else {
			param.EnableMultiEnv = false
		}

	} else {
		config = isExtensionFile(param.Filename, &param)
	}

	if param.FileType == UNSET {
		param.FileType = defaultFiletype
	}

	return toViper(rawVal, config, param)

}

func toViper(rawVal any, config string, param Param) *viper.Viper {
	v := viper.New()

	if param.EnableMultiEnv {
		multiEnvViper := viper.New()
		multiEnvViper.SetConfigFile(getFile(config, param))
		multiEnvViper.SetConfigType(param.FileType.String())
		err := multiEnvViper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
		envValue := multiEnvViper.Get("env")
		if envValue != nil {
			mode := ParseMode(envValue.(string))
			config = SwitchMode(config, mode, param.FileType)
		}

		multiEnvViper.OnConfigChange(func(e fsnotify.Event) {

			envValue := multiEnvViper.Get("env")
			if envValue != "" {
				mode := ParseMode(envValue.(string))
				config = SwitchMode(config, mode, param.FileType)
			}
			v.SetConfigFile(config)
			readConfig(v, rawVal)
		})

	}

	v.SetConfigFile(getFile(config, param))
	v.SetConfigType(param.FileType.String())
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

func SwitchMode(config string, mode Mode, fileType FileType) string {

	if mode == DEFAULT {
		return config
	}
	return fmt.Sprintf("%v.%v", config, mode.String())

}

func isExtensionFile(config string, param *Param) string {
	is, _ := assert.Is(param.Filename, assert.FILENAME)
	if is {
		config = handleExtensionFile(param.Filename, param)
	} else {
		if param.FileType == UNSET {
			param.FileType = defaultFiletype
		}
	}

	return config
}

func handleExtensionFile(filename string, param *Param) string {
	fileName := filepath.Base(filename)
	fileExt := filepath.Ext(fileName)
	param.FileType = ParseFileType(fileExt[1:])

	return filename

}

func getFile(config string, param Param) string {
	return fmt.Sprintf("%v.%v", config, param.FileType.String())

}
