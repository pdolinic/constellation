// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.20.1
// source: activation.proto

package activationproto

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

type ActivateWorkerNodeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DiskUuid string `protobuf:"bytes,1,opt,name=disk_uuid,json=diskUuid,proto3" json:"disk_uuid,omitempty"`
	NodeName string `protobuf:"bytes,2,opt,name=node_name,json=nodeName,proto3" json:"node_name,omitempty"`
}

func (x *ActivateWorkerNodeRequest) Reset() {
	*x = ActivateWorkerNodeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_activation_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ActivateWorkerNodeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ActivateWorkerNodeRequest) ProtoMessage() {}

func (x *ActivateWorkerNodeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_activation_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ActivateWorkerNodeRequest.ProtoReflect.Descriptor instead.
func (*ActivateWorkerNodeRequest) Descriptor() ([]byte, []int) {
	return file_activation_proto_rawDescGZIP(), []int{0}
}

func (x *ActivateWorkerNodeRequest) GetDiskUuid() string {
	if x != nil {
		return x.DiskUuid
	}
	return ""
}

func (x *ActivateWorkerNodeRequest) GetNodeName() string {
	if x != nil {
		return x.NodeName
	}
	return ""
}

type ActivateWorkerNodeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StateDiskKey             []byte `protobuf:"bytes,1,opt,name=state_disk_key,json=stateDiskKey,proto3" json:"state_disk_key,omitempty"`
	OwnerId                  []byte `protobuf:"bytes,2,opt,name=owner_id,json=ownerId,proto3" json:"owner_id,omitempty"`
	ClusterId                []byte `protobuf:"bytes,3,opt,name=cluster_id,json=clusterId,proto3" json:"cluster_id,omitempty"`
	KubeletKey               []byte `protobuf:"bytes,4,opt,name=kubelet_key,json=kubeletKey,proto3" json:"kubelet_key,omitempty"`
	KubeletCert              []byte `protobuf:"bytes,5,opt,name=kubelet_cert,json=kubeletCert,proto3" json:"kubelet_cert,omitempty"`
	ApiServerEndpoint        string `protobuf:"bytes,6,opt,name=api_server_endpoint,json=apiServerEndpoint,proto3" json:"api_server_endpoint,omitempty"`
	Token                    string `protobuf:"bytes,7,opt,name=token,proto3" json:"token,omitempty"`
	DiscoveryTokenCaCertHash string `protobuf:"bytes,8,opt,name=discovery_token_ca_cert_hash,json=discoveryTokenCaCertHash,proto3" json:"discovery_token_ca_cert_hash,omitempty"`
}

func (x *ActivateWorkerNodeResponse) Reset() {
	*x = ActivateWorkerNodeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_activation_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ActivateWorkerNodeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ActivateWorkerNodeResponse) ProtoMessage() {}

func (x *ActivateWorkerNodeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_activation_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ActivateWorkerNodeResponse.ProtoReflect.Descriptor instead.
func (*ActivateWorkerNodeResponse) Descriptor() ([]byte, []int) {
	return file_activation_proto_rawDescGZIP(), []int{1}
}

func (x *ActivateWorkerNodeResponse) GetStateDiskKey() []byte {
	if x != nil {
		return x.StateDiskKey
	}
	return nil
}

func (x *ActivateWorkerNodeResponse) GetOwnerId() []byte {
	if x != nil {
		return x.OwnerId
	}
	return nil
}

func (x *ActivateWorkerNodeResponse) GetClusterId() []byte {
	if x != nil {
		return x.ClusterId
	}
	return nil
}

func (x *ActivateWorkerNodeResponse) GetKubeletKey() []byte {
	if x != nil {
		return x.KubeletKey
	}
	return nil
}

func (x *ActivateWorkerNodeResponse) GetKubeletCert() []byte {
	if x != nil {
		return x.KubeletCert
	}
	return nil
}

func (x *ActivateWorkerNodeResponse) GetApiServerEndpoint() string {
	if x != nil {
		return x.ApiServerEndpoint
	}
	return ""
}

func (x *ActivateWorkerNodeResponse) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *ActivateWorkerNodeResponse) GetDiscoveryTokenCaCertHash() string {
	if x != nil {
		return x.DiscoveryTokenCaCertHash
	}
	return ""
}

type ActivateControlPlaneNodeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DiskUuid string `protobuf:"bytes,1,opt,name=disk_uuid,json=diskUuid,proto3" json:"disk_uuid,omitempty"`
	NodeName string `protobuf:"bytes,2,opt,name=node_name,json=nodeName,proto3" json:"node_name,omitempty"`
}

func (x *ActivateControlPlaneNodeRequest) Reset() {
	*x = ActivateControlPlaneNodeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_activation_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ActivateControlPlaneNodeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ActivateControlPlaneNodeRequest) ProtoMessage() {}

func (x *ActivateControlPlaneNodeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_activation_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ActivateControlPlaneNodeRequest.ProtoReflect.Descriptor instead.
func (*ActivateControlPlaneNodeRequest) Descriptor() ([]byte, []int) {
	return file_activation_proto_rawDescGZIP(), []int{2}
}

func (x *ActivateControlPlaneNodeRequest) GetDiskUuid() string {
	if x != nil {
		return x.DiskUuid
	}
	return ""
}

func (x *ActivateControlPlaneNodeRequest) GetNodeName() string {
	if x != nil {
		return x.NodeName
	}
	return ""
}

type ActivateControlPlaneNodeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StateDiskKey             []byte `protobuf:"bytes,1,opt,name=state_disk_key,json=stateDiskKey,proto3" json:"state_disk_key,omitempty"`
	OwnerId                  []byte `protobuf:"bytes,2,opt,name=owner_id,json=ownerId,proto3" json:"owner_id,omitempty"`
	ClusterId                []byte `protobuf:"bytes,3,opt,name=cluster_id,json=clusterId,proto3" json:"cluster_id,omitempty"`
	KubeletKey               []byte `protobuf:"bytes,4,opt,name=kubelet_key,json=kubeletKey,proto3" json:"kubelet_key,omitempty"`
	KubeletCert              []byte `protobuf:"bytes,5,opt,name=kubelet_cert,json=kubeletCert,proto3" json:"kubelet_cert,omitempty"`
	ApiServerEndpoint        string `protobuf:"bytes,6,opt,name=api_server_endpoint,json=apiServerEndpoint,proto3" json:"api_server_endpoint,omitempty"`
	Token                    string `protobuf:"bytes,7,opt,name=token,proto3" json:"token,omitempty"`
	DiscoveryTokenCaCertHash string `protobuf:"bytes,8,opt,name=discovery_token_ca_cert_hash,json=discoveryTokenCaCertHash,proto3" json:"discovery_token_ca_cert_hash,omitempty"`
	CertificateKey           string `protobuf:"bytes,9,opt,name=certificate_key,json=certificateKey,proto3" json:"certificate_key,omitempty"`
}

func (x *ActivateControlPlaneNodeResponse) Reset() {
	*x = ActivateControlPlaneNodeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_activation_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ActivateControlPlaneNodeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ActivateControlPlaneNodeResponse) ProtoMessage() {}

func (x *ActivateControlPlaneNodeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_activation_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ActivateControlPlaneNodeResponse.ProtoReflect.Descriptor instead.
func (*ActivateControlPlaneNodeResponse) Descriptor() ([]byte, []int) {
	return file_activation_proto_rawDescGZIP(), []int{3}
}

func (x *ActivateControlPlaneNodeResponse) GetStateDiskKey() []byte {
	if x != nil {
		return x.StateDiskKey
	}
	return nil
}

func (x *ActivateControlPlaneNodeResponse) GetOwnerId() []byte {
	if x != nil {
		return x.OwnerId
	}
	return nil
}

func (x *ActivateControlPlaneNodeResponse) GetClusterId() []byte {
	if x != nil {
		return x.ClusterId
	}
	return nil
}

func (x *ActivateControlPlaneNodeResponse) GetKubeletKey() []byte {
	if x != nil {
		return x.KubeletKey
	}
	return nil
}

func (x *ActivateControlPlaneNodeResponse) GetKubeletCert() []byte {
	if x != nil {
		return x.KubeletCert
	}
	return nil
}

func (x *ActivateControlPlaneNodeResponse) GetApiServerEndpoint() string {
	if x != nil {
		return x.ApiServerEndpoint
	}
	return ""
}

func (x *ActivateControlPlaneNodeResponse) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *ActivateControlPlaneNodeResponse) GetDiscoveryTokenCaCertHash() string {
	if x != nil {
		return x.DiscoveryTokenCaCertHash
	}
	return ""
}

func (x *ActivateControlPlaneNodeResponse) GetCertificateKey() string {
	if x != nil {
		return x.CertificateKey
	}
	return ""
}

var File_activation_proto protoreflect.FileDescriptor

var file_activation_proto_rawDesc = []byte{
	0x0a, 0x10, 0x61, 0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0a, 0x61, 0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x55,
	0x0a, 0x19, 0x41, 0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x65, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72,
	0x4e, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x64,
	0x69, 0x73, 0x6b, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x64, 0x69, 0x73, 0x6b, 0x55, 0x75, 0x69, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x6e, 0x6f, 0x64, 0x65,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6e, 0x6f, 0x64,
	0x65, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0xc6, 0x02, 0x0a, 0x1a, 0x41, 0x63, 0x74, 0x69, 0x76, 0x61,
	0x74, 0x65, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x24, 0x0a, 0x0e, 0x73, 0x74, 0x61, 0x74, 0x65, 0x5f, 0x64, 0x69,
	0x73, 0x6b, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x73, 0x74,
	0x61, 0x74, 0x65, 0x44, 0x69, 0x73, 0x6b, 0x4b, 0x65, 0x79, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x77,
	0x6e, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x6f, 0x77,
	0x6e, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72,
	0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x63, 0x6c, 0x75, 0x73, 0x74,
	0x65, 0x72, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x6b, 0x75, 0x62, 0x65, 0x6c, 0x65, 0x74, 0x5f,
	0x6b, 0x65, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x6b, 0x75, 0x62, 0x65, 0x6c,
	0x65, 0x74, 0x4b, 0x65, 0x79, 0x12, 0x21, 0x0a, 0x0c, 0x6b, 0x75, 0x62, 0x65, 0x6c, 0x65, 0x74,
	0x5f, 0x63, 0x65, 0x72, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0b, 0x6b, 0x75, 0x62,
	0x65, 0x6c, 0x65, 0x74, 0x43, 0x65, 0x72, 0x74, 0x12, 0x2e, 0x0a, 0x13, 0x61, 0x70, 0x69, 0x5f,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x61, 0x70, 0x69, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x3e,
	0x0a, 0x1c, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x5f, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x5f, 0x63, 0x61, 0x5f, 0x63, 0x65, 0x72, 0x74, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x08,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x18, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x43, 0x61, 0x43, 0x65, 0x72, 0x74, 0x48, 0x61, 0x73, 0x68, 0x22, 0x5b,
	0x0a, 0x1f, 0x41, 0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f,
	0x6c, 0x50, 0x6c, 0x61, 0x6e, 0x65, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1b, 0x0a, 0x09, 0x64, 0x69, 0x73, 0x6b, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x64, 0x69, 0x73, 0x6b, 0x55, 0x75, 0x69, 0x64, 0x12, 0x1b,
	0x0a, 0x09, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x6e, 0x6f, 0x64, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0xf5, 0x02, 0x0a, 0x20,
	0x41, 0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x50,
	0x6c, 0x61, 0x6e, 0x65, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x24, 0x0a, 0x0e, 0x73, 0x74, 0x61, 0x74, 0x65, 0x5f, 0x64, 0x69, 0x73, 0x6b, 0x5f, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x73, 0x74, 0x61, 0x74, 0x65, 0x44,
	0x69, 0x73, 0x6b, 0x4b, 0x65, 0x79, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x49,
	0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x1f, 0x0a, 0x0b, 0x6b, 0x75, 0x62, 0x65, 0x6c, 0x65, 0x74, 0x5f, 0x6b, 0x65, 0x79, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x6b, 0x75, 0x62, 0x65, 0x6c, 0x65, 0x74, 0x4b, 0x65,
	0x79, 0x12, 0x21, 0x0a, 0x0c, 0x6b, 0x75, 0x62, 0x65, 0x6c, 0x65, 0x74, 0x5f, 0x63, 0x65, 0x72,
	0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0b, 0x6b, 0x75, 0x62, 0x65, 0x6c, 0x65, 0x74,
	0x43, 0x65, 0x72, 0x74, 0x12, 0x2e, 0x0a, 0x13, 0x61, 0x70, 0x69, 0x5f, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x5f, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x11, 0x61, 0x70, 0x69, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x45, 0x6e, 0x64, 0x70,
	0x6f, 0x69, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x3e, 0x0a, 0x1c, 0x64, 0x69,
	0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x63, 0x61,
	0x5f, 0x63, 0x65, 0x72, 0x74, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x18, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x43, 0x61, 0x43, 0x65, 0x72, 0x74, 0x48, 0x61, 0x73, 0x68, 0x12, 0x27, 0x0a, 0x0f, 0x63, 0x65,
	0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x09, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0e, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65,
	0x4b, 0x65, 0x79, 0x32, 0xe1, 0x01, 0x0a, 0x03, 0x41, 0x50, 0x49, 0x12, 0x63, 0x0a, 0x12, 0x41,
	0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x65, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x4e, 0x6f, 0x64,
	0x65, 0x12, 0x25, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x41,
	0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x65, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x4e, 0x6f, 0x64,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x26, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x76,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x65, 0x57, 0x6f,
	0x72, 0x6b, 0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x75, 0x0a, 0x18, 0x41, 0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x74,
	0x72, 0x6f, 0x6c, 0x50, 0x6c, 0x61, 0x6e, 0x65, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x2b, 0x2e, 0x61,
	0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x76, 0x61,
	0x74, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x50, 0x6c, 0x61, 0x6e, 0x65, 0x4e, 0x6f,
	0x64, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2c, 0x2e, 0x61, 0x63, 0x74, 0x69,
	0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x65, 0x43,
	0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x50, 0x6c, 0x61, 0x6e, 0x65, 0x4e, 0x6f, 0x64, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x48, 0x5a, 0x46, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x65, 0x64, 0x67, 0x65, 0x6c, 0x65, 0x73, 0x73, 0x73, 0x79,
	0x73, 0x2f, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x65, 0x6c, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f,
	0x61, 0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_activation_proto_rawDescOnce sync.Once
	file_activation_proto_rawDescData = file_activation_proto_rawDesc
)

func file_activation_proto_rawDescGZIP() []byte {
	file_activation_proto_rawDescOnce.Do(func() {
		file_activation_proto_rawDescData = protoimpl.X.CompressGZIP(file_activation_proto_rawDescData)
	})
	return file_activation_proto_rawDescData
}

var file_activation_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_activation_proto_goTypes = []interface{}{
	(*ActivateWorkerNodeRequest)(nil),        // 0: activation.ActivateWorkerNodeRequest
	(*ActivateWorkerNodeResponse)(nil),       // 1: activation.ActivateWorkerNodeResponse
	(*ActivateControlPlaneNodeRequest)(nil),  // 2: activation.ActivateControlPlaneNodeRequest
	(*ActivateControlPlaneNodeResponse)(nil), // 3: activation.ActivateControlPlaneNodeResponse
}
var file_activation_proto_depIdxs = []int32{
	0, // 0: activation.API.ActivateWorkerNode:input_type -> activation.ActivateWorkerNodeRequest
	2, // 1: activation.API.ActivateControlPlaneNode:input_type -> activation.ActivateControlPlaneNodeRequest
	1, // 2: activation.API.ActivateWorkerNode:output_type -> activation.ActivateWorkerNodeResponse
	3, // 3: activation.API.ActivateControlPlaneNode:output_type -> activation.ActivateControlPlaneNodeResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_activation_proto_init() }
func file_activation_proto_init() {
	if File_activation_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_activation_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ActivateWorkerNodeRequest); i {
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
		file_activation_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ActivateWorkerNodeResponse); i {
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
		file_activation_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ActivateControlPlaneNodeRequest); i {
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
		file_activation_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ActivateControlPlaneNodeResponse); i {
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
			RawDescriptor: file_activation_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_activation_proto_goTypes,
		DependencyIndexes: file_activation_proto_depIdxs,
		MessageInfos:      file_activation_proto_msgTypes,
	}.Build()
	File_activation_proto = out.File
	file_activation_proto_rawDesc = nil
	file_activation_proto_goTypes = nil
	file_activation_proto_depIdxs = nil
}