// +build linux

/*
Copyright 2021 Kazım SARIKAYA

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package zfsbackup

import (
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
)

var _ = Describe("Send Methods Tests", func() {

	Context("Send Test", func() {

		var remoteConfig RemoteConfig

		BeforeEach(func() {
			data, err := ioutil.ReadFile("./tests/remote_config_test.json")
			Expect(err).To(BeNil(), "cannot read remote test config")

			err = json.Unmarshal(data, &remoteConfig)
			Expect(err).To(BeNil(), "cannot parse remote test config")
		})

		Describe("Test sending executable", func() {
			It("Should be succeed", func() {
				err := SendHandler(&remoteConfig, nil, nil)
				Expect(err).To(BeNil(), "cannot send executable")
			})
		})
	})
})
