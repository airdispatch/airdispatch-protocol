// Code generated by protoc-gen-go.
// source: source/server.proto
// DO NOT EDIT!

package wire

import proto "code.google.com/p/goprotobuf/proto"
import json "encoding/json"
import math "math"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type TransferMessage struct {
	Name             *string `protobuf:"bytes,1,req,name=name" json:"name,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *TransferMessage) Reset()         { *m = TransferMessage{} }
func (m *TransferMessage) String() string { return proto.CompactTextString(m) }
func (*TransferMessage) ProtoMessage()    {}

func (m *TransferMessage) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

type TransferMessageList struct {
	LastUpdated      *uint64 `protobuf:"varint,1,req,name=last_updated" json:"last_updated,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *TransferMessageList) Reset()         { *m = TransferMessageList{} }
func (m *TransferMessageList) String() string { return proto.CompactTextString(m) }
func (*TransferMessageList) ProtoMessage()    {}

func (m *TransferMessageList) GetLastUpdated() uint64 {
	if m != nil && m.LastUpdated != nil {
		return *m.LastUpdated
	}
	return 0
}

type MessageDescription struct {
	Location         *string `protobuf:"bytes,1,req,name=location" json:"location,omitempty"`
	Name             *string `protobuf:"bytes,2,req,name=name" json:"name,omitempty"`
	Nonce            *uint64 `protobuf:"varint,3,opt,name=nonce" json:"nonce,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *MessageDescription) Reset()         { *m = MessageDescription{} }
func (m *MessageDescription) String() string { return proto.CompactTextString(m) }
func (*MessageDescription) ProtoMessage()    {}

func (m *MessageDescription) GetLocation() string {
	if m != nil && m.Location != nil {
		return *m.Location
	}
	return ""
}

func (m *MessageDescription) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *MessageDescription) GetNonce() uint64 {
	if m != nil && m.Nonce != nil {
		return *m.Nonce
	}
	return 0
}

type MessageList struct {
	Messages         []*MessageDescription `protobuf:"bytes,1,rep,name=messages" json:"messages,omitempty"`
	XXX_unrecognized []byte                `json:"-"`
}

func (m *MessageList) Reset()         { *m = MessageList{} }
func (m *MessageList) String() string { return proto.CompactTextString(m) }
func (*MessageList) ProtoMessage()    {}

func (m *MessageList) GetMessages() []*MessageDescription {
	if m != nil {
		return m.Messages
	}
	return nil
}

func init() {
}