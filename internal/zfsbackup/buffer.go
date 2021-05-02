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
	"io"
	klog "k8s.io/klog/v2"
	"net"
	"os"
	"strings"
	"time"
)

type bufferedReceviver struct {
	bufsize     int
	port        int
	listener    net.Listener
	destination string
}

func start_buffered_receiver(bufsize int, dest string) (*bufferedReceviver, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		klog.V(5).Error(err, "cannot start listener")
		return nil, err
	}

	port := listener.Addr().(*net.TCPAddr).Port

	return &bufferedReceviver{
		bufsize:     bufsize,
		port:        port,
		listener:    listener,
		destination: dest,
	}, nil
}

func (br *bufferedReceviver) stop() {
	br.listener.Close()
}

func (br *bufferedReceviver) accept() error {
	timer := time.NewTimer(60 * time.Second)
	go func() {
		<-timer.C
		klog.V(5).Infof("any connection comed")
		br.listener.Close()
	}()
	conn, err := br.listener.Accept()
	if err != nil {
		klog.V(5).Error(err, "cannot accept connection")
		return err
	}
	timer.Stop()

	if strings.HasPrefix(br.destination, "file:") {
		d := br.destination[5:]
		klog.V(5).Infof("receiving file %s", d)
		f, err := os.Create(d)
		if err != nil {
			klog.V(5).Error(err, "cannot create destination file")
			return err
		}
		conn_r := bufio.NewReaderSize(conn, br.bufsize)
		r, w := io.Pipe()
		buf_r := bufio.NewReaderSize(r, br.bufsize)
		go func() {
			io.Copy(w, conn_r)
			w.Close()
		}()
		io.Copy(f, buf_r)
		f.Close()
		klog.V(5).Infof("receiving ended")
	} else if strings.HasPrefix(br.destination, "zfs:") {
		d := br.destination[4:]
		klog.V(5).Infof("receiving zfs %s", d)

		// TODO: implement zfs receive

		klog.V(5).Infof("receiving ended")
	} else {
		err = errors.New("unknown destination")
		klog.V(5).Error(err, "unknown destination")
		return err
	}

	return nil
}
