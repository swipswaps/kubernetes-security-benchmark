// Copyright © 2018 Jimmi Dyson <jdyson@mesosphere.com>
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

package cmd

import (
	"flag"
	"testing"

	"github.com/onsi/ginkgo/config"
	"github.com/spf13/cobra"

	"github.com/mesosphere/kubernetes-security-benchmark/pkg/cis"
)

// cisCmd represents the cis command
var cisCmd = &cobra.Command{
	Use:   "cis",
	Short: "Run Kubernetes CIS Benchmark tests",
	Long:  `Run Kubernetes CIS Benchmark tests.`,
	Run: func(cmd *cobra.Command, args []string) {
		testing.Main(func(pat, str string) (bool, error) { return true, nil },
			[]testing.InternalTest{{"KubernetesCISBenchmark", cis.CISBenchmark}},
			[]testing.InternalBenchmark{},
			[]testing.InternalExample{})
	},
}

func init() {

	ginkgoFlagSet := flag.NewFlagSet("ginkgo", flag.ContinueOnError)
	config.Flags(ginkgoFlagSet, "ginkgo", false)

	cisCmd.Flags().AddGoFlagSet(ginkgoFlagSet)

	rootCmd.AddCommand(cisCmd)
}
