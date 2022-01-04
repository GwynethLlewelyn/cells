/*
 * Copyright (c) 2019-2021. Abstrium SAS <team (at) pydio.com>
 * This file is part of Pydio Cells.
 *
 * Pydio Cells is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Pydio Cells is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Pydio Cells.  If not, see <http://www.gnu.org/licenses/>.
 *
 * The latest code can be found at <https://pydio.com>.
 */

package cmd

import (
	"github.com/spf13/cobra"
)

// AdminCmd groups the data manipulation commands
// The sub-commands are connecting via gRPC to a **running** Cells instance.
var AdminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Direct Read/Write access to Cells data",
	Long: `
DESCRIPTION

  Set of commands with direct access to Cells data.
	
  These commands require a running Cells instance. They connect directly to low-level services
  using gRPC connections. They are not authenticated.
`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		bindViperFlags(cmd.Flags(), map[string]string{})

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	// Registry / Broker Flags
	addRegistryFlags(AdminCmd.PersistentFlags())
	RootCmd.AddCommand(AdminCmd)
}
