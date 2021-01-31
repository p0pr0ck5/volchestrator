package cmd

/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/p0pr0ck5/volchestrator/client"
	"github.com/p0pr0ck5/volchestrator/config"
	"github.com/spf13/cobra"
)

var clientConfigPath string

func clientRun(cmd *cobra.Command, args []string) {
	var config config.ClientConfig
	err := hclsimple.DecodeFile(clientConfigPath, nil, &config)
	if err != nil {
		log.Fatalf("failed to decode config: %s", err)
	}

	c, err := client.NewClient(config)
	if err != nil {
		log.Fatalf("failed to create a new client: %s", err)
	}

	c.Run()

	sigs := make(chan os.Signal, 1)
	done := make(chan struct{})

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs

	go func() {
		c.Stop()
		close(done)
	}()

	t := time.After(time.Second)

	select {
	case <-done:
	case <-t:
		log.Println("timeout")
	}
}

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Dummy volchestrator client",
	Long:  `TBD`,
	Run:   clientRun,
}

func init() {
	rootCmd.AddCommand(clientCmd)

	clientCmd.Flags().StringVarP(&clientConfigPath, "config-path", "c", "config/examples/client.hcl", "Path for the config file")
}
