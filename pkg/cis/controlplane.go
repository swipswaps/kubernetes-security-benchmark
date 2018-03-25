// Copyright © 2018 Jimmi Dyson <jimmidyson@gmail.com>
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

package cis

import (
	. "github.com/onsi/ginkgo"

	"github.com/mesosphere/kubernetes-security-benchmark/pkg/cis/controlplane"
)

const CISVersion = "1.2.0"

func CISDescribe(text string, body func()) bool {
	return Describe("[CIS 1.2.0] "+text, body)
}

var _ = CISDescribe("[1] Control plane", func() {
	Context("[1.1] API Server", func() {
		controlplane.APIServer(1, 1)
	})
	Context("[1.2] Scheduler", func() {
		controlplane.Scheduler(1, 2)
	})
	Context("[1.3] Controller Manager", func() {
		controlplane.ControllerManager(1, 3)
	})
})
