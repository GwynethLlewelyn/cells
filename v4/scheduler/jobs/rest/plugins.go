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

// Package rest provides a REST gateway to the job definition repository.
package rest

import (
	"context"

	"github.com/pydio/cells/v4/common/nodes"
	"github.com/pydio/cells/v4/common/nodes/compose"
	servicecontext "github.com/pydio/cells/v4/common/service/context"

	"github.com/pydio/cells/v4/common"
	"github.com/pydio/cells/v4/common/plugins"
	"github.com/pydio/cells/v4/common/service"
)

func init() {
	plugins.Register("main", func(ctx context.Context) {
		service.NewService(
			service.Name(common.ServiceRestNamespace_+common.ServiceJobs),
			service.Context(ctx),
			service.Tag(common.ServiceTagScheduler),
			service.Description("REST gateway to the scheduler service"),
			service.Dependency(common.ServiceGrpcNamespace_+common.ServiceJobs, []string{}),
			service.Dependency(common.ServiceGrpcNamespace_+common.ServiceTasks, []string{}),
			service.WithWeb(func(ctx context.Context) service.WebHandler {
				// Init router with current registry
				router = compose.PathClient(
					nodes.WithContext(ctx),
					nodes.WithRegistryWatch(servicecontext.GetRegistry(ctx)),
				)
				return new(JobsHandler)
			}),
		)
	})
}
