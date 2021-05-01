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

package main

import (
	"bytes"
	"flag"
	"fmt"
	. "github.com/kazimsarikaya/go-zfsbackup/cmd/recv"
	. "github.com/kazimsarikaya/go-zfsbackup/cmd/send"
	"github.com/kazimsarikaya/go-zfsbackup/internal/zfsbackup"
	"github.com/spf13/cobra"
	cobradoc "github.com/spf13/cobra/doc"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	klog "k8s.io/klog/v2"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	rootCmd = &cobra.Command{
		Use:   "zfsbackup",
		Short: "A simple zfs backup tool",
		Long:  `A beatiful zfs backuptool`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd)
		},
	}
	version   = ""
	buildTime = ""
	goVersion = ""

	readMeHeader = `# ZFSBACKUP

ZFSBACKUP is a tool creates snapshots sends target machine via ssh and retain old
backups using configuration

# Usage

`

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("ZFS Backup tool\n")
			fmt.Printf("Version: %v\n", version)
			fmt.Printf("Build Time: %v\n", buildTime)
			fmt.Printf("%v\n", goVersion)
		},
	}

	readmeCmd = &cobra.Command{
		Use:    "readme",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {

			lh := func(name string) string {
				base := strings.TrimSuffix(name, path.Ext(name))
				base = strings.Replace(base, "_", "-", -1)
				return "#" + base
			}

			var genReadMe func(cmd *cobra.Command, out *bytes.Buffer) error
			genReadMe = func(cmd *cobra.Command, out *bytes.Buffer) error {
				cmd.DisableAutoGenTag = true
				if err := cobradoc.GenMarkdownCustom(cmd, out, lh); err != nil {
					return err
				}
				for _, subcmd := range cmd.Commands() {
					if err := genReadMe(subcmd, out); err != nil {
						return err
					}
				}
				return nil
			}
			out := new(bytes.Buffer)
			genReadMe(rootCmd, out)
			if rm, err := os.Create("README.md"); err == nil {
				rm.Write([]byte(readMeHeader))
				rm.Write(out.Bytes())
				rm.Close()
			} else {
				klog.V(0).Error(err, "cannot generate readme")
			}
		},
	}
)

func init() {

	klog.InitFlags(nil)

	ex, _ := os.Executable()
	zfsbackup.SetOwnExecutable(ex)

	rootCmd.PersistentFlags().StringP("config", "", "", "configuration file")
	pflag.CommandLine.AddGoFlag(flag.CommandLine.Lookup("v"))
	pflag.CommandLine.AddGoFlag(flag.CommandLine.Lookup("logtostderr"))
	pflag.CommandLine.Set("logtostderr", "true")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(readmeCmd)
	rootCmd.AddCommand(GetRecvCmd())
	rootCmd.AddCommand(GetSendCmd())

}

func Execute() error {
	return rootCmd.Execute()
}

func initializeConfig(cmd *cobra.Command) error {
	klog.V(6).Infof("initialize config")
	progName := filepath.Base(os.Args[0])
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/" + progName)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	} else {
		klog.V(6).Infof("config file loaded from one of default locations")
	}

	configFile, err := cmd.Flags().GetString("config")
	if err != nil {
		return err
	}

	if configFile != "" {
		configType := filepath.Ext(configFile)
		if configType[0] == '.' {
			configType = configType[1:]
		}
		klog.V(6).Infof("a config file given as parameter: %v type: %v", configFile, configType)
		if r, err := os.Open(configFile); err == nil {
			v.SetConfigType(configType)
			err = v.MergeConfig(r)
			if err != nil {
				klog.V(6).Error(err, "cannot merge config file")
				return err
			}
			r.Close()
			klog.V(6).Infof("config merged from %s", configFile)
		} else {
			klog.V(6).Error(err, "cannot open config file")
			return err
		}
	}

	v.SetEnvPrefix(strings.ToUpper(progName))
	v.AutomaticEnv()

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			v.BindEnv(f.Name, fmt.Sprintf("%s_%s", strings.ToUpper(progName), envVarSuffix))
		}

		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
	klog.V(6).Infof("config initialized")
	return nil
}

func main() {
	if err := Execute(); err != nil {
		klog.Errorf("zfsbackup command failed err=%v", err)
	}
}
