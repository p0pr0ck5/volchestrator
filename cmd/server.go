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
	"github.com/spf13/cobra"

	"github.com/p0pr0ck5/volchestrator/config"
	"github.com/p0pr0ck5/volchestrator/server/wrapper"
)

var configPath string

func run(cmd *cobra.Command, args []string) {
	var config config.ServerConfig
	err := hclsimple.DecodeFile(configPath, nil, &config)
	if err != nil {
		log.Fatalf("failed to decode config: %s", err)
	}

	w, err := wrapper.NewWrapper(config)
	if err != nil {
		log.Fatalf("failed to create new wrapper: %s", err)
	}

	w.Start()

	sigs := make(chan os.Signal, 1)
	done := make(chan struct{})

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs

	go func() {
		w.Stop()
		close(done)
	}()

	t := time.After(time.Second)

	select {
	case <-done:
	case <-t:
		log.Println("timeout")
	}
}

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the volchestrator server",
	Long:  `TBD`,
	Run:   run,
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVarP(&configPath, "config-path", "c", "config/examples/server.hcl", "Path for the config file")
}
