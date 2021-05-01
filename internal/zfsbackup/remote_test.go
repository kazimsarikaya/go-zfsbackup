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
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
	"io/ioutil"
	"os"
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
				err := SendInput2Command(&remoteConfig, "echo hi", nil, outw)
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
				err := SendInput2Command(&remoteConfig, "gzip -d -c -", inr, outw)
				Expect(err).To(BeNil(), "error occured")
				outb, err := io.ReadAll(outr)
				Expect(err).To(BeNil(), "error occured")
				out := string(outb)
				Expect(out).To(Equal("hi"), "output different")
			})
		})

		Describe("Test file existence at remote", func() {
			It("Should be succeed", func() {
				f, err := os.CreateTemp(os.TempDir(), "file-exists.*.txt")
				Expect(err).To(BeNil(), "cannot create test file")
				defer os.Remove(f.Name()) // clean up
				found, err := is_file_exists_at_remote(&remoteConfig, f.Name())
				Expect(err).To(BeNil(), "error occured while testing")
				Expect(found).To(BeTrue(), "file not found")
				found, err = is_file_exists_at_remote(&remoteConfig, fmt.Sprintf("%s/not-exists-file", os.TempDir()))
				Expect(err).To(BeNil(), "error occured while testing")
				Expect(found).NotTo(BeTrue(), "file not found")
			})
		})

		Describe("Test file existence at remote", func() {
			It("Should be succeed", func() {
				f, err := os.CreateTemp(os.TempDir(), "file-exists.*.txt")
				Expect(err).To(BeNil(), "cannot create test file")
				defer os.Remove(f.Name()) // clean up
				found, err := is_file_exists_at_remote(&remoteConfig, f.Name())
				Expect(err).To(BeNil(), "error occured while testing")
				Expect(found).To(BeTrue(), "file not found")
				found, err = is_file_exists_at_remote(&remoteConfig, fmt.Sprintf("%s/not-exists-file", os.TempDir()))
				Expect(err).To(BeNil(), "error occured while testing")
				Expect(found).NotTo(BeTrue(), "file not found")
			})
		})

		Describe("Test remote file's sha256", func() {
			It("Should be succeed", func() {
				f, err := os.CreateTemp(os.TempDir(), "remote-sha256-test.*.txt")
				Expect(err).To(BeNil(), "cannot create test file")
				defer os.Remove(f.Name()) // clean up
				_, err = f.Write([]byte("test content\n"))
				Expect(err).To(BeNil(), "cannot write test file")
				err = f.Close()
				Expect(err).To(BeNil(), "cannot close test file")
				hash, err := get_file_hash_at_remote(&remoteConfig, f.Name())
				Expect(err).To(BeNil(), "cannot calculate hash of test file")
				Expect(hash).To(Equal("a1fff0ffefb9eace7230c24e50731f0a91c62f9cefdfe77121c2f607125dffae"), "verify of hash failed")
			})
		})

		Describe("Test remote sending file", func() {
			It("Should be succeed", func() {
				f_src, err := os.CreateTemp(os.TempDir(), "remote-send-test-src.*.txt")
				Expect(err).To(BeNil(), "cannot create test file")
				defer os.Remove(f_src.Name()) // clean up
				_, err = f_src.Write([]byte("test content\n"))
				Expect(err).To(BeNil(), "cannot write test file")
				err = f_src.Close()
				Expect(err).To(BeNil(), "cannot close test file")

				f_dst, err := os.CreateTemp(os.TempDir(), "remote-send-test-dst.*.txt")
				Expect(err).To(BeNil(), "cannot create test file")
				defer os.Remove(f_dst.Name()) // clean up

				err = send_file_to_remote(&remoteConfig, f_src.Name(), f_dst.Name())
				Expect(err).To(BeNil(), "cannot calculate hash of test file")

				hash, err := get_file_hash_at_remote(&remoteConfig, f_dst.Name())
				Expect(err).To(BeNil(), "cannot calculate hash of test file")
				Expect(hash).To(Equal("a1fff0ffefb9eace7230c24e50731f0a91c62f9cefdfe77121c2f607125dffae"), "verify of hash failed")
			})
		})
	})
})
