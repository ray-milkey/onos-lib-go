// Copyright 2020-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	"context"
	"time"

	api "github.com/onosproject/onos-lib-go/api/logging"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func getSetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set {level, debug} [args]",
		Short: "Sets a logger level or enable debug mode for logging package",
	}
	cmd.AddCommand(getSetLevelCommand())
	cmd.AddCommand(getSetDebugCommand())
	return cmd
}

func getSetDebugCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "debug [args]",
		Short: "enable/disable debug mode for logging package",
		Args:  cobra.ExactArgs(1),
		RunE:  runSetDebugCommand,
	}
	return cmd
}

func runSetDebugCommand(cmd *cobra.Command, args []string) error {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer func() {
		err = conn.Close()
	}()

	debugMode := args[0]
	client := api.NewLoggerClient(conn)

	if debugMode == "enable" {
		req := api.SetDebugModeRequest{
			Debug: true,
		}
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		_, err = client.SetDebug(ctx, &req)
		return err

	} else if debugMode == "disable" {
		req := api.SetDebugModeRequest{
			Debug: false,
		}
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		_, err = client.SetDebug(ctx, &req)
		return err

	}
	return nil
}

func getSetLevelCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "level [args]",
		Short: "Sets a level to a logger",
		Args:  cobra.ExactArgs(1),
		RunE:  runSetLevelCommand,
	}
	cmd.Flags().StringP("level", "l", "INFO", "the logger level")

	return cmd
}

func runSetLevelCommand(cmd *cobra.Command, args []string) error {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer func() {
		err = conn.Close()
	}()

	name := args[0]
	if name == "" {
		return errors.New("The logger name should be provided")
	}
	level, err := cmd.Flags().GetString("level")
	if err != nil {
		return err
	}
	var apiLevel api.Level
	switch level {
	case api.Level_INFO.String():
		apiLevel = api.Level_INFO
	case api.Level_DEBUG.String():
		apiLevel = api.Level_DEBUG
	case api.Level_ERROR.String():
		apiLevel = api.Level_ERROR
	case api.Level_PANIC.String():
		apiLevel = api.Level_PANIC
	case api.Level_DPANIC.String():
		apiLevel = api.Level_DPANIC
	case api.Level_FATAL.String():
		apiLevel = api.Level_FATAL
	case api.Level_WARN.String():
		apiLevel = api.Level_WARN

	}

	client := api.NewLoggerClient(conn)
	req := api.SetLevelRequest{
		LoggerName: name,
		Level:      apiLevel,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err = client.SetLevel(ctx, &req)
	return err
}
