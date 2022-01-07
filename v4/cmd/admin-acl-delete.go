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
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/spf13/cobra"

	"github.com/pydio/cells/v4/common"
	"github.com/pydio/cells/v4/common/client/grpc"
	"github.com/pydio/cells/v4/common/proto/idm"
	"github.com/pydio/cells/v4/common/proto/service"
)

var deleteAclCmd = &cobra.Command{
	Use:   "delete",
	Short: "Remove one or more ACLs",
	Long: `
DESCRIPTION
  
  Remove one or more ACLs by querying the ACL API.
  Flags allow you to query the gRPC service and delete the resulting ACLs.

`,
	Run: func(cmd *cobra.Command, args []string) {
		client := idm.NewACLServiceClient(grpc.GetClientConnFromCtx(ctx, common.ServiceAcl))

		var aclActions []*idm.ACLAction
		for _, action := range actions {
			aclActions = append(aclActions, &idm.ACLAction{
				Name:  action,
				Value: "1",
			})
		}

		query, _ := ptypes.MarshalAny(&idm.ACLSingleQuery{
			Actions:      aclActions,
			RoleIDs:      roleIDs,
			WorkspaceIDs: workspaceIDs,
			NodeIDs:      nodeIDs,
		})

		if len(query.Value) == 0 {
			return
		}

		if response, err := client.DeleteACL(context.Background(), &idm.DeleteACLRequest{
			Query: &service.Query{
				SubQueries: []*any.Any{query},
			},
		}); err != nil {
			fmt.Println(err)
			return
		} else {
			fmt.Printf("Successfully deleted ACL - [%v]\n", response.RowsDeleted)
			return
		}
	},
}

func init() {
	deleteAclCmd.Flags().StringArrayVarP(&actions, "action", "a", []string{}, "Action value")
	deleteAclCmd.Flags().StringArrayVarP(&roleIDs, "role_id", "r", []string{}, "RoleIDs")
	deleteAclCmd.Flags().StringArrayVarP(&workspaceIDs, "workspace_id", "w", []string{}, "WorkspaceIDs")
	deleteAclCmd.Flags().StringArrayVarP(&nodeIDs, "node_id", "n", []string{}, "NodeIDs")

	AclCmd.AddCommand(deleteAclCmd)
}
