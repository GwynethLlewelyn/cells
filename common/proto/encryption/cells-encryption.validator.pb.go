// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: cells-encryption.proto

package encryption

import (
	fmt "fmt"
	math "math"
	proto "google.golang.org/protobuf/proto"
	github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (this *Export) Validate() error {
	return nil
}
func (this *Import) Validate() error {
	return nil
}
func (this *KeyInfo) Validate() error {
	for _, item := range this.Exports {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Exports", err)
			}
		}
	}
	for _, item := range this.Imports {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Imports", err)
			}
		}
	}
	return nil
}
func (this *Key) Validate() error {
	if this.Info != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Info); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Info", err)
		}
	}
	return nil
}
func (this *AddKeyRequest) Validate() error {
	if this.Key != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Key); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Key", err)
		}
	}
	return nil
}
func (this *AddKeyResponse) Validate() error {
	return nil
}
func (this *GetKeyRequest) Validate() error {
	return nil
}
func (this *GetKeyResponse) Validate() error {
	if this.Key != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Key); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Key", err)
		}
	}
	return nil
}
func (this *AdminListKeysRequest) Validate() error {
	return nil
}
func (this *AdminListKeysResponse) Validate() error {
	for _, item := range this.Keys {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Keys", err)
			}
		}
	}
	return nil
}
func (this *AdminDeleteKeyRequest) Validate() error {
	return nil
}
func (this *AdminDeleteKeyResponse) Validate() error {
	return nil
}
func (this *AdminExportKeyRequest) Validate() error {
	return nil
}
func (this *AdminExportKeyResponse) Validate() error {
	if this.Key != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Key); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Key", err)
		}
	}
	return nil
}
func (this *AdminImportKeyRequest) Validate() error {
	if this.Key != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Key); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Key", err)
		}
	}
	return nil
}
func (this *AdminImportKeyResponse) Validate() error {
	return nil
}
func (this *AdminCreateKeyRequest) Validate() error {
	return nil
}
func (this *AdminCreateKeyResponse) Validate() error {
	return nil
}
func (this *NodeKey) Validate() error {
	return nil
}
func (this *Node) Validate() error {
	return nil
}
func (this *NodeInfo) Validate() error {
	if this.Node != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Node); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Node", err)
		}
	}
	if this.NodeKey != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.NodeKey); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("NodeKey", err)
		}
	}
	if this.Block != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Block); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Block", err)
		}
	}
	return nil
}
func (this *Block) Validate() error {
	return nil
}
func (this *RangedBlock) Validate() error {
	return nil
}
func (this *GetNodeInfoRequest) Validate() error {
	return nil
}
func (this *GetNodeInfoResponse) Validate() error {
	if this.NodeInfo != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.NodeInfo); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("NodeInfo", err)
		}
	}
	return nil
}
func (this *GetNodePlainSizeRequest) Validate() error {
	return nil
}
func (this *GetNodePlainSizeResponse) Validate() error {
	return nil
}
func (this *SetNodeInfoRequest) Validate() error {
	if this.SetNodeKey != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.SetNodeKey); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("SetNodeKey", err)
		}
	}
	if this.SetBlock != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.SetBlock); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("SetBlock", err)
		}
	}
	return nil
}
func (this *SetNodeInfoResponse) Validate() error {
	if this.SetNodeKey != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.SetNodeKey); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("SetNodeKey", err)
		}
	}
	if this.SetBlock != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.SetBlock); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("SetBlock", err)
		}
	}
	return nil
}
func (this *DeleteNodeRequest) Validate() error {
	return nil
}
func (this *DeleteNodeResponse) Validate() error {
	return nil
}
func (this *DeleteNodeKeyRequest) Validate() error {
	return nil
}
func (this *DeleteNodeKeyResponse) Validate() error {
	return nil
}
func (this *DeleteNodeSharedKeyRequest) Validate() error {
	return nil
}
func (this *DeleteNodeSharedKeyResponse) Validate() error {
	return nil
}
func (this *SetNodeKeyRequest) Validate() error {
	if this.NodeKey != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.NodeKey); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("NodeKey", err)
		}
	}
	return nil
}
func (this *SetNodeKeyResponse) Validate() error {
	return nil
}
func (this *SetNodeBlockRequest) Validate() error {
	if this.Block != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Block); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Block", err)
		}
	}
	return nil
}
func (this *SetNodeBlockResponse) Validate() error {
	return nil
}
func (this *CopyNodeInfoRequest) Validate() error {
	return nil
}
func (this *CopyNodeInfoResponse) Validate() error {
	return nil
}
