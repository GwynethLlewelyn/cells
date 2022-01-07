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
	"context"
	"os"
	"path"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/manifoldco/promptui"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/pydio/cells/v4/common"
	"github.com/pydio/cells/v4/common/client/grpc"
	"github.com/pydio/cells/v4/common/proto/tree"
)

var (
	lsPath       string
	lsRecursive  bool
	lsShowHidden bool
	lsShowUuid   bool
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List files",
	Long: `
DESCRIPTION
  
  List Nodes by querying the tree microservice. Paths are computed starting from the root, their first segment is always
  a datasource name.

EXAMPLE

  List all files at the root of the "Common Files" workspace

  $ ` + os.Args[0] + ` admin files ls --path pydiods1 --uuid
	+--------+---------------+--------------------------------------+--------+-----------------+
	|  TYPE  |     PATH      |                 UUID                 |  SIZE  |    MODIFIED     |
	+--------+---------------+--------------------------------------+--------+-----------------+
	| Folder | Shared Folder | dcabadd0-7d32-4e45-9d1f-a927d3d0c174 | 147 MB | 21 Oct 21 08:58 |
	+--------+---------------+--------------------------------------+--------+-----------------+

 `,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := tree.NewNodeProviderClient(grpc.GetClientConnFromCtx(ctx, common.ServiceTree))

		// List all children and move them all
		streamer, err := client.ListNodes(context.Background(), &tree.ListNodesRequest{Node: &tree.Node{Path: lsPath}, Recursive: lsRecursive})
		if err != nil {
			return err
		}

		cmd.Println("")
		cmd.Println("Listing nodes under " + promptui.Styler(promptui.FGUnderline)(lsPath))
		table := tablewriter.NewWriter(cmd.OutOrStdout())
		hh := []string{"Type", "Path", "Size", "Modified"}
		if lsShowUuid {
			hh = []string{"Type", "Path", "Uuid", "Size", "Modified"}
		}
		table.SetHeader(hh)
		res := 0
		for {
			resp, err := streamer.Recv()
			if err != nil {
				break
			}
			res++
			node := resp.GetNode()
			if path.Base(node.GetPath()) == common.PydioSyncHiddenFile && !lsShowHidden {
				continue
			}
			var t, p, s, m string
			p = strings.TrimLeft(strings.TrimPrefix(node.GetPath(), lsPath), "/")
			t = "Folder"
			s = humanize.Bytes(uint64(node.GetSize()))
			if node.GetSize() == 0 {
				s = "-"
			}
			m = time.Unix(node.GetMTime(), 0).Format("02 Jan 06 15:04")
			if node.GetMTime() == 0 {
				m = "-"
			}
			if node.IsLeaf() {
				t = "File"
			}
			if lsShowUuid {
				table.Append([]string{t, p, node.GetUuid(), s, m})
			} else {
				table.Append([]string{t, p, s, m})
			}
		}
		if res > 0 {
			table.Render()
		} else {
			cmd.Println("No results")
		}
		return nil
	},
}

func init() {
	lsCmd.Flags().StringVarP(&lsPath, "path", "p", "/", "List nodes under given path")
	lsCmd.Flags().BoolVarP(&lsRecursive, "recursive", "", false, "List nodes recursively")
	lsCmd.Flags().BoolVarP(&lsShowUuid, "uuid", "", false, "Show UUIDs")
	lsCmd.Flags().BoolVarP(&lsShowHidden, "hidden", "", false, "Show hidden files (.pydio)")
	FilesCmd.AddCommand(lsCmd)
}
