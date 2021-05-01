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
	"errors"
	klog "k8s.io/klog/v2"
)

func SendHandler(rcfg *RemoteConfig, srcdataset, dstdataset *string) error {
	klog.V(5).Infof("inside send handler")
	found, err := is_file_exists_at_remote(rcfg, rcfg.ExecutablePath)
	if err != nil {
		return nil
	}
	if !found {
		err = send_file_to_remote(rcfg, ownExecutable, rcfg.ExecutablePath)
		if err != nil {
			return err
		}
	} else {
		dest_hash, err := get_file_hash_at_remote(rcfg, rcfg.ExecutablePath)
		if err != nil {
			return err
		}

		source_hash, err := get_local_hash_of_file(ownExecutable)
		if err != nil {
			return err
		}

		klog.V(5).Infof("source_hash: %s dest_hash: %s", source_hash, dest_hash)
		if dest_hash != source_hash {
			return errors.New("source and destination executables different")
		}
		klog.V(5).Infof("source and destination executables are same")
	}

	klog.V(5).Infof("ending send handler")
	return nil
}
