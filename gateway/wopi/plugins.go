/*
 * Copyright (c) 2018. Abstrium SAS <team (at) pydio.com>
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

// Package wopi serves files using the WOPI protocol.
package wopi

import (
	"context"
	"net/http"

	"github.com/pydio/cells/common"
	"github.com/pydio/cells/common/nodes"
	"github.com/pydio/cells/common/plugins"
	"github.com/pydio/cells/common/service"
)

var (
	viewsRouter *nodes.Router
)

func init() {
	plugins.Register("main", func(ctx context.Context) {
		service.NewService(
			service.Name(common.ServiceGatewayWopi),
			service.Context(ctx),
			service.Tag(common.ServiceTagGateway),
			service.RouterDependencies(),
			service.Description("WOPI REST Gateway to tree service"),
			service.WithHTTP(func() http.Handler {
				viewsRouter = nodes.NewUuidRouter(nodes.RouterOptions{WatchRegistry: true, AuditEvent: true})

				return NewRouter()
			}),
		)
	})
}
