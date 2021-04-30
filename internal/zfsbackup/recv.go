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
	"bufio"
	"fmt"
	klog "k8s.io/klog/v2"
	"os"
)

func ReceiveHandler() error {
	klog.V(5).Infof("inside receive handler")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		if err := scanner.Err(); err != nil {
			klog.V(5).Error(err, "cannot read input")
			return err
		}
		fmt.Printf("hello: %s\n", scanner.Text())
	}

	klog.V(5).Infof("ending receive handler")
	return nil
}
