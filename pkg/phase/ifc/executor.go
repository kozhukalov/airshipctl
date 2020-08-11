/*
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package ifc

import (
	"io"
	"time"

	"opendev.org/airship/airshipctl/pkg/document"
	"opendev.org/airship/airshipctl/pkg/environment"
	"opendev.org/airship/airshipctl/pkg/events"
	"opendev.org/airship/airshipctl/pkg/k8s/kubeconfig"
)

// Executor interface should be implemented by each runner
type Executor interface {
	Run(RunOptions) <-chan events.Event
	Render(io.Writer, RenderOptions) error
	Validate() error
	Wait(WaitOptions) <-chan events.Event
}

// RunOptions holds options for run method
type RunOptions struct {
	Debug  bool
	DryRun bool

	Timeout time.Duration
}

// RenderOptions is empty for now, but may hold things like format in future
type RenderOptions struct{}

// WaitOptions holds only timeout now, but may be extended in the future
type WaitOptions struct {
	Timeout time.Duration
}

// ExecutorFactory for executor instantiation
// First argument is document object which represents executor
// configuration.
// Second argument is document bundle used by executor.
// Third argument airship configuration settings since each phase
// has to be aware of execution context and global settings
type ExecutorFactory func(
	document.Document,
	document.Bundle,
	*environment.AirshipCTLSettings,
	kubeconfig.Provider,
) (Executor, error)
