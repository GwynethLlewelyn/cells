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

// Package views provides high-level clients for talking to the main data tree in certain context.
//
// It follows the "wrapper" pattern of http handlers to filter all requests inputs and outputs. The "Router" is object
// is used by all services or gateways when accessing to data as a given user.
// Between others, it will
// - Load ACLs and perform checks to make sure user is allowed to read/write the data
// - Perform other meta-related or acl-related checks like Quota management, locks, etc..
// - Perform encryption/decryption of actual data on the fly
// - Compress / Decompress archives,
// - Add metadata collected from any services on the nodes outputted by the responses,
// - etc...
package nodes

import (
	"context"
	"io"

	"github.com/golang/protobuf/proto"
	"github.com/pydio/minio-go"
	"go.uber.org/zap/zapcore"

	"github.com/pydio/cells/common/nodes/models"
	"github.com/pydio/cells/common/proto/idm"
	"github.com/pydio/cells/common/proto/object"
	"github.com/pydio/cells/common/proto/tree"
)

const (
	VIEWS_LIBRARY_NAME = "pydio.lib.views"
)

var (
	// IsUnitTestEnv flag prevents among others the ClientPool to look for declared
	// datasources in the registry. As none is present while running unit tests, it
	// otherwise times out.
	IsUnitTestEnv = false
)

// These keys may be enriched in Context depending on the middleware
type (
	LoadedSource struct {
		object.DataSource
		Client *minio.Core
	}

	SourcesPool interface {
		Close()
		GetTreeClient() tree.NodeProviderClient
		GetTreeClientWrite() tree.NodeReceiverClient
		GetDataSourceInfo(dsName string, retries ...int) (LoadedSource, error)
		GetDataSources() map[string]LoadedSource
		LoadDataSources()
	}

	BranchInfo struct {
		LoadedSource
		idm.Workspace
		Root              *tree.Node
		Binary            bool
		TransparentBinary bool
		AncestorsList     map[string][]*tree.Node
	}
)

type NodeFilter func(ctx context.Context, inputNode *tree.Node, identifier string) (context.Context, *tree.Node, error)
type NodesCallback func(inputFilter NodeFilter, outputFilter NodeFilter) error

type WalkFunc func(ctx context.Context, node *tree.Node, err error) error
type WalkFilter func(ctx context.Context, node *tree.Node) bool

type Handler interface {
	tree.NodeProviderClient
	tree.NodeReceiverClient
	tree.NodeChangesStreamerClient
	GetObject(ctx context.Context, node *tree.Node, requestData *models.GetRequestData) (io.ReadCloser, error)
	PutObject(ctx context.Context, node *tree.Node, reader io.Reader, requestData *models.PutRequestData) (int64, error)
	CopyObject(ctx context.Context, from *tree.Node, to *tree.Node, requestData *models.CopyRequestData) (int64, error)

	MultipartCreate(ctx context.Context, target *tree.Node, requestData *models.MultipartRequestData) (string, error)
	MultipartPutObjectPart(ctx context.Context, target *tree.Node, uploadID string, partNumberMarker int, reader io.Reader, requestData *models.PutRequestData) (models.MultipartObjectPart, error)
	MultipartAbort(ctx context.Context, target *tree.Node, uploadID string, requestData *models.MultipartRequestData) error
	MultipartComplete(ctx context.Context, target *tree.Node, uploadID string, uploadedParts []models.MultipartObjectPart) (models.S3ObjectInfo, error)

	MultipartList(ctx context.Context, prefix string, requestData *models.MultipartRequestData) (models.ListMultipartUploadsResult, error)
	MultipartListObjectParts(ctx context.Context, target *tree.Node, uploadID string, partNumberMarker int, maxParts int) (models.ListObjectPartsResult, error)

	ExecuteWrapped(inputFilter NodeFilter, outputFilter NodeFilter, provider NodesCallback) error
	WrappedCanApply(srcCtx context.Context, targetCtx context.Context, operation *tree.NodeChangeEvent) error
	ListNodesWithCallback(ctx context.Context, request *tree.ListNodesRequest, callback WalkFunc, ignoreCbError bool, filters ...WalkFilter) error

	SetClientsPool(p SourcesPool)
}

// IsInternal check if either datasource is internal or branch has Binary flag
func (b BranchInfo) IsInternal() bool {
	return b.Binary || b.LoadedSource.IsInternal()
}

// WithBucketName creates a copy of a LoadedSource with a bucket name
func WithBucketName(s LoadedSource, bucket string) LoadedSource {
	out := LoadedSource{
		Client: s.Client,
	}
	c := proto.Clone(&s.DataSource).(*object.DataSource)
	c.ObjectsBucket = bucket
	out.DataSource = *c
	return out
}

func (s LoadedSource) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	return encoder.AddObject("DataSource", &s.DataSource)
}
