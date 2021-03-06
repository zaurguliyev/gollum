// Copyright 2015-2016 trivago GmbH
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

package core

import (
	"github.com/trivago/tgo/tlog"
)

// A Modulator defines a modification or analysis step inside the message
// ModulateResult. It may alter messages or stop the ModulateResult for this message.
type Modulator interface {
	// Modulate processes the given message for analysis or modification.
	// The result of this function defines how and if a message proceeds
	// along the ModulateResult.
	Modulate(msg *Message) ModulateResult
}

// ScopedModulator extends the Modulator interface by adding a log scope.
// This interface is implemented by modulators that are embedded in plugins
// that already have their own scope.
type ScopedModulator interface {
	Modulator

	// SetLogScope defines the log scope for this modulator.
	SetLogScope(log tlog.LogScope)
}

// ModulateResult defines a set of results used to control the message flow
// induced by Modulator actions.
type ModulateResult int

const (
	// ModulateResultContinue indicates that a message can be passed along.
	ModulateResultContinue = ModulateResult(iota)
	// ModulateResultRoute indicates that a message requires routing
	ModulateResultRoute = ModulateResult(iota)
	// ModulateResultDrop has to act like ModulateResultRoute but also
	// indicates that no further modluators should be called.
	ModulateResultDrop = ModulateResult(iota)
	// ModulateResultDiscard indicates that a message should be discarded and
	// that no further modulators should be called.
	ModulateResultDiscard = ModulateResult(iota)
	// ModulateResultHandled is used inside a Modulate chain call when a
	// message has already been processed and does not require further
	// processing (ignore)
	ModulateResultHandled = ModulateResult(iota)
)

// ModulatorArray is a type wrapper to []Modulator to make array of modulators
// compatible with the modulator interface
type ModulatorArray []Modulator

// Modulate calls Modulate on every Modulator in the array and react according
// to the definition of each ModulateResult state.
func (modulators ModulatorArray) Modulate(msg *Message) ModulateResult {
	action := ModulateResultContinue
	for _, modulator := range modulators {
		switch modRes := modulator.Modulate(msg); modRes {
		case ModulateResultDiscard, ModulateResultDrop:
			return modRes

		case ModulateResultRoute:
			action = modRes
		}
	}
	return action
}
