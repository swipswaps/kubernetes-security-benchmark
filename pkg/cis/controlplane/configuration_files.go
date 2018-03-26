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

package controlplane

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mesosphere/kubernetes-security-benchmark/pkg/framework"
	. "github.com/mesosphere/kubernetes-security-benchmark/pkg/matcher"
	"github.com/mesosphere/kubernetes-security-benchmark/pkg/util"
)

func configFilePermissionsContext(directory, fileName string, specFunc func(filePath string)) {
	Context("", func() {
		var filePath string

		BeforeEach(func() {
			filePath = filepath.Join(directory, fileName)
			_, err := os.Stat(filePath)
			if err != nil {
				if os.IsNotExist(err) {
					Skip(fmt.Sprintf("%s does not exist", filePath))
				}
				Expect(err).NotTo(HaveOccurred())
			}
			Expect(filePath).To(BeARegularFile())
		})

		specFunc(filePath)
	})
}

func ConfigurationFiles(index, subIndex int, missingProcessFunc framework.MissingProcessHandlerFunc) {
	Context("", func() {
		kubelet := framework.New("kubelet", missingProcessFunc)
		BeforeEach(kubelet.BeforeEach)

		Context("", func() {
			var kubeletPodManifestPath string

			BeforeEach(func() {
				kmp, err := util.FlagValueFromProcess(kubelet.Process, "pod-manifest-path")
				Expect(err).NotTo(HaveOccurred())
				if kmp == "" {
					Skip(fmt.Sprintf("Flag --%s is unset", "pod-manifest-path"))
				}
				Expect(kmp).To(BeADirectory())
				kubeletPodManifestPath = kmp.(string)
			})

			configFilePermissionsContext(kubeletPodManifestPath, "api-server.yaml", func(filePath string) {
				It(fmt.Sprintf("[%d.%d.1] Ensure that the API server pod specification file permissions are set to 644 or more restrictive", index, subIndex), func() {
					Expect(filePath).To(HavePermissionsNumerically("<=", os.FileMode(0644)))
				})

				It(fmt.Sprintf("[%d.%d.2] Ensure that the API server pod specification file ownership is set to root:root", index, subIndex), func() {
					Expect(filePath).To(BeOwnedBy("root", "root"))
				})
			})

			configFilePermissionsContext(kubeletPodManifestPath, "kube-controller-manager.yaml", func(filePath string) {
				It(fmt.Sprintf("[%d.%d.3] Ensure that the controller manager pod specification file permissions are set to 644 or more restrictive", index, subIndex), func() {
					Expect(filePath).To(HavePermissionsNumerically("<=", os.FileMode(0644)))
				})

				It(fmt.Sprintf("[%d.%d.4] Ensure that the controller manager pod specification file ownership is set to root:root", index, subIndex), func() {
					Expect(filePath).To(BeOwnedBy("root", "root"))
				})
			})

			configFilePermissionsContext(kubeletPodManifestPath, "kube-scheduler.yaml", func(filePath string) {
				It(fmt.Sprintf("[%d.%d.5] Ensure that the scheduler pod specification file permissions are set to 644 or more restrictive", index, subIndex), func() {
					Expect(filePath).To(HavePermissionsNumerically("<=", os.FileMode(0644)))
				})

				It(fmt.Sprintf("[%d.%d.6] Ensure that the scheduler pod specification file ownership is set to root:root", index, subIndex), func() {
					Expect(filePath).To(BeOwnedBy("root", "root"))
				})
			})

			configFilePermissionsContext(kubeletPodManifestPath, "etcd.yaml", func(filePath string) {
				It(fmt.Sprintf("[%d.%d.7] Ensure that the etcd pod specification file permissions are set to 644 or more restrictive", index, subIndex), func() {
					Expect(filePath).To(HavePermissionsNumerically("<=", os.FileMode(0644)))
				})

				It(fmt.Sprintf("[%d.%d.8] Ensure that the etcd pod specification file ownership is set to root:root", index, subIndex), func() {
					Expect(filePath).To(BeOwnedBy("root", "root"))
				})
			})
		})

		Context("", func() {
			var cniConfDir string

			BeforeEach(func() {
				networkPlugin, err := util.FlagValueFromProcess(kubelet.Process, "network-plugin")
				Expect(err).NotTo(HaveOccurred())
				if networkPlugin == "" {
					Skip("Flag --network-plugin is unset")
				}
				if networkPlugin != "cni" {
					Skip("Flag --network-plugin is not set to cni")
				}

				ccd, err := util.FlagValueFromProcess(kubelet.Process, "cni-conf-dir")
				Expect(err).NotTo(HaveOccurred())
				if ccd == "" {
					Skip("Flag --cni-conf-dir is unset")
				}
				Expect(ccd).To(BeADirectory())

				cniConfDir = ccd.(string)
			})

			It(fmt.Sprintf("[%d.%d.9] Ensure that the Container Network Interface file permissions are set to 644 or more restrictive", index, subIndex), func() {
				err := filepath.Walk(cniConfDir, func(path string, info os.FileInfo, err error) error {
					ExpectWithOffset(1, err).NotTo(HaveOccurred())
					if !info.IsDir() {
						Expect(path).To(HavePermissionsNumerically("<=", os.FileMode(0644)))
					}
					return nil
				})
				Expect(err).NotTo(HaveOccurred())
			})

			It(fmt.Sprintf("[%d.%d.10] Ensure that the Container Network Interface file ownership is set to root:root", index, subIndex), func() {
				err := filepath.Walk(cniConfDir, func(path string, info os.FileInfo, err error) error {
					ExpectWithOffset(1, err).NotTo(HaveOccurred())
					if !info.IsDir() {
						Expect(path).To(BeOwnedBy("root", "root"))
					}
					return nil
				})
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Context("", func() {
		etcd := framework.New("etcd", missingProcessFunc)
		BeforeEach(etcd.BeforeEach)

		Context("", func() {
			var etcdDataDir string

			BeforeEach(func() {
				edd, err := util.FlagValueFromProcess(etcd.Process, "data-dir")
				Expect(err).NotTo(HaveOccurred())
				if edd == "" {
					Skip("Flag --data-dir is unset")
				}
				Expect(edd).To(BeADirectory())

				etcdDataDir = edd.(string)
			})

			It(fmt.Sprintf("[%d.%d.11] Ensure that the etcd data directory permissions are set to 700 or more restrictive", index, subIndex), func() {
				Expect(etcdDataDir).To(HavePermissionsNumerically("<=", os.FileMode(0700)))
			})

			It(fmt.Sprintf("[%d.%d.12] Ensure that the etcd data directory ownership is set to etcd:etcd", index, subIndex), func() {
				Expect(etcdDataDir).To(BeOwnedBy("etcd", "etcd"))
			})
		})
	})

	Context("", func() {
		const adminFilePath = "/etc/kubernetes/admin.conf"
		BeforeEach(func() {
			_, err := os.Stat(adminFilePath)
			if err != nil {
				if os.IsNotExist(err) {
					Skip(fmt.Sprintf("%s does not exist", adminFilePath))
				}
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It(fmt.Sprintf("[%d.%d.13] Ensure that the admin.conf file permissions are set to 644 or more restrictive", index, subIndex), func() {
			Expect(adminFilePath).To(HavePermissionsNumerically("<=", os.FileMode(0644)))
		})

		It(fmt.Sprintf("[%d.%d.14] Ensure that the admin.conf file ownership is set to root:root", index, subIndex), func() {
			Expect(adminFilePath).To(BeOwnedBy("root", "root"))
		})
	})

	Context("", func() {
		const schedulerFilePath = "/etc/kubernetes/scheduler.conf"
		BeforeEach(func() {
			_, err := os.Stat(schedulerFilePath)
			if err != nil {
				if os.IsNotExist(err) {
					Skip(fmt.Sprintf("%s does not exist", schedulerFilePath))
				}
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It(fmt.Sprintf("[%d.%d.15] Ensure that the scheduler.conf file permissions are set to 644 or more restrictive", index, subIndex), func() {
			Expect(schedulerFilePath).To(HavePermissionsNumerically("<=", os.FileMode(0644)))
		})

		It(fmt.Sprintf("[%d.%d.16] Ensure that the scheduler.conf file ownership is set to root:root", index, subIndex), func() {
			Expect(schedulerFilePath).To(BeOwnedBy("root", "root"))
		})
	})

	Context("", func() {
		const controllerManagerFilePath = "/etc/kubernetes/controller-manager.conf"
		BeforeEach(func() {
			_, err := os.Stat(controllerManagerFilePath)
			if err != nil {
				if os.IsNotExist(err) {
					Skip(fmt.Sprintf("%s does not exist", controllerManagerFilePath))
				}
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It(fmt.Sprintf("[%d.%d.17] Ensure that the controller-manager.conf file permissions are set to 644 or more restrictive", index, subIndex), func() {
			Expect(controllerManagerFilePath).To(HavePermissionsNumerically("<=", os.FileMode(0644)))
		})

		It(fmt.Sprintf("[%d.%d.18] Ensure that the controller-manager.conf file ownership is set to root:root", index, subIndex), func() {
			Expect(controllerManagerFilePath).To(BeOwnedBy("root", "root"))
		})
	})
}