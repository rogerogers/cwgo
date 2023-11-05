/*
 *
 *  * Copyright 2022 CloudWeGo Authors
 *  *
 *  * Licensed under the Apache License, Version 2.0 (the "License");
 *  * you may not use this file except in compliance with the License.
 *  * You may obtain a copy of the License at
 *  *
 *  *     http://www.apache.org/licenses/LICENSE-2.0
 *  *
 *  * Unless required by applicable law or agreed to in writing, software
 *  * distributed under the License is distributed on an "AS IS" BASIS,
 *  * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  * See the License for the specific language governing permissions and
 *  * limitations under the License.
 *
 */

package svc

import (
	"github.com/cloudwego/cwgo/platform/server/cmd/api/pkg/manager"
	"github.com/cloudwego/cwgo/platform/server/shared/config"
	"github.com/cloudwego/cwgo/platform/server/shared/dao"
	"github.com/cloudwego/cwgo/platform/server/shared/logger"
	"github.com/cloudwego/cwgo/platform/server/shared/registry"
	"github.com/cloudwego/cwgo/platform/server/shared/repository"
	"go.uber.org/zap"
)

type ServiceContext struct {
	DaoManager      *dao.Manager
	RepoManager     *repository.Manager
	BuiltinRegistry *registry.BuiltinRegistry
	Manager         *manager.Manager
}

var Svc *ServiceContext

func InitServiceContext() {
	// init dao manager
	logger.Logger.Info("initializing dao manager")
	daoManager, err := dao.NewDaoManager(config.GetManager().Config.Store)
	if err != nil {
		logger.Logger.Fatal("initialize dao manager failed", zap.Error(err))
	}
	logger.Logger.Info("initialize dao manager successfully")

	logger.Logger.Info("initializing dao manager")
	repoManager, err := repository.NewRepoManager(daoManager)
	if err != nil {
		logger.Logger.Fatal("initialize repository manager failed")
	}
	logger.Logger.Info("initialize dao manager successfully")

	Svc = &ServiceContext{
		DaoManager:      daoManager,
		RepoManager:     repoManager,
		BuiltinRegistry: config.GetManager().ApiConfigManager.RegistryConfigManager.GetRegistry().(*registry.BuiltinRegistry),
		Manager: manager.NewManager(
			config.GetManager().Config.App,
			daoManager,
			config.GetManager().Config.Api.Dispatcher.NewDispatcher(),
			config.GetManager().ApiConfigManager.RegistryConfigManager.GetRegistry(),
			config.GetManager().ApiConfigManager.RegistryConfigManager.GetDiscoveryResolver(),
		),
	}
}
