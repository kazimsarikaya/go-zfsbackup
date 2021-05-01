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
	"crypto/sha256"
	"fmt"
	"io"
	klog "k8s.io/klog/v2"
	"os"
)

var ownExecutable string

func SetOwnExecutable(ex string) {
	ownExecutable = ex
}

func get_local_hash_of_file(file string) (string, error) {
	file_r, err := os.Open(file)
	if err != nil {
		klog.V(5).Error(err, "cannot open local file for hash256 sum")
		return "", err
	}

	h := sha256.New()
	if _, err := io.Copy(h, file_r); err != nil {
		klog.V(5).Error(err, "cannot calculate hash of local file")
		return "", err
	}

	source_hash := fmt.Sprintf("%x", h.Sum(nil))
	return source_hash, nil
}
