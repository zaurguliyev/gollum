// Copyright 2015-2017 trivago GmbH
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

package router

import (
	"github.com/trivago/gollum/core"
	"sync"
	"sync/atomic"
)

// RoundRobin stream plugin
// Messages will be sent to one of the producers attached to this router.
// Producers will be switched one-by-one.
type RoundRobin struct {
	core.SimpleRouter
	index         int32
	indexByStream map[core.MessageStreamID]*int32
	mapInitLock   *sync.Mutex
}

func init() {
	core.TypeRegistry.Register(RoundRobin{})
}

// Configure initializes this distributor with values from a plugin config.
func (router *RoundRobin) Configure(conf core.PluginConfigReader) error {
	router.SimpleRouter.Configure(conf)

	router.index = 0
	router.indexByStream = make(map[core.MessageStreamID]*int32)
	router.mapInitLock = new(sync.Mutex)
	return conf.Errors.OrNil()
}

// Enqueue enques a message to the router
func (router *RoundRobin) Enqueue(msg *core.Message) error {
	producers := router.GetProducers()
	if len(producers) == 0 {
		return core.NewModulateResultError("No producers configured for stream %s", router.GetID())
	}
	index := atomic.AddInt32(&router.index, 1) % int32(len(producers))
	producers[index].Enqueue(msg, router.Timeout)
	return nil
}
