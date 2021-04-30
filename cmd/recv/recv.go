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

package recvcmd

import (
	"github.com/kazimsarikaya/go-zfsbackup/internal/zfsbackup"
	"github.com/spf13/cobra"
	klog "k8s.io/klog/v2"
)

var (
	recvCmd = &cobra.Command{
		Use:   "recv",
		Short: "Receives recv commands from stdin",
		Long: `Receives commands from stdin pipelined by ssh which
initited by send command`,
		RunE: func(cmd *cobra.Command, args []string) error {
			klog.V(5).Infof("Recevie started")
			return zfsbackup.ReceiveHandler()
		},
	}
)

func GetRecvCmd() *cobra.Command {
	return recvCmd
}
