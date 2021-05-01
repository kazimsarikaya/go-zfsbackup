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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Utils Methods Tests", func() {

	Context("Send Test", func() {

		Describe("Test local file's sha256", func() {
			It("Should be succeed", func() {
				f, err := os.CreateTemp(os.TempDir(), "local-sha256-test.*.txt")
				Expect(err).To(BeNil(), "cannot create test file")
				defer os.Remove(f.Name()) // clean up
				_, err = f.Write([]byte("test content\n"))
				Expect(err).To(BeNil(), "cannot write test file")
				err = f.Close()
				Expect(err).To(BeNil(), "cannot close test file")
				hash, err := get_local_hash_of_file(f.Name())
				Expect(err).To(BeNil(), "cannot calculate hash of test file")
				Expect(hash).To(Equal("a1fff0ffefb9eace7230c24e50731f0a91c62f9cefdfe77121c2f607125dffae"), "verify of hash failed")
			})
		})
	})
})
