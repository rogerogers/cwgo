/*
 * Copyright 2022 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config

import (
	"fmt"
	"github.com/cloudwego/cwgo/platform/server/shared/config/internal/logger"
	"github.com/cloudwego/cwgo/platform/server/shared/config/internal/registry"
	"github.com/cloudwego/cwgo/platform/server/shared/config/internal/store"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
	"github.com/spf13/viper"
)

type Manager struct {
	ServerType         consts.ServerType
	ServerMode         consts.ServerMode
	ServiceId          string
	Config             Config
	StoreConfigManager *store.StoreConfigManager
}

type Config struct {
	Logger   logger.Config   `mapstructure:"logger"`
	Registry registry.Config `mapstructure:"registry"`
	Store    store.Config    `mapstructure:"store"`
}

var manager *Manager

type FileConfig struct {
	Path string
}

func InitManager(serverType consts.ServerType, serverMode consts.ServerMode, configType consts.ConfigType, metadata ...interface{}) error {
	var config Config

	switch configType {
	case consts.ConfigTypeNumFile:
		var configPath string

		if metadata != nil {
			if fileConfig, ok := metadata[0].(FileConfig); ok {
				configPath = fileConfig.Path
			}
		}

		configPath = fmt.Sprintf("config-%s.yaml", consts.ServerModeMapToStr[serverMode])

		fmt.Printf("get Config path: %s", configPath)

		v := viper.New()
		v.SetConfigType("yaml")
		v.SetConfigFile(configPath)
		err := v.ReadInConfig()
		if err != nil {
			panic(fmt.Sprintf("get Config file failed, err: %v", err))
		}

		if err := v.Unmarshal(&config); err != nil {
			return fmt.Errorf("unmarshal Config failed, err: %v", err)
		}

	case consts.ConfigTypeNumApollo:
		// TODO: to be implemented
		panic("to be implemented")
	default:

	}

	serviceId, err := utils.NewServiceId()
	if err != nil {
		return err
	}

	manager = &Manager{
		ServerType:         serverType,
		ServerMode:         serverMode,
		ServiceId:          serviceId,
		Config:             config,
		StoreConfigManager: nil,
	}

	return nil
}

func GetManager() *Manager {
	if manager == nil {
		panic("Config manager.go not initialized")
	}

	return manager
}
