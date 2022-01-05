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
	"fmt"
	"log"
	"sort"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/pydio/cells/v4/common/config"
)

var listConfigCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configurations",
	Long: `
DESCRIPTION

  Display all configuration items registered by the application.
  Configuration items are listed as truple [serviceName, configName, configValue]. The configuration value is json encoded.

`,
	Run: func(cmd *cobra.Command, args []string) {

		var m map[string]map[string]interface{}
		if err := config.Get("services").Scan(&m); err != nil {
			log.Fatal(err)
		}

		table := tablewriter.NewWriter(cmd.OutOrStdout())
		table.SetHeader([]string{"Service", "Configuration", "Value"})

		var skeys []string
		for k := range m {
			skeys = append(skeys, k)
		}

		sort.Strings(skeys)

		for _, sk := range skeys {
			var ckeys []string
			for k := range m[sk] {
				ckeys = append(ckeys, k)
			}
			sort.Strings(ckeys)
			for _, ck := range ckeys {
				table.Append([]string{sk, ck, fmt.Sprintf("%v", m[sk][ck])})
			}
		}

		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.Render()
	},
}

func init() {
	ConfigCmd.AddCommand(listConfigCmd)
}
