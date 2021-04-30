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
	"compress/gzip"
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
	"io/ioutil"
)

var _ = Describe("Remote Methods Tests", func() {

	Context("Remote Test", func() {

		var remoteConfig RemoteConfig

		BeforeEach(func() {
			data, err := ioutil.ReadFile("./tests/remote_config_test.json")
			Expect(err).To(BeNil(), "cannot read remote test config")

			err = json.Unmarshal(data, &remoteConfig)
			Expect(err).To(BeNil(), "cannot parse remote test config")
		})

		Describe("Test without input", func() {
			It("Should be succeed", func() {
				outr, outw := io.Pipe()
				err := SendInput2Command(remoteConfig, "echo hi", nil, outw)
				Expect(err).To(BeNil(), "error occured")
				outb, err := io.ReadAll(outr)
				Expect(err).To(BeNil(), "error occured")
				out := string(outb)
				Expect(out).To(Equal("hi\n"), "output different")
			})
		})

		Describe("Test with input", func() {
			It("Should be succeed", func() {
				inr, inw := io.Pipe()
				go func() {
					gw := gzip.NewWriter(inw)
					gw.Write([]byte("hi"))
					gw.Close()
					inw.Close()
				}()
				outr, outw := io.Pipe()
				err := SendInput2Command(remoteConfig, "gzip -d -c -", inr, outw)
				Expect(err).To(BeNil(), "error occured")
				outb, err := io.ReadAll(outr)
				Expect(err).To(BeNil(), "error occured")
				out := string(outb)
				Expect(out).To(Equal("hi"), "output different")
			})
		})
	})
})
