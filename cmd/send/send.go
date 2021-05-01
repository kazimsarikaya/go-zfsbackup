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

package sendcmd

import (
	"github.com/kazimsarikaya/go-zfsbackup/internal/zfsbackup"
	"github.com/spf13/cobra"
	klog "k8s.io/klog/v2"
)

var (
	sendCmd = &cobra.Command{
		Use:   "send",
		Short: "Sends zfs backup to remote",
		Long: `Receives commands from stdin pipelined by ssh which
initited by send command`,
		RunE: func(cmd *cobra.Command, args []string) error {

			hostPort, err := cmd.Flags().GetString("host-port")
			if err != nil {
				return err
			}
			klog.V(9).Infof("host-port: %v", hostPort)

			hostKey, err := cmd.Flags().GetString("host-key")
			if err != nil {
				return err
			}
			klog.V(9).Infof("host-key: %v", hostKey)

			user, err := cmd.Flags().GetString("user")
			if err != nil {
				return err
			}
			klog.V(9).Infof("user: %v", user)

			keyFile, err := cmd.Flags().GetString("key-file")
			if err != nil {
				return err
			}
			klog.V(9).Infof("key-file: %v", keyFile)

			executablePath, err := cmd.Flags().GetString("executable-path")
			if err != nil {
				return err
			}
			klog.V(9).Infof("executable-path: %v", executablePath)

			src_ds, err := cmd.Flags().GetString("source-dataset")
			if err != nil {
				return err
			}
			klog.V(9).Infof("source-dataset: %v", src_ds)

			dst_ds, err := cmd.Flags().GetString("destination-dataset")
			if err != nil {
				return err
			}
			klog.V(9).Infof("destination-dataset: %v", dst_ds)

			rcfg := &zfsbackup.RemoteConfig{
				HostPort:       hostPort,
				HostKey:        hostKey,
				User:           user,
				KeyFile:        keyFile,
				ExecutablePath: executablePath,
			}
			klog.V(5).Infof("Send started")
			return zfsbackup.SendHandler(rcfg, &src_ds, &dst_ds)
		},
	}
)

func GetSendCmd() *cobra.Command {
	sendCmd.Flags().StringP("host-port", "", "localhost:22", "Destination with format host:port")
	sendCmd.Flags().StringP("host-key", "", "", "SSH host public key with format 'host keytype key'")
	sendCmd.Flags().StringP("user", "", "root", "Destination root user")
	sendCmd.Flags().StringP("key-file", "", "/root/.ssh/id_ecdsa", "Users ssh private key")
	sendCmd.Flags().StringP("executable-path", "", "/usr/local/sbin/zfsbackup", "Destination zfs executable path")
	sendCmd.Flags().StringP("source-dataset", "s", "", "Source zfs dataset to send for backup")
	sendCmd.Flags().StringP("destination-dataset", "d", "", "Destination zfs dataset for receive backup")
	return sendCmd
}
