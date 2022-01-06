// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.18.0
// source: echo.proto

// Generate with protoc --go_out=. echo.proto && protoc --go-grpc_out=:. ./echo.proto

package proto

import (
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
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

type EchoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *EchoRequest) Reset() {
	*x = EchoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_echo_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EchoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EchoRequest) ProtoMessage() {}

func (x *EchoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_echo_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EchoRequest.ProtoReflect.Descriptor instead.
func (*EchoRequest) Descriptor() ([]byte, []int) {
	return file_echo_proto_rawDescGZIP(), []int{0}
}

func (x *EchoRequest) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type EchoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *EchoResponse) Reset() {
	*x = EchoResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_echo_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EchoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EchoResponse) ProtoMessage() {}

func (x *EchoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_echo_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EchoResponse.ProtoReflect.Descriptor instead.
func (*EchoResponse) Descriptor() ([]byte, []int) {
	return file_echo_proto_rawDescGZIP(), []int{1}
}

func (x *EchoResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type Header struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *Header) Reset() {
	*x = Header{}
	if protoimpl.UnsafeEnabled {
		mi := &file_echo_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Header) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Header) ProtoMessage() {}

func (x *Header) ProtoReflect() protoreflect.Message {
	mi := &file_echo_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Header.ProtoReflect.Descriptor instead.
func (*Header) Descriptor() ([]byte, []int) {
	return file_echo_proto_rawDescGZIP(), []int{2}
}

func (x *Header) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *Header) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type ForwardEchoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Count         int32     `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
	Qps           int32     `protobuf:"varint,2,opt,name=qps,proto3" json:"qps,omitempty"`
	TimeoutMicros int64     `protobuf:"varint,3,opt,name=timeout_micros,json=timeoutMicros,proto3" json:"timeout_micros,omitempty"`
	Url           string    `protobuf:"bytes,4,opt,name=url,proto3" json:"url,omitempty"`
	Headers       []*Header `protobuf:"bytes,5,rep,name=headers,proto3" json:"headers,omitempty"`
	Message       string    `protobuf:"bytes,6,opt,name=message,proto3" json:"message,omitempty"`
	// Method for the request. Valid only for HTTP
	Method string `protobuf:"bytes,9,opt,name=method,proto3" json:"method,omitempty"`
	// If true, requests will be sent using h2c prior knowledge
	Http2 bool `protobuf:"varint,7,opt,name=http2,proto3" json:"http2,omitempty"`
	// If true, requests will be sent using http3
	Http3 bool `protobuf:"varint,15,opt,name=http3,proto3" json:"http3,omitempty"`
	// If true, requests will not be sent until magic string is received
	ServerFirst bool `protobuf:"varint,8,opt,name=serverFirst,proto3" json:"serverFirst,omitempty"`
	// If true, 301 redirects will be followed
	FollowRedirects bool `protobuf:"varint,14,opt,name=followRedirects,proto3" json:"followRedirects,omitempty"`
	// If non-empty, make the request with the corresponding cert and key.
	Cert string `protobuf:"bytes,10,opt,name=cert,proto3" json:"cert,omitempty"`
	Key  string `protobuf:"bytes,11,opt,name=key,proto3" json:"key,omitempty"`
	// If non-empty, verify the server CA
	CaCert string `protobuf:"bytes,12,opt,name=caCert,proto3" json:"caCert,omitempty"`
	// If non-empty, make the request with the corresponding cert and key file.
	CertFile string `protobuf:"bytes,16,opt,name=certFile,proto3" json:"certFile,omitempty"`
	KeyFile  string `protobuf:"bytes,17,opt,name=keyFile,proto3" json:"keyFile,omitempty"`
	// If non-empty, verify the server CA with the ca cert file.
	CaCertFile string `protobuf:"bytes,18,opt,name=caCertFile,proto3" json:"caCertFile,omitempty"`
	// Skip verifying peer's certificate.
	InsecureSkipVerify bool `protobuf:"varint,19,opt,name=insecureSkipVerify,proto3" json:"insecureSkipVerify,omitempty"`
	// List of ALPNs to present. If not set, this will be automatically be set based on the protocol
	Alpn *Alpn `protobuf:"bytes,13,opt,name=alpn,proto3" json:"alpn,omitempty"`
	// Server name (SNI) to present in TLS connections. If not set, Host will be used for http requests.
	ServerName string `protobuf:"bytes,20,opt,name=serverName,proto3" json:"serverName,omitempty"`
	// Expected response determines what string to look for in the response to validate TCP requests succeeded.
	// If not set, defaults to "StatusCode=200"
	ExpectedResponse *wrappers.StringValue `protobuf:"bytes,21,opt,name=expectedResponse,proto3" json:"expectedResponse,omitempty"`
}

func (x *ForwardEchoRequest) Reset() {
	*x = ForwardEchoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_echo_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ForwardEchoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ForwardEchoRequest) ProtoMessage() {}

func (x *ForwardEchoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_echo_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ForwardEchoRequest.ProtoReflect.Descriptor instead.
func (*ForwardEchoRequest) Descriptor() ([]byte, []int) {
	return file_echo_proto_rawDescGZIP(), []int{3}
}

func (x *ForwardEchoRequest) GetCount() int32 {
	if x != nil {
		return x.Count
	}
	return 0
}

func (x *ForwardEchoRequest) GetQps() int32 {
	if x != nil {
		return x.Qps
	}
	return 0
}

func (x *ForwardEchoRequest) GetTimeoutMicros() int64 {
	if x != nil {
		return x.TimeoutMicros
	}
	return 0
}

func (x *ForwardEchoRequest) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *ForwardEchoRequest) GetHeaders() []*Header {
	if x != nil {
		return x.Headers
	}
	return nil
}

func (x *ForwardEchoRequest) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *ForwardEchoRequest) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *ForwardEchoRequest) GetHttp2() bool {
	if x != nil {
		return x.Http2
	}
	return false
}

func (x *ForwardEchoRequest) GetHttp3() bool {
	if x != nil {
		return x.Http3
	}
	return false
}

func (x *ForwardEchoRequest) GetServerFirst() bool {
	if x != nil {
		return x.ServerFirst
	}
	return false
}

func (x *ForwardEchoRequest) GetFollowRedirects() bool {
	if x != nil {
		return x.FollowRedirects
	}
	return false
}

func (x *ForwardEchoRequest) GetCert() string {
	if x != nil {
		return x.Cert
	}
	return ""
}

func (x *ForwardEchoRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *ForwardEchoRequest) GetCaCert() string {
	if x != nil {
		return x.CaCert
	}
	return ""
}

func (x *ForwardEchoRequest) GetCertFile() string {
	if x != nil {
		return x.CertFile
	}
	return ""
}

func (x *ForwardEchoRequest) GetKeyFile() string {
	if x != nil {
		return x.KeyFile
	}
	return ""
}

func (x *ForwardEchoRequest) GetCaCertFile() string {
	if x != nil {
		return x.CaCertFile
	}
	return ""
}

func (x *ForwardEchoRequest) GetInsecureSkipVerify() bool {
	if x != nil {
		return x.InsecureSkipVerify
	}
	return false
}

func (x *ForwardEchoRequest) GetAlpn() *Alpn {
	if x != nil {
		return x.Alpn
	}
	return nil
}

func (x *ForwardEchoRequest) GetServerName() string {
	if x != nil {
		return x.ServerName
	}
	return ""
}

func (x *ForwardEchoRequest) GetExpectedResponse() *wrappers.StringValue {
	if x != nil {
		return x.ExpectedResponse
	}
	return nil
}

type Alpn struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value []string `protobuf:"bytes,1,rep,name=value,proto3" json:"value,omitempty"`
}

func (x *Alpn) Reset() {
	*x = Alpn{}
	if protoimpl.UnsafeEnabled {
		mi := &file_echo_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Alpn) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Alpn) ProtoMessage() {}

func (x *Alpn) ProtoReflect() protoreflect.Message {
	mi := &file_echo_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Alpn.ProtoReflect.Descriptor instead.
func (*Alpn) Descriptor() ([]byte, []int) {
	return file_echo_proto_rawDescGZIP(), []int{4}
}

func (x *Alpn) GetValue() []string {
	if x != nil {
		return x.Value
	}
	return nil
}

type ForwardEchoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Output []string `protobuf:"bytes,1,rep,name=output,proto3" json:"output,omitempty"`
}

func (x *ForwardEchoResponse) Reset() {
	*x = ForwardEchoResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_echo_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ForwardEchoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ForwardEchoResponse) ProtoMessage() {}

func (x *ForwardEchoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_echo_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ForwardEchoResponse.ProtoReflect.Descriptor instead.
func (*ForwardEchoResponse) Descriptor() ([]byte, []int) {
	return file_echo_proto_rawDescGZIP(), []int{5}
}

func (x *ForwardEchoResponse) GetOutput() []string {
	if x != nil {
		return x.Output
	}
	return nil
}

var File_echo_proto protoreflect.FileDescriptor

var file_echo_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x65, 0x63, 0x68, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x27, 0x0a, 0x0b, 0x45, 0x63, 0x68, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x28, 0x0a, 0x0c,
	0x45, 0x63, 0x68, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x30, 0x0a, 0x06, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x97, 0x05, 0x0a, 0x12, 0x46, 0x6f, 0x72,
	0x77, 0x61, 0x72, 0x64, 0x45, 0x63, 0x68, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x71, 0x70, 0x73, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x03, 0x71, 0x70, 0x73, 0x12, 0x25, 0x0a, 0x0e, 0x74, 0x69, 0x6d, 0x65, 0x6f,
	0x75, 0x74, 0x5f, 0x6d, 0x69, 0x63, 0x72, 0x6f, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x0d, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x4d, 0x69, 0x63, 0x72, 0x6f, 0x73, 0x12, 0x10,
	0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c,
	0x12, 0x27, 0x0a, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x52, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x09, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x68,
	0x74, 0x74, 0x70, 0x32, 0x18, 0x07, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x68, 0x74, 0x74, 0x70,
	0x32, 0x12, 0x14, 0x0a, 0x05, 0x68, 0x74, 0x74, 0x70, 0x33, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x05, 0x68, 0x74, 0x74, 0x70, 0x33, 0x12, 0x20, 0x0a, 0x0b, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x46, 0x69, 0x72, 0x73, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x46, 0x69, 0x72, 0x73, 0x74, 0x12, 0x28, 0x0a, 0x0f, 0x66, 0x6f, 0x6c,
	0x6c, 0x6f, 0x77, 0x52, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x73, 0x18, 0x0e, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x0f, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x52, 0x65, 0x64, 0x69, 0x72, 0x65,
	0x63, 0x74, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x65, 0x72, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x63, 0x65, 0x72, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x0b,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x61, 0x43,
	0x65, 0x72, 0x74, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x63, 0x61, 0x43, 0x65, 0x72,
	0x74, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x65, 0x72, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x18, 0x10, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x65, 0x72, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x6b, 0x65, 0x79, 0x46, 0x69, 0x6c, 0x65, 0x18, 0x11, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x6b, 0x65, 0x79, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x61, 0x43, 0x65, 0x72,
	0x74, 0x46, 0x69, 0x6c, 0x65, 0x18, 0x12, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x61, 0x43,
	0x65, 0x72, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x2e, 0x0a, 0x12, 0x69, 0x6e, 0x73, 0x65, 0x63,
	0x75, 0x72, 0x65, 0x53, 0x6b, 0x69, 0x70, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x18, 0x13, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x12, 0x69, 0x6e, 0x73, 0x65, 0x63, 0x75, 0x72, 0x65, 0x53, 0x6b, 0x69,
	0x70, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x12, 0x1f, 0x0a, 0x04, 0x61, 0x6c, 0x70, 0x6e, 0x18,
	0x0d, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x6c,
	0x70, 0x6e, 0x52, 0x04, 0x61, 0x6c, 0x70, 0x6e, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x14, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x48, 0x0a, 0x10, 0x65, 0x78, 0x70, 0x65,
	0x63, 0x74, 0x65, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x15, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x52, 0x10, 0x65, 0x78, 0x70, 0x65, 0x63, 0x74, 0x65, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x1c, 0x0a, 0x04, 0x41, 0x6c, 0x70, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x22, 0x2d, 0x0a, 0x13, 0x46, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x45, 0x63, 0x68, 0x6f, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x75, 0x74, 0x70, 0x75,
	0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x32,
	0x88, 0x01, 0x0a, 0x0f, 0x45, 0x63, 0x68, 0x6f, 0x54, 0x65, 0x73, 0x74, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x2f, 0x0a, 0x04, 0x45, 0x63, 0x68, 0x6f, 0x12, 0x12, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x45, 0x63, 0x68, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x45, 0x63, 0x68, 0x6f, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x44, 0x0a, 0x0b, 0x46, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x45,
	0x63, 0x68, 0x6f, 0x12, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x46, 0x6f, 0x72, 0x77,
	0x61, 0x72, 0x64, 0x45, 0x63, 0x68, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x46, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x45, 0x63,
	0x68, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2e,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_echo_proto_rawDescOnce sync.Once
	file_echo_proto_rawDescData = file_echo_proto_rawDesc
)

func file_echo_proto_rawDescGZIP() []byte {
	file_echo_proto_rawDescOnce.Do(func() {
		file_echo_proto_rawDescData = protoimpl.X.CompressGZIP(file_echo_proto_rawDescData)
	})
	return file_echo_proto_rawDescData
}

var file_echo_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_echo_proto_goTypes = []interface{}{
	(*EchoRequest)(nil),          // 0: proto.EchoRequest
	(*EchoResponse)(nil),         // 1: proto.EchoResponse
	(*Header)(nil),               // 2: proto.Header
	(*ForwardEchoRequest)(nil),   // 3: proto.ForwardEchoRequest
	(*Alpn)(nil),                 // 4: proto.Alpn
	(*ForwardEchoResponse)(nil),  // 5: proto.ForwardEchoResponse
	(*wrappers.StringValue)(nil), // 6: google.protobuf.StringValue
}
var file_echo_proto_depIdxs = []int32{
	2, // 0: proto.ForwardEchoRequest.headers:type_name -> proto.Header
	4, // 1: proto.ForwardEchoRequest.alpn:type_name -> proto.Alpn
	6, // 2: proto.ForwardEchoRequest.expectedResponse:type_name -> google.protobuf.StringValue
	0, // 3: proto.EchoTestService.Echo:input_type -> proto.EchoRequest
	3, // 4: proto.EchoTestService.ForwardEcho:input_type -> proto.ForwardEchoRequest
	1, // 5: proto.EchoTestService.Echo:output_type -> proto.EchoResponse
	5, // 6: proto.EchoTestService.ForwardEcho:output_type -> proto.ForwardEchoResponse
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_echo_proto_init() }
func file_echo_proto_init() {
	if File_echo_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_echo_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EchoRequest); i {
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
		file_echo_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EchoResponse); i {
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
		file_echo_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Header); i {
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
		file_echo_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ForwardEchoRequest); i {
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
		file_echo_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Alpn); i {
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
		file_echo_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ForwardEchoResponse); i {
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
			RawDescriptor: file_echo_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_echo_proto_goTypes,
		DependencyIndexes: file_echo_proto_depIdxs,
		MessageInfos:      file_echo_proto_msgTypes,
	}.Build()
	File_echo_proto = out.File
	file_echo_proto_rawDesc = nil
	file_echo_proto_goTypes = nil
	file_echo_proto_depIdxs = nil
}
