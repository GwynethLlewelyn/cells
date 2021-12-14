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

// Package rest provides access to the TokenServiceApi
package rest

import (
	"context"

	"github.com/pydio/cells/v4/common"
	"github.com/pydio/cells/v4/common/plugins"
	"github.com/pydio/cells/v4/common/service"
)

func init() {
	plugins.Register("main", func(ctx context.Context) {
		service.NewService(
			service.Name(common.ServiceRestNamespace_+common.ServiceAuth),
			service.Context(ctx),
			service.Tag(common.ServiceTagIdm),
			service.Description("RESTful Gateway to token service"),
			service.Dependency(common.ServiceGrpcNamespace_+common.ServiceUser, []string{}),
			service.Dependency(common.ServiceGrpcNamespace_+common.ServiceDocStore, []string{}),
			service.Dependency(common.ServiceGrpcNamespace_+common.ServiceMailer, []string{}),
			service.WithWeb(func(c context.Context) service.WebHandler {
				return new(TokenHandler)
			}),
		)
	})
}
