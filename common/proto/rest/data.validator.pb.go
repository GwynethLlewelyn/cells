// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: data.proto

package rest

import github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/pydio/cells/common/proto/docstore"
import _ "github.com/pydio/cells/common/proto/tree"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (this *SearchResults) Validate() error {
	for _, item := range this.Results {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Results", err)
			}
		}
	}
	return nil
}
func (this *Pagination) Validate() error {
	return nil
}
func (this *Metadata) Validate() error {
	return nil
}
func (this *MetaCollection) Validate() error {
	for _, item := range this.Metadatas {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Metadatas", err)
			}
		}
	}
	return nil
}
func (this *MetaNamespaceRequest) Validate() error {
	return nil
}
func (this *GetBulkMetaRequest) Validate() error {
	return nil
}
func (this *BulkMetaResponse) Validate() error {
	for _, item := range this.Nodes {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Nodes", err)
			}
		}
	}
	if this.Pagination != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Pagination); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Pagination", err)
		}
	}
	return nil
}
func (this *HeadNodeRequest) Validate() error {
	return nil
}
func (this *HeadNodeResponse) Validate() error {
	if this.Node != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Node); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Node", err)
		}
	}
	return nil
}
func (this *CreateNodesRequest) Validate() error {
	for _, item := range this.Nodes {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Nodes", err)
			}
		}
	}
	return nil
}
func (this *CreateSelectionRequest) Validate() error {
	for _, item := range this.Nodes {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Nodes", err)
			}
		}
	}
	return nil
}
func (this *CreateSelectionResponse) Validate() error {
	for _, item := range this.Nodes {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Nodes", err)
			}
		}
	}
	return nil
}
func (this *NodesCollection) Validate() error {
	if this.Parent != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Parent); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Parent", err)
		}
	}
	for _, item := range this.Children {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Children", err)
			}
		}
	}
	return nil
}
func (this *DeleteNodesRequest) Validate() error {
	for _, item := range this.Nodes {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Nodes", err)
			}
		}
	}
	return nil
}
func (this *BackgroundJobResult) Validate() error {
	return nil
}
func (this *DeleteNodesResponse) Validate() error {
	for _, item := range this.DeleteJobs {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("DeleteJobs", err)
			}
		}
	}
	return nil
}
func (this *RestoreNodesRequest) Validate() error {
	for _, item := range this.Nodes {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Nodes", err)
			}
		}
	}
	return nil
}
func (this *RestoreNodesResponse) Validate() error {
	for _, item := range this.RestoreJobs {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("RestoreJobs", err)
			}
		}
	}
	return nil
}
func (this *ListDocstoreRequest) Validate() error {
	if this.Query != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Query); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Query", err)
		}
	}
	return nil
}
func (this *DocstoreCollection) Validate() error {
	for _, item := range this.Docs {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Docs", err)
			}
		}
	}
	return nil
}
