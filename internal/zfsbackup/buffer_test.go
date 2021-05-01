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
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
	"net"
	"os"
)

var _ = Describe("Buffer Methods Tests", func() {

	Context("Receive File", func() {

		Describe("Test receving file", func() {
			It("Should be succeed", func() {
				src_f, err := os.CreateTemp(os.TempDir(), "buffer-test-src.*.bin")
				Expect(err).To(BeNil(), "cannot create test file")
				defer os.Remove(src_f.Name()) // clean up
				r_data, err := os.Open("/dev/urandom")
				Expect(err).To(BeNil(), "cannot open urandom")
				io.CopyN(src_f, r_data, 1024*1024)
				r_data.Close()
				src_f.Close()
				src_hash, err := get_local_hash_of_file(src_f.Name())
				Expect(err).To(BeNil(), "cannot calculate hash of src file")

				dst_f, err := os.CreateTemp(os.TempDir(), "buffer-test-dst.*.bin")
				Expect(err).To(BeNil(), "cannot create test file")
				defer os.Remove(dst_f.Name()) // clean up
				br, err := start_buffered_receiver(1024, fmt.Sprintf("file:%s", dst_f.Name()))

				wait_ch := make(chan string)

				go func() {
					err := br.accept()
					Expect(err).To(BeNil(), "cannot recevie test file")
					br.stop()
					dst_hash, err := get_local_hash_of_file(dst_f.Name())
					Expect(err).To(BeNil(), "cannot calculate hash of dst file")
					Expect(dst_hash).To(Equal(src_hash), "file content mismatch")
					wait_ch <- "ok"
				}()
				conn, err := net.Dial("tcp", fmt.Sprintf(":%d", br.port))
				defer conn.Close()
				Expect(err).To(BeNil(), "cannot connect")
				src_in, err := os.Open(src_f.Name())
				Expect(err).To(BeNil(), "cannot open source file")
				_, err = io.Copy(conn, src_in)
				conn.Close()
				Expect(err).To(BeNil(), "cannot send file content")
				res := <-wait_ch
				Expect(res).To(Equal("ok"), "test failed")
			})
		})
	})
})
