// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: sf/gateway/payment/v1/gateway.proto

package pbgateway

import (
	v1 "github.com/streamingfast/dmetering/pb/sf/metering/v1"
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

type ServicesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ServicesRequest) Reset() {
	*x = ServicesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sf_gateway_payment_v1_gateway_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServicesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServicesRequest) ProtoMessage() {}

func (x *ServicesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_sf_gateway_payment_v1_gateway_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServicesRequest.ProtoReflect.Descriptor instead.
func (*ServicesRequest) Descriptor() ([]byte, []int) {
	return file_sf_gateway_payment_v1_gateway_proto_rawDescGZIP(), []int{0}
}

type ServicesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UsageEndpoint      string `protobuf:"bytes,1,opt,name=usage_endpoint,json=usageEndpoint,proto3" json:"usage_endpoint,omitempty"`                // points to the "Usage" service endpoint
	AuthEndpoint       string `protobuf:"bytes,2,opt,name=auth_endpoint,json=authEndpoint,proto3" json:"auth_endpoint,omitempty"`                   // .. for "Auth"
	RevocationEndpoint string `protobuf:"bytes,3,opt,name=revocation_endpoint,json=revocationEndpoint,proto3" json:"revocation_endpoint,omitempty"` // .. for "Revocation"
}

func (x *ServicesResponse) Reset() {
	*x = ServicesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sf_gateway_payment_v1_gateway_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServicesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServicesResponse) ProtoMessage() {}

func (x *ServicesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_sf_gateway_payment_v1_gateway_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServicesResponse.ProtoReflect.Descriptor instead.
func (*ServicesResponse) Descriptor() ([]byte, []int) {
	return file_sf_gateway_payment_v1_gateway_proto_rawDescGZIP(), []int{1}
}

func (x *ServicesResponse) GetUsageEndpoint() string {
	if x != nil {
		return x.UsageEndpoint
	}
	return ""
}

func (x *ServicesResponse) GetAuthEndpoint() string {
	if x != nil {
		return x.AuthEndpoint
	}
	return ""
}

func (x *ServicesResponse) GetRevocationEndpoint() string {
	if x != nil {
		return x.RevocationEndpoint
	}
	return ""
}

type ReportRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Events []*v1.Event `protobuf:"bytes,1,rep,name=events,proto3" json:"events,omitempty"`
}

func (x *ReportRequest) Reset() {
	*x = ReportRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sf_gateway_payment_v1_gateway_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReportRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReportRequest) ProtoMessage() {}

func (x *ReportRequest) ProtoReflect() protoreflect.Message {
	mi := &file_sf_gateway_payment_v1_gateway_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReportRequest.ProtoReflect.Descriptor instead.
func (*ReportRequest) Descriptor() ([]byte, []int) {
	return file_sf_gateway_payment_v1_gateway_proto_rawDescGZIP(), []int{2}
}

func (x *ReportRequest) GetEvents() []*v1.Event {
	if x != nil {
		return x.Events
	}
	return nil
}

type ReportResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Revoked          bool   `protobuf:"varint,1,opt,name=revoked,proto3" json:"revoked,omitempty"`
	RevocationReason string `protobuf:"bytes,2,opt,name=revocation_reason,json=revocationReason,proto3" json:"revocation_reason,omitempty"`
}

func (x *ReportResponse) Reset() {
	*x = ReportResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sf_gateway_payment_v1_gateway_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReportResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReportResponse) ProtoMessage() {}

func (x *ReportResponse) ProtoReflect() protoreflect.Message {
	mi := &file_sf_gateway_payment_v1_gateway_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReportResponse.ProtoReflect.Descriptor instead.
func (*ReportResponse) Descriptor() ([]byte, []int) {
	return file_sf_gateway_payment_v1_gateway_proto_rawDescGZIP(), []int{3}
}

func (x *ReportResponse) GetRevoked() bool {
	if x != nil {
		return x.Revoked
	}
	return false
}

func (x *ReportResponse) GetRevocationReason() string {
	if x != nil {
		return x.RevocationReason
	}
	return ""
}

type AuthenticateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ApiKey          string  `protobuf:"bytes,1,opt,name=api_key,json=apiKey,proto3" json:"api_key,omitempty"`
	Origin          *string `protobuf:"bytes,2,opt,name=origin,proto3,oneof" json:"origin,omitempty"`
	IpAddr          *string `protobuf:"bytes,3,opt,name=ip_addr,json=ipAddr,proto3,oneof" json:"ip_addr,omitempty"`
	LifetimeSeconds uint64  `protobuf:"varint,4,opt,name=lifetime_seconds,json=lifetimeSeconds,proto3" json:"lifetime_seconds,omitempty"` // usually long-lived
}

func (x *AuthenticateRequest) Reset() {
	*x = AuthenticateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sf_gateway_payment_v1_gateway_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthenticateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthenticateRequest) ProtoMessage() {}

func (x *AuthenticateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_sf_gateway_payment_v1_gateway_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthenticateRequest.ProtoReflect.Descriptor instead.
func (*AuthenticateRequest) Descriptor() ([]byte, []int) {
	return file_sf_gateway_payment_v1_gateway_proto_rawDescGZIP(), []int{4}
}

func (x *AuthenticateRequest) GetApiKey() string {
	if x != nil {
		return x.ApiKey
	}
	return ""
}

func (x *AuthenticateRequest) GetOrigin() string {
	if x != nil && x.Origin != nil {
		return *x.Origin
	}
	return ""
}

func (x *AuthenticateRequest) GetIpAddr() string {
	if x != nil && x.IpAddr != nil {
		return *x.IpAddr
	}
	return ""
}

func (x *AuthenticateRequest) GetLifetimeSeconds() uint64 {
	if x != nil {
		return x.LifetimeSeconds
	}
	return 0
}

type AuthenticateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success    bool   `protobuf:"varint,2,opt,name=success,proto3" json:"success,omitempty"`
	Jwt        string `protobuf:"bytes,1,opt,name=jwt,proto3" json:"jwt,omitempty"`
	Expiration uint64 `protobuf:"varint,3,opt,name=expiration,proto3" json:"expiration,omitempty"`
	UserId     string `protobuf:"bytes,4,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	JwtId      string `protobuf:"bytes,5,opt,name=jwt_id,json=jwtId,proto3" json:"jwt_id,omitempty"` //bytes key_authority = 6; // Who signed this? What's the public key? For us to validate whether we will honor this partner. Say this is E&N's, are we configured to honor those keys?
}

func (x *AuthenticateResponse) Reset() {
	*x = AuthenticateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sf_gateway_payment_v1_gateway_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthenticateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthenticateResponse) ProtoMessage() {}

func (x *AuthenticateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_sf_gateway_payment_v1_gateway_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthenticateResponse.ProtoReflect.Descriptor instead.
func (*AuthenticateResponse) Descriptor() ([]byte, []int) {
	return file_sf_gateway_payment_v1_gateway_proto_rawDescGZIP(), []int{5}
}

func (x *AuthenticateResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *AuthenticateResponse) GetJwt() string {
	if x != nil {
		return x.Jwt
	}
	return ""
}

func (x *AuthenticateResponse) GetExpiration() uint64 {
	if x != nil {
		return x.Expiration
	}
	return 0
}

func (x *AuthenticateResponse) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *AuthenticateResponse) GetJwtId() string {
	if x != nil {
		return x.JwtId
	}
	return ""
}

var File_sf_gateway_payment_v1_gateway_proto protoreflect.FileDescriptor

var file_sf_gateway_payment_v1_gateway_proto_rawDesc = []byte{
	0x0a, 0x23, 0x73, 0x66, 0x2f, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x70, 0x61, 0x79,
	0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x15, 0x73, 0x66, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61,
	0x79, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x1a, 0x1d, 0x73, 0x66,
	0x2f, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x65, 0x74,
	0x65, 0x72, 0x69, 0x6e, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x11, 0x0a, 0x0f, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x8f,
	0x01, 0x0a, 0x10, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x75, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x65, 0x6e, 0x64,
	0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x75, 0x73, 0x61,
	0x67, 0x65, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x61, 0x75,
	0x74, 0x68, 0x5f, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0c, 0x61, 0x75, 0x74, 0x68, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12,
	0x2f, 0x0a, 0x13, 0x72, 0x65, 0x76, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x65, 0x6e,
	0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x72, 0x65,
	0x76, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74,
	0x22, 0x3e, 0x0a, 0x0d, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x2d, 0x0a, 0x06, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x15, 0x2e, 0x73, 0x66, 0x2e, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x2e,
	0x76, 0x31, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x06, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73,
	0x22, 0x57, 0x0a, 0x0e, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x72, 0x65, 0x76, 0x6f, 0x6b, 0x65, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x07, 0x72, 0x65, 0x76, 0x6f, 0x6b, 0x65, 0x64, 0x12, 0x2b, 0x0a, 0x11,
	0x72, 0x65, 0x76, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x72, 0x65, 0x61, 0x73, 0x6f,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x72, 0x65, 0x76, 0x6f, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x22, 0xab, 0x01, 0x0a, 0x13, 0x41, 0x75,
	0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x17, 0x0a, 0x07, 0x61, 0x70, 0x69, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x61, 0x70, 0x69, 0x4b, 0x65, 0x79, 0x12, 0x1b, 0x0a, 0x06, 0x6f, 0x72,
	0x69, 0x67, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x06, 0x6f, 0x72,
	0x69, 0x67, 0x69, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x1c, 0x0a, 0x07, 0x69, 0x70, 0x5f, 0x61, 0x64,
	0x64, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x06, 0x69, 0x70, 0x41, 0x64,
	0x64, 0x72, 0x88, 0x01, 0x01, 0x12, 0x29, 0x0a, 0x10, 0x6c, 0x69, 0x66, 0x65, 0x74, 0x69, 0x6d,
	0x65, 0x5f, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x0f, 0x6c, 0x69, 0x66, 0x65, 0x74, 0x69, 0x6d, 0x65, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73,
	0x42, 0x09, 0x0a, 0x07, 0x5f, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x42, 0x0a, 0x0a, 0x08, 0x5f,
	0x69, 0x70, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x22, 0x92, 0x01, 0x0a, 0x14, 0x41, 0x75, 0x74, 0x68,
	0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x6a, 0x77,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6a, 0x77, 0x74, 0x12, 0x1e, 0x0a, 0x0a,
	0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x17, 0x0a, 0x07,
	0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x15, 0x0a, 0x06, 0x6a, 0x77, 0x74, 0x5f, 0x69, 0x64, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6a, 0x77, 0x74, 0x49, 0x64, 0x32, 0x6e, 0x0a, 0x0f,
	0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x5b, 0x0a, 0x08, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x12, 0x26, 0x2e, 0x73, 0x66,
	0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74,
	0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x27, 0x2e, 0x73, 0x66, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79,
	0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0x65, 0x0a, 0x0c,
	0x55, 0x73, 0x61, 0x67, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x55, 0x0a, 0x06,
	0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x24, 0x2e, 0x73, 0x66, 0x2e, 0x67, 0x61, 0x74, 0x65,
	0x77, 0x61, 0x79, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x52,
	0x65, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x73,
	0x66, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e,
	0x74, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x32, 0x76, 0x0a, 0x0b, 0x41, 0x75, 0x74, 0x68, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x67, 0x0a, 0x0c, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61,
	0x74, 0x65, 0x12, 0x2a, 0x2e, 0x73, 0x66, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e,
	0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x65,
	0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2b,
	0x2e, 0x73, 0x66, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x70, 0x61, 0x79, 0x6d,
	0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63,
	0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x4d, 0x5a, 0x4b, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x69, 0x6e, 0x67, 0x66, 0x61, 0x73, 0x74, 0x2f, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2d,
	0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x70, 0x62, 0x2f, 0x73, 0x66, 0x2f, 0x67, 0x61,
	0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x76, 0x31,
	0x3b, 0x70, 0x62, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_sf_gateway_payment_v1_gateway_proto_rawDescOnce sync.Once
	file_sf_gateway_payment_v1_gateway_proto_rawDescData = file_sf_gateway_payment_v1_gateway_proto_rawDesc
)

func file_sf_gateway_payment_v1_gateway_proto_rawDescGZIP() []byte {
	file_sf_gateway_payment_v1_gateway_proto_rawDescOnce.Do(func() {
		file_sf_gateway_payment_v1_gateway_proto_rawDescData = protoimpl.X.CompressGZIP(file_sf_gateway_payment_v1_gateway_proto_rawDescData)
	})
	return file_sf_gateway_payment_v1_gateway_proto_rawDescData
}

var file_sf_gateway_payment_v1_gateway_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_sf_gateway_payment_v1_gateway_proto_goTypes = []interface{}{
	(*ServicesRequest)(nil),      // 0: sf.gateway.payment.v1.ServicesRequest
	(*ServicesResponse)(nil),     // 1: sf.gateway.payment.v1.ServicesResponse
	(*ReportRequest)(nil),        // 2: sf.gateway.payment.v1.ReportRequest
	(*ReportResponse)(nil),       // 3: sf.gateway.payment.v1.ReportResponse
	(*AuthenticateRequest)(nil),  // 4: sf.gateway.payment.v1.AuthenticateRequest
	(*AuthenticateResponse)(nil), // 5: sf.gateway.payment.v1.AuthenticateResponse
	(*v1.Event)(nil),             // 6: sf.metering.v1.Event
}
var file_sf_gateway_payment_v1_gateway_proto_depIdxs = []int32{
	6, // 0: sf.gateway.payment.v1.ReportRequest.events:type_name -> sf.metering.v1.Event
	0, // 1: sf.gateway.payment.v1.DiscoverService.Services:input_type -> sf.gateway.payment.v1.ServicesRequest
	2, // 2: sf.gateway.payment.v1.UsageService.Report:input_type -> sf.gateway.payment.v1.ReportRequest
	4, // 3: sf.gateway.payment.v1.AuthService.Authenticate:input_type -> sf.gateway.payment.v1.AuthenticateRequest
	1, // 4: sf.gateway.payment.v1.DiscoverService.Services:output_type -> sf.gateway.payment.v1.ServicesResponse
	3, // 5: sf.gateway.payment.v1.UsageService.Report:output_type -> sf.gateway.payment.v1.ReportResponse
	5, // 6: sf.gateway.payment.v1.AuthService.Authenticate:output_type -> sf.gateway.payment.v1.AuthenticateResponse
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_sf_gateway_payment_v1_gateway_proto_init() }
func file_sf_gateway_payment_v1_gateway_proto_init() {
	if File_sf_gateway_payment_v1_gateway_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_sf_gateway_payment_v1_gateway_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServicesRequest); i {
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
		file_sf_gateway_payment_v1_gateway_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServicesResponse); i {
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
		file_sf_gateway_payment_v1_gateway_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReportRequest); i {
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
		file_sf_gateway_payment_v1_gateway_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReportResponse); i {
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
		file_sf_gateway_payment_v1_gateway_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthenticateRequest); i {
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
		file_sf_gateway_payment_v1_gateway_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthenticateResponse); i {
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
	file_sf_gateway_payment_v1_gateway_proto_msgTypes[4].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_sf_gateway_payment_v1_gateway_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   3,
		},
		GoTypes:           file_sf_gateway_payment_v1_gateway_proto_goTypes,
		DependencyIndexes: file_sf_gateway_payment_v1_gateway_proto_depIdxs,
		MessageInfos:      file_sf_gateway_payment_v1_gateway_proto_msgTypes,
	}.Build()
	File_sf_gateway_payment_v1_gateway_proto = out.File
	file_sf_gateway_payment_v1_gateway_proto_rawDesc = nil
	file_sf_gateway_payment_v1_gateway_proto_goTypes = nil
	file_sf_gateway_payment_v1_gateway_proto_depIdxs = nil
}
