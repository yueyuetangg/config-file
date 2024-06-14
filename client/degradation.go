// Copyright 2024 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	kitexclient "github.com/cloudwego/kitex/client"
	"github.com/kitex-contrib/config-file/monitor"
	degradation "github.com/kitex-contrib/config-file/pkg"
)

func WithDegradation(watcher monitor.ConfigMonitor) []kitexclient.Option {
	container, keyDegradation := initDegradationOptions(watcher)
	return []kitexclient.Option{
		kitexclient.WithACLRules(container.GetAclRule()),
		kitexclient.WithCloseCallbacks(func() error {
			watcher.DeregisterCallback(keyDegradation)
			return nil
		}),
	}
}

func initDegradationOptions(watcher monitor.ConfigMonitor) (*degradation.Container, int64) {
	degradationContainer := degradation.NewContainer()
	onChangeCallback := func() {
		config := getFileConfig(watcher)
		if config == nil {
			return // config is nil, do nothing, log will be printed in getFileConfig
		}
		degradationContainer.NotifyPolicyChange(config.Degradation)
	}
	keyDegradation := watcher.RegisterCallback(onChangeCallback)
	return degradationContainer, keyDegradation
}
