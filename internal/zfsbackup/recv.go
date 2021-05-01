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
	"errors"
	"fmt"
	klog "k8s.io/klog/v2"
	"os"
	"strconv"
	"strings"
)

func ReceiveHandler() error {
	klog.V(5).Infof("inside receive handler")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {

		cmd := scanner.Text()
		klog.V(9).Infof("received command %s", cmd)
		if cmd == "END" {
			fmt.Printf("OK\n")
			klog.V(5).Infof("stopping recv")
			break
		} else if strings.HasPrefix(cmd, "START") {
			klog.V(5).Infof("starting recv")
			parts := strings.Split(cmd, " ")
			if len(parts) != 3 {
				err := errors.New("Command format invalid")
				klog.V(5).Error(err, "cannot start receiver")
				fmt.Printf("ERR\n")
				return err
			}
			bufsize, err := strconv.Atoi(parts[1])
			if err != nil {
				klog.V(5).Error(err, "cannot get buffer size")
				fmt.Printf("ERR\n")
				return err
			}
			br, err := start_buffered_receiver(bufsize, parts[2])
			if err != nil {
				klog.V(5).Error(err, "cannot start receiver")
				fmt.Printf("ERR\n")
				return err
			}
			fmt.Printf("PORT %d\n", br.port)
			err = br.accept()
			if err != nil {
				klog.V(5).Error(err, "cannot accept file")
				fmt.Printf("ERR\n")
				return err
			}
			br.stop()
			fmt.Printf("OK\n")
		} else {
			fmt.Printf("UNKNOWN\n")
		}
	}

	if err := scanner.Err(); err != nil {
		klog.V(5).Error(err, "cannot read input")
		return err
	}

	klog.V(5).Infof("ending receive handler")
	return nil
}
