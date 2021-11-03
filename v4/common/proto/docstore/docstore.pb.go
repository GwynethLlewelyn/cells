// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.18.1
// source: docstore.proto

package docstore

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type DocumentType int32

const (
	DocumentType_JSON   DocumentType = 0
	DocumentType_BINARY DocumentType = 1
)

// Enum value maps for DocumentType.
var (
	DocumentType_name = map[int32]string{
		0: "JSON",
		1: "BINARY",
	}
	DocumentType_value = map[string]int32{
		"JSON":   0,
		"BINARY": 1,
	}
)

func (x DocumentType) Enum() *DocumentType {
	p := new(DocumentType)
	*p = x
	return p
}

func (x DocumentType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (DocumentType) Descriptor() protoreflect.EnumDescriptor {
	return file_docstore_proto_enumTypes[0].Descriptor()
}

func (DocumentType) Type() protoreflect.EnumType {
	return &file_docstore_proto_enumTypes[0]
}

func (x DocumentType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use DocumentType.Descriptor instead.
func (DocumentType) EnumDescriptor() ([]byte, []int) {
	return file_docstore_proto_rawDescGZIP(), []int{0}
}

type Document struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID            string       `protobuf:"bytes,2,opt,name=ID,proto3" json:"ID,omitempty"`
	Type          DocumentType `protobuf:"varint,3,opt,name=Type,proto3,enum=docstore.DocumentType" json:"Type,omitempty"`
	Owner         string       `protobuf:"bytes,4,opt,name=Owner,proto3" json:"Owner,omitempty"`
	Data          string       `protobuf:"bytes,5,opt,name=Data,proto3" json:"Data,omitempty"`
	IndexableMeta string       `protobuf:"bytes,6,opt,name=IndexableMeta,proto3" json:"IndexableMeta,omitempty"`
}

func (x *Document) Reset() {
	*x = Document{}
	if protoimpl.UnsafeEnabled {
		mi := &file_docstore_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Document) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Document) ProtoMessage() {}

func (x *Document) ProtoReflect() protoreflect.Message {
	mi := &file_docstore_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Document.ProtoReflect.Descriptor instead.
func (*Document) Descriptor() ([]byte, []int) {
	return file_docstore_proto_rawDescGZIP(), []int{0}
}

func (x *Document) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *Document) GetType() DocumentType {
	if x != nil {
		return x.Type
	}
	return DocumentType_JSON
}

func (x *Document) GetOwner() string {
	if x != nil {
		return x.Owner
	}
	return ""
}

func (x *Document) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

func (x *Document) GetIndexableMeta() string {
	if x != nil {
		return x.IndexableMeta
	}
	return ""
}

type DocumentQuery struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID        string `protobuf:"bytes,2,opt,name=ID,proto3" json:"ID,omitempty"`
	Owner     string `protobuf:"bytes,3,opt,name=Owner,proto3" json:"Owner,omitempty"`
	MetaQuery string `protobuf:"bytes,4,opt,name=MetaQuery,proto3" json:"MetaQuery,omitempty"`
}

func (x *DocumentQuery) Reset() {
	*x = DocumentQuery{}
	if protoimpl.UnsafeEnabled {
		mi := &file_docstore_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DocumentQuery) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DocumentQuery) ProtoMessage() {}

func (x *DocumentQuery) ProtoReflect() protoreflect.Message {
	mi := &file_docstore_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DocumentQuery.ProtoReflect.Descriptor instead.
func (*DocumentQuery) Descriptor() ([]byte, []int) {
	return file_docstore_proto_rawDescGZIP(), []int{1}
}

func (x *DocumentQuery) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *DocumentQuery) GetOwner() string {
	if x != nil {
		return x.Owner
	}
	return ""
}

func (x *DocumentQuery) GetMetaQuery() string {
	if x != nil {
		return x.MetaQuery
	}
	return ""
}

type PutDocumentRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StoreID    string    `protobuf:"bytes,1,opt,name=StoreID,proto3" json:"StoreID,omitempty"`
	DocumentID string    `protobuf:"bytes,2,opt,name=DocumentID,proto3" json:"DocumentID,omitempty"`
	Document   *Document `protobuf:"bytes,3,opt,name=Document,proto3" json:"Document,omitempty"`
}

func (x *PutDocumentRequest) Reset() {
	*x = PutDocumentRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_docstore_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PutDocumentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PutDocumentRequest) ProtoMessage() {}

func (x *PutDocumentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_docstore_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PutDocumentRequest.ProtoReflect.Descriptor instead.
func (*PutDocumentRequest) Descriptor() ([]byte, []int) {
	return file_docstore_proto_rawDescGZIP(), []int{2}
}

func (x *PutDocumentRequest) GetStoreID() string {
	if x != nil {
		return x.StoreID
	}
	return ""
}

func (x *PutDocumentRequest) GetDocumentID() string {
	if x != nil {
		return x.DocumentID
	}
	return ""
}

func (x *PutDocumentRequest) GetDocument() *Document {
	if x != nil {
		return x.Document
	}
	return nil
}

type PutDocumentResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Document *Document `protobuf:"bytes,1,opt,name=Document,proto3" json:"Document,omitempty"`
}

func (x *PutDocumentResponse) Reset() {
	*x = PutDocumentResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_docstore_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PutDocumentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PutDocumentResponse) ProtoMessage() {}

func (x *PutDocumentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_docstore_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PutDocumentResponse.ProtoReflect.Descriptor instead.
func (*PutDocumentResponse) Descriptor() ([]byte, []int) {
	return file_docstore_proto_rawDescGZIP(), []int{3}
}

func (x *PutDocumentResponse) GetDocument() *Document {
	if x != nil {
		return x.Document
	}
	return nil
}

type GetDocumentRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StoreID    string `protobuf:"bytes,1,opt,name=StoreID,proto3" json:"StoreID,omitempty"`
	DocumentID string `protobuf:"bytes,2,opt,name=DocumentID,proto3" json:"DocumentID,omitempty"`
}

func (x *GetDocumentRequest) Reset() {
	*x = GetDocumentRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_docstore_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetDocumentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDocumentRequest) ProtoMessage() {}

func (x *GetDocumentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_docstore_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDocumentRequest.ProtoReflect.Descriptor instead.
func (*GetDocumentRequest) Descriptor() ([]byte, []int) {
	return file_docstore_proto_rawDescGZIP(), []int{4}
}

func (x *GetDocumentRequest) GetStoreID() string {
	if x != nil {
		return x.StoreID
	}
	return ""
}

func (x *GetDocumentRequest) GetDocumentID() string {
	if x != nil {
		return x.DocumentID
	}
	return ""
}

type GetDocumentResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Document  *Document `protobuf:"bytes,1,opt,name=Document,proto3" json:"Document,omitempty"`
	BinaryUrl string    `protobuf:"bytes,2,opt,name=BinaryUrl,proto3" json:"BinaryUrl,omitempty"`
}

func (x *GetDocumentResponse) Reset() {
	*x = GetDocumentResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_docstore_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetDocumentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDocumentResponse) ProtoMessage() {}

func (x *GetDocumentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_docstore_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDocumentResponse.ProtoReflect.Descriptor instead.
func (*GetDocumentResponse) Descriptor() ([]byte, []int) {
	return file_docstore_proto_rawDescGZIP(), []int{5}
}

func (x *GetDocumentResponse) GetDocument() *Document {
	if x != nil {
		return x.Document
	}
	return nil
}

func (x *GetDocumentResponse) GetBinaryUrl() string {
	if x != nil {
		return x.BinaryUrl
	}
	return ""
}

type DeleteDocumentsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StoreID    string         `protobuf:"bytes,1,opt,name=StoreID,proto3" json:"StoreID,omitempty"`
	DocumentID string         `protobuf:"bytes,2,opt,name=DocumentID,proto3" json:"DocumentID,omitempty"`
	Query      *DocumentQuery `protobuf:"bytes,3,opt,name=Query,proto3" json:"Query,omitempty"`
}

func (x *DeleteDocumentsRequest) Reset() {
	*x = DeleteDocumentsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_docstore_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteDocumentsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteDocumentsRequest) ProtoMessage() {}

func (x *DeleteDocumentsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_docstore_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteDocumentsRequest.ProtoReflect.Descriptor instead.
func (*DeleteDocumentsRequest) Descriptor() ([]byte, []int) {
	return file_docstore_proto_rawDescGZIP(), []int{6}
}

func (x *DeleteDocumentsRequest) GetStoreID() string {
	if x != nil {
		return x.StoreID
	}
	return ""
}

func (x *DeleteDocumentsRequest) GetDocumentID() string {
	if x != nil {
		return x.DocumentID
	}
	return ""
}

func (x *DeleteDocumentsRequest) GetQuery() *DocumentQuery {
	if x != nil {
		return x.Query
	}
	return nil
}

type DeleteDocumentsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success       bool  `protobuf:"varint,1,opt,name=Success,proto3" json:"Success,omitempty"`
	DeletionCount int32 `protobuf:"varint,2,opt,name=DeletionCount,proto3" json:"DeletionCount,omitempty"`
}

func (x *DeleteDocumentsResponse) Reset() {
	*x = DeleteDocumentsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_docstore_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteDocumentsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteDocumentsResponse) ProtoMessage() {}

func (x *DeleteDocumentsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_docstore_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteDocumentsResponse.ProtoReflect.Descriptor instead.
func (*DeleteDocumentsResponse) Descriptor() ([]byte, []int) {
	return file_docstore_proto_rawDescGZIP(), []int{7}
}

func (x *DeleteDocumentsResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *DeleteDocumentsResponse) GetDeletionCount() int32 {
	if x != nil {
		return x.DeletionCount
	}
	return 0
}

type ListDocumentsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StoreID string         `protobuf:"bytes,1,opt,name=StoreID,proto3" json:"StoreID,omitempty"`
	Query   *DocumentQuery `protobuf:"bytes,2,opt,name=Query,proto3" json:"Query,omitempty"`
}

func (x *ListDocumentsRequest) Reset() {
	*x = ListDocumentsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_docstore_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListDocumentsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListDocumentsRequest) ProtoMessage() {}

func (x *ListDocumentsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_docstore_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListDocumentsRequest.ProtoReflect.Descriptor instead.
func (*ListDocumentsRequest) Descriptor() ([]byte, []int) {
	return file_docstore_proto_rawDescGZIP(), []int{8}
}

func (x *ListDocumentsRequest) GetStoreID() string {
	if x != nil {
		return x.StoreID
	}
	return ""
}

func (x *ListDocumentsRequest) GetQuery() *DocumentQuery {
	if x != nil {
		return x.Query
	}
	return nil
}

type ListDocumentsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Document  *Document `protobuf:"bytes,1,opt,name=Document,proto3" json:"Document,omitempty"`
	BinaryUrl string    `protobuf:"bytes,2,opt,name=BinaryUrl,proto3" json:"BinaryUrl,omitempty"`
	Score     int32     `protobuf:"varint,3,opt,name=Score,proto3" json:"Score,omitempty"`
}

func (x *ListDocumentsResponse) Reset() {
	*x = ListDocumentsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_docstore_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListDocumentsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListDocumentsResponse) ProtoMessage() {}

func (x *ListDocumentsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_docstore_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListDocumentsResponse.ProtoReflect.Descriptor instead.
func (*ListDocumentsResponse) Descriptor() ([]byte, []int) {
	return file_docstore_proto_rawDescGZIP(), []int{9}
}

func (x *ListDocumentsResponse) GetDocument() *Document {
	if x != nil {
		return x.Document
	}
	return nil
}

func (x *ListDocumentsResponse) GetBinaryUrl() string {
	if x != nil {
		return x.BinaryUrl
	}
	return ""
}

func (x *ListDocumentsResponse) GetScore() int32 {
	if x != nil {
		return x.Score
	}
	return 0
}

type CountDocumentsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Total int64 `protobuf:"varint,1,opt,name=Total,proto3" json:"Total,omitempty"`
}

func (x *CountDocumentsResponse) Reset() {
	*x = CountDocumentsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_docstore_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CountDocumentsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CountDocumentsResponse) ProtoMessage() {}

func (x *CountDocumentsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_docstore_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CountDocumentsResponse.ProtoReflect.Descriptor instead.
func (*CountDocumentsResponse) Descriptor() ([]byte, []int) {
	return file_docstore_proto_rawDescGZIP(), []int{10}
}

func (x *CountDocumentsResponse) GetTotal() int64 {
	if x != nil {
		return x.Total
	}
	return 0
}

var File_docstore_proto protoreflect.FileDescriptor

var file_docstore_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x64, 0x6f, 0x63, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x08, 0x64, 0x6f, 0x63, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x22, 0x96, 0x01, 0x0a, 0x08, 0x44,
	0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x44, 0x12, 0x2a, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x16, 0x2e, 0x64, 0x6f, 0x63, 0x73, 0x74, 0x6f, 0x72, 0x65,
	0x2e, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x44, 0x61, 0x74,
	0x61, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x24, 0x0a,
	0x0d, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x61, 0x62, 0x6c, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x61, 0x62, 0x6c, 0x65, 0x4d,
	0x65, 0x74, 0x61, 0x22, 0x53, 0x0a, 0x0d, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x51,
	0x75, 0x65, 0x72, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x49, 0x44, 0x12, 0x14, 0x0a, 0x05, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x12, 0x1c, 0x0a, 0x09, 0x4d, 0x65,
	0x74, 0x61, 0x51, 0x75, 0x65, 0x72, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x4d,
	0x65, 0x74, 0x61, 0x51, 0x75, 0x65, 0x72, 0x79, 0x22, 0x7e, 0x0a, 0x12, 0x50, 0x75, 0x74, 0x44,
	0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18,
	0x0a, 0x07, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x49, 0x44, 0x12, 0x1e, 0x0a, 0x0a, 0x44, 0x6f, 0x63, 0x75,
	0x6d, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x44, 0x6f,
	0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x12, 0x2e, 0x0a, 0x08, 0x44, 0x6f, 0x63, 0x75,
	0x6d, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x64, 0x6f, 0x63,
	0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x08,
	0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x22, 0x45, 0x0a, 0x13, 0x50, 0x75, 0x74, 0x44,
	0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x2e, 0x0a, 0x08, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x12, 0x2e, 0x64, 0x6f, 0x63, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x44, 0x6f, 0x63,
	0x75, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x08, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x22,
	0x4e, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x49, 0x44,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x49, 0x44, 0x12,
	0x1e, 0x0a, 0x0a, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x22,
	0x63, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x08, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65,
	0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x64, 0x6f, 0x63, 0x73, 0x74,
	0x6f, 0x72, 0x65, 0x2e, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x08, 0x44, 0x6f,
	0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x42, 0x69, 0x6e, 0x61, 0x72, 0x79,
	0x55, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x42, 0x69, 0x6e, 0x61, 0x72,
	0x79, 0x55, 0x72, 0x6c, 0x22, 0x81, 0x01, 0x0a, 0x16, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x44,
	0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x18, 0x0a, 0x07, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x49, 0x44, 0x12, 0x1e, 0x0a, 0x0a, 0x44, 0x6f, 0x63,
	0x75, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x44,
	0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x12, 0x2d, 0x0a, 0x05, 0x51, 0x75, 0x65,
	0x72, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x64, 0x6f, 0x63, 0x73, 0x74,
	0x6f, 0x72, 0x65, 0x2e, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x51, 0x75, 0x65, 0x72,
	0x79, 0x52, 0x05, 0x51, 0x75, 0x65, 0x72, 0x79, 0x22, 0x59, 0x0a, 0x17, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x24, 0x0a,
	0x0d, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x0d, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f,
	0x75, 0x6e, 0x74, 0x22, 0x5f, 0x0a, 0x14, 0x4c, 0x69, 0x73, 0x74, 0x44, 0x6f, 0x63, 0x75, 0x6d,
	0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x53,
	0x74, 0x6f, 0x72, 0x65, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x53, 0x74,
	0x6f, 0x72, 0x65, 0x49, 0x44, 0x12, 0x2d, 0x0a, 0x05, 0x51, 0x75, 0x65, 0x72, 0x79, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x64, 0x6f, 0x63, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e,
	0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x05, 0x51,
	0x75, 0x65, 0x72, 0x79, 0x22, 0x7b, 0x0a, 0x15, 0x4c, 0x69, 0x73, 0x74, 0x44, 0x6f, 0x63, 0x75,
	0x6d, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a,
	0x08, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x12, 0x2e, 0x64, 0x6f, 0x63, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x44, 0x6f, 0x63, 0x75, 0x6d,
	0x65, 0x6e, 0x74, 0x52, 0x08, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x1c, 0x0a,
	0x09, 0x42, 0x69, 0x6e, 0x61, 0x72, 0x79, 0x55, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x42, 0x69, 0x6e, 0x61, 0x72, 0x79, 0x55, 0x72, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x53,
	0x63, 0x6f, 0x72, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x53, 0x63, 0x6f, 0x72,
	0x65, 0x22, 0x2e, 0x0a, 0x16, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65,
	0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x54,
	0x6f, 0x74, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x54, 0x6f, 0x74, 0x61,
	0x6c, 0x2a, 0x24, 0x0a, 0x0c, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x08, 0x0a, 0x04, 0x4a, 0x53, 0x4f, 0x4e, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x42,
	0x49, 0x4e, 0x41, 0x52, 0x59, 0x10, 0x01, 0x32, 0xac, 0x03, 0x0a, 0x08, 0x44, 0x6f, 0x63, 0x53,
	0x74, 0x6f, 0x72, 0x65, 0x12, 0x4c, 0x0a, 0x0b, 0x50, 0x75, 0x74, 0x44, 0x6f, 0x63, 0x75, 0x6d,
	0x65, 0x6e, 0x74, 0x12, 0x1c, 0x2e, 0x64, 0x6f, 0x63, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x50,
	0x75, 0x74, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1d, 0x2e, 0x64, 0x6f, 0x63, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x50, 0x75, 0x74,
	0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x4c, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e,
	0x74, 0x12, 0x1c, 0x2e, 0x64, 0x6f, 0x63, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x47, 0x65, 0x74,
	0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1d, 0x2e, 0x64, 0x6f, 0x63, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x44, 0x6f,
	0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x12, 0x58, 0x0a, 0x0f, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65,
	0x6e, 0x74, 0x73, 0x12, 0x20, 0x2e, 0x64, 0x6f, 0x63, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x64, 0x6f, 0x63, 0x73, 0x74, 0x6f, 0x72, 0x65,
	0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x54, 0x0a, 0x0e, 0x43, 0x6f,
	0x75, 0x6e, 0x74, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x1e, 0x2e, 0x64,
	0x6f, 0x63, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x44, 0x6f, 0x63, 0x75,
	0x6d, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x64,
	0x6f, 0x63, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x44, 0x6f, 0x63,
	0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x12, 0x54, 0x0a, 0x0d, 0x4c, 0x69, 0x73, 0x74, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74,
	0x73, 0x12, 0x1e, 0x2e, 0x64, 0x6f, 0x63, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x4c, 0x69, 0x73,
	0x74, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1f, 0x2e, 0x64, 0x6f, 0x63, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x4c, 0x69, 0x73,
	0x74, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x30, 0x01, 0x42, 0x31, 0x5a, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x79, 0x64, 0x69, 0x6f, 0x2f, 0x63, 0x65, 0x6c, 0x6c, 0x73,
	0x2f, 0x76, 0x34, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x64, 0x6f, 0x63, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_docstore_proto_rawDescOnce sync.Once
	file_docstore_proto_rawDescData = file_docstore_proto_rawDesc
)

func file_docstore_proto_rawDescGZIP() []byte {
	file_docstore_proto_rawDescOnce.Do(func() {
		file_docstore_proto_rawDescData = protoimpl.X.CompressGZIP(file_docstore_proto_rawDescData)
	})
	return file_docstore_proto_rawDescData
}

var file_docstore_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_docstore_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_docstore_proto_goTypes = []interface{}{
	(DocumentType)(0),               // 0: docstore.DocumentType
	(*Document)(nil),                // 1: docstore.Document
	(*DocumentQuery)(nil),           // 2: docstore.DocumentQuery
	(*PutDocumentRequest)(nil),      // 3: docstore.PutDocumentRequest
	(*PutDocumentResponse)(nil),     // 4: docstore.PutDocumentResponse
	(*GetDocumentRequest)(nil),      // 5: docstore.GetDocumentRequest
	(*GetDocumentResponse)(nil),     // 6: docstore.GetDocumentResponse
	(*DeleteDocumentsRequest)(nil),  // 7: docstore.DeleteDocumentsRequest
	(*DeleteDocumentsResponse)(nil), // 8: docstore.DeleteDocumentsResponse
	(*ListDocumentsRequest)(nil),    // 9: docstore.ListDocumentsRequest
	(*ListDocumentsResponse)(nil),   // 10: docstore.ListDocumentsResponse
	(*CountDocumentsResponse)(nil),  // 11: docstore.CountDocumentsResponse
}
var file_docstore_proto_depIdxs = []int32{
	0,  // 0: docstore.Document.Type:type_name -> docstore.DocumentType
	1,  // 1: docstore.PutDocumentRequest.Document:type_name -> docstore.Document
	1,  // 2: docstore.PutDocumentResponse.Document:type_name -> docstore.Document
	1,  // 3: docstore.GetDocumentResponse.Document:type_name -> docstore.Document
	2,  // 4: docstore.DeleteDocumentsRequest.Query:type_name -> docstore.DocumentQuery
	2,  // 5: docstore.ListDocumentsRequest.Query:type_name -> docstore.DocumentQuery
	1,  // 6: docstore.ListDocumentsResponse.Document:type_name -> docstore.Document
	3,  // 7: docstore.DocStore.PutDocument:input_type -> docstore.PutDocumentRequest
	5,  // 8: docstore.DocStore.GetDocument:input_type -> docstore.GetDocumentRequest
	7,  // 9: docstore.DocStore.DeleteDocuments:input_type -> docstore.DeleteDocumentsRequest
	9,  // 10: docstore.DocStore.CountDocuments:input_type -> docstore.ListDocumentsRequest
	9,  // 11: docstore.DocStore.ListDocuments:input_type -> docstore.ListDocumentsRequest
	4,  // 12: docstore.DocStore.PutDocument:output_type -> docstore.PutDocumentResponse
	6,  // 13: docstore.DocStore.GetDocument:output_type -> docstore.GetDocumentResponse
	8,  // 14: docstore.DocStore.DeleteDocuments:output_type -> docstore.DeleteDocumentsResponse
	11, // 15: docstore.DocStore.CountDocuments:output_type -> docstore.CountDocumentsResponse
	10, // 16: docstore.DocStore.ListDocuments:output_type -> docstore.ListDocumentsResponse
	12, // [12:17] is the sub-list for method output_type
	7,  // [7:12] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_docstore_proto_init() }
func file_docstore_proto_init() {
	if File_docstore_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_docstore_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Document); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_docstore_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DocumentQuery); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_docstore_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PutDocumentRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_docstore_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PutDocumentResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_docstore_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetDocumentRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_docstore_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetDocumentResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_docstore_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteDocumentsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_docstore_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteDocumentsResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_docstore_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListDocumentsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_docstore_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListDocumentsResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_docstore_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CountDocumentsResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_docstore_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_docstore_proto_goTypes,
		DependencyIndexes: file_docstore_proto_depIdxs,
		EnumInfos:         file_docstore_proto_enumTypes,
		MessageInfos:      file_docstore_proto_msgTypes,
	}.Build()
	File_docstore_proto = out.File
	file_docstore_proto_rawDesc = nil
	file_docstore_proto_goTypes = nil
	file_docstore_proto_depIdxs = nil
}
