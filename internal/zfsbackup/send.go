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
	"io"
	klog "k8s.io/klog/v2"
	"net"
	"os"
	"strconv"
	"strings"
)

func SendHandler(rcfg *RemoteConfig, src, dst *string, buffer_size int) error {
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

	cmd := fmt.Sprintf("START %d %s\nEND\n", buffer_size, *dst)
	klog.V(9).Infof("remote input: %q", cmd)
	cmd_in := strings.NewReader(cmd)

	out_r, out_w := io.Pipe()

	err_ch := make(chan error)

	go func() {
		err_ch <- SendInput2Command(rcfg, fmt.Sprintf("%s recv", rcfg.ExecutablePath), cmd_in, out_w)
	}()

	scanner := bufio.NewScanner(out_r)
	for {
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			klog.V(5).Error(err, "remote answer error")
			continue
		}
		res := scanner.Text()
		if res == "FIN" {
			klog.V(5).Infof("session ended")
			break
		} else if strings.HasPrefix(res, "PORT") {
			port, err := strconv.Atoi(res[5:])
			if err != nil {
				klog.V(5).Error(err, "cannot parse port")
				break
			}
			host_port := rcfg.HostPort[:strings.Index(rcfg.HostPort, ":")]
			host_port = fmt.Sprintf("%s:%d", host_port, port)
			err = send_source(*src, host_port, buffer_size)
			if err != nil {
				klog.V(5).Error(err, "cannot send source")
				break
			}
		} else if res == "ERR" {
			err := errors.New("remote endpoint error")
			klog.V(5).Error(err, "error occured at remote endpoint")
		} else if res == "OK" {
			klog.V(5).Infof("command ended without error")
		} else {
			klog.V(5).Error(errors.New("unknown answer"), "remote answer unknown")
			break
		}
	}

	err = <-err_ch
	if err != nil {
		klog.V(5).Error(err, "ssh command failed")
		return err
	}

	klog.V(5).Infof("ending send handler")
	return nil
}

func send_source(src, host_port string, buffer_size int) error {
	klog.V(5).Infof("sending %s to %s", src, host_port)
	conn, err := net.Dial("tcp", host_port)
	if err != nil {
		klog.V(5).Error(err, "cannot connect remote endpoint")
		return err
	}
	defer conn.Close()
	if err != nil {
		klog.V(5).Error(err, "cannot connect remote end")
		return err
	}

	if strings.HasPrefix(src, "file:") {
		file := src[5:]
		klog.V(5).Infof("sending file %s", file)
		in_file, err := os.Open(file)
		if err != nil {
			klog.V(5).Error(err, "cannot open input")
			return err
		}
		defer in_file.Close()
		buf_in_file := bufio.NewReaderSize(in_file, buffer_size)
		_, err = io.Copy(conn, buf_in_file)
		if err != nil {
			klog.V(5).Error(err, "cannot send source file")
			return err
		}
	} else if strings.HasPrefix(src, "zfs:") {
		// TODO: implement zfs
	} else {
		return errors.New("unknown source")
	}
	klog.V(5).Infof("%s sended", src)
	return nil
}
