// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: cellsapi-share.proto

package rest

import (
	fmt "fmt"
	math "math"
	proto "google.golang.org/protobuf/proto"
	_ "github.com/pydio/cells/v4/common/proto/idm"
	_ "github.com/pydio/cells/v4/common/proto/tree"
	_ "github.com/mwitkow/go-proto-validators"
	_ "github.com/pydio/cells/v4/common/proto/service"
	github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (this *CellAcl) Validate() error {
	for _, item := range this.Actions {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Actions", err)
			}
		}
	}
	if this.User != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.User); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("User", err)
		}
	}
	if this.Group != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Group); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Group", err)
		}
	}
	if this.Role != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Role); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Role", err)
		}
	}
	return nil
}
func (this *Cell) Validate() error {
	if !(len(this.Label) < 500) {
		return github_com_mwitkow_go_proto_validators.FieldError("Label", fmt.Errorf(`value '%v' must have a length smaller than '500'`, this.Label))
	}
	if !(len(this.Description) < 1000) {
		return github_com_mwitkow_go_proto_validators.FieldError("Description", fmt.Errorf(`value '%v' must have a length smaller than '1000'`, this.Description))
	}
	for _, item := range this.RootNodes {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("RootNodes", err)
			}
		}
	}
	// Validation of proto3 map<> fields is unsupported.
	for _, item := range this.Policies {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Policies", err)
			}
		}
	}
	return nil
}
func (this *ShareLinkTargetUser) Validate() error {
	return nil
}
func (this *ShareLink) Validate() error {
	if !(len(this.Label) < 500) {
		return github_com_mwitkow_go_proto_validators.FieldError("Label", fmt.Errorf(`value '%v' must have a length smaller than '500'`, this.Label))
	}
	if !(len(this.Description) < 1000) {
		return github_com_mwitkow_go_proto_validators.FieldError("Description", fmt.Errorf(`value '%v' must have a length smaller than '1000'`, this.Description))
	}
	// Validation of proto3 map<> fields is unsupported.
	for _, item := range this.RootNodes {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("RootNodes", err)
			}
		}
	}
	for _, item := range this.Policies {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Policies", err)
			}
		}
	}
	return nil
}
func (this *PutCellRequest) Validate() error {
	if this.Room != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Room); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Room", err)
		}
	}
	return nil
}
func (this *GetCellRequest) Validate() error {
	return nil
}
func (this *DeleteCellRequest) Validate() error {
	return nil
}
func (this *DeleteCellResponse) Validate() error {
	return nil
}
func (this *GetShareLinkRequest) Validate() error {
	return nil
}
func (this *PutShareLinkRequest) Validate() error {
	if this.ShareLink != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.ShareLink); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("ShareLink", err)
		}
	}
	return nil
}
func (this *DeleteShareLinkRequest) Validate() error {
	return nil
}
func (this *DeleteShareLinkResponse) Validate() error {
	return nil
}
func (this *ListSharedResourcesRequest) Validate() error {
	return nil
}
func (this *ListSharedResourcesResponse) Validate() error {
	for _, item := range this.Resources {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Resources", err)
			}
		}
	}
	return nil
}
func (this *ListSharedResourcesResponse_SharedResource) Validate() error {
	if this.Node != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Node); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Node", err)
		}
	}
	if this.Link != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Link); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Link", err)
		}
	}
	for _, item := range this.Cells {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Cells", err)
			}
		}
	}
	return nil
}
func (this *UpdateSharePoliciesRequest) Validate() error {
	for _, item := range this.Policies {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Policies", err)
			}
		}
	}
	return nil
}
func (this *UpdateSharePoliciesResponse) Validate() error {
	for _, item := range this.Policies {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Policies", err)
			}
		}
	}
	return nil
}
