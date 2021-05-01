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
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	klog "k8s.io/klog/v2"
	"os"
)

type RemoteConfig struct {
	HostPort       string
	HostKey        string
	User           string
	KeyFile        string
	ExecutablePath string
}

func SendInput2Command(remoteConfig *RemoteConfig, cmd string, in io.ReadCloser, out io.WriteCloser) error {
	klog.V(8).Infof("SendInput2Command called")
	_, _, hostPubKey, _, _, err := ssh.ParseKnownHosts([]byte(remoteConfig.HostKey))
	if err != nil {
		klog.V(0).Error(err, "unable to parse host public key")
		return err
	}
	klog.V(8).Infof("host public key pared")

	key, err := ioutil.ReadFile(remoteConfig.KeyFile)
	if err != nil {
		klog.V(0).Error(err, "unable to read private key")
		return err
	}
	klog.V(8).Infof("private key readed")

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		klog.V(0).Error(err, "unable to parse private key")
		return err
	}
	klog.V(8).Infof("private key parsed")

	config := &ssh.ClientConfig{
		User: remoteConfig.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.FixedHostKey(hostPubKey),
	}

	client, err := ssh.Dial("tcp", remoteConfig.HostPort, config)
	if err != nil {
		klog.V(0).Error(err, "unable to connect")
		return err
	}
	defer client.Close()
	klog.V(5).Infof("ssh connection established")

	sess, err := client.NewSession()
	if err != nil {
		klog.V(0).Error(err, "unable to create session")
		return err
	}
	defer sess.Close()
	klog.V(5).Infof("ssh session started")

	klog.Flush()

	r, err := sess.StdoutPipe()
	if err != nil {
		klog.V(0).Error(err, "unable to create ssh output pipe")
		return err
	}

	w, err := sess.StdinPipe()
	if err != nil {
		klog.V(0).Error(err, "unable to create ssh input pipe")
		return err
	}

	err = sess.Start(cmd)
	if err != nil {
		klog.V(0).Error(err, "unable to start ssh command")
		return err
	}
	klog.V(5).Infof("ssh command started")

	if in != nil {
		go func() {
			klog.V(5).Infof("start coping input")
			io.Copy(w, in)
			w.Close()
			klog.V(5).Infof("coping input finished")
		}()
	}

	if out != nil {
		go func() {
			klog.V(5).Infof("start coping output")
			io.Copy(out, r)
			out.Close()
			klog.V(5).Infof("ssh output gathered")
		}()
	}

	err = sess.Wait()
	if err != nil {
		klog.V(0).Error(err, "ssh command failed")
		return err
	}
	klog.V(5).Infof("ssh command ended")

	return nil
}

func is_file_exists_at_remote(rcfg *RemoteConfig, file string) (bool, error) {
	outr, outw := io.Pipe()
	err := SendInput2Command(rcfg, fmt.Sprintf("[[ -f %s ]] && echo -n found || echo -n notfound", file), nil, outw)
	if err != nil {
		klog.V(5).Error(err, "cannot verify file existance")
		return false, err
	}
	outb, err := io.ReadAll(outr)
	if err != nil {
		klog.V(5).Error(err, "cannot verify file existance")
		return false, err
	}

	out := string(outb)
	if out == "found" {
		return true, nil
	}
	return false, nil
}

func send_file_to_remote(rcfg *RemoteConfig, src_file, dst_file string) error {
	klog.V(5).Infof("sending file %s", src_file)
	fi, err := os.Stat(src_file)
	if err != nil {
		klog.V(5).Error(err, "cannot stat source file for perms")
		return err
	}
	klog.V(9).Infof("source file's permision is %o", fi.Mode().Perm())
	file_r, err := os.Open(src_file)
	if err != nil {
		klog.V(5).Error(err, "cannot open for reading file")
		return err
	}
	err = SendInput2Command(rcfg, fmt.Sprintf("cat > %s; chmod %o %s", dst_file, fi.Mode().Perm(), dst_file), file_r, nil)
	if err != nil {
		klog.V(5).Error(err, "cannot send file")
		return err
	}
	return nil
}

func get_file_hash_at_remote(rcfg *RemoteConfig, file string) (string, error) {
	klog.V(9).Infof("calculate remote file %s sha256 sum", file)
	outr, outw := io.Pipe()
	err := SendInput2Command(rcfg, fmt.Sprintf("openssl sha256 -r  %s", file), nil, outw)
	if err != nil {
		klog.V(5).Error(err, "cannot calculate remote file's sha256 sum")
		return "", err
	}
	outb, err := io.ReadAll(outr)
	if err != nil {
		klog.V(5).Error(err, "cannot read remote file's sha256 sum")
		return "", err
	}
	outs := string(outb)
	if len(outs) >= 64 {
		dest_hash := outs[:64]
		return dest_hash, nil
	}
	return "", errors.New("remote hash calculation error")
}
