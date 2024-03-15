// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: axiome/v1beta1/events.proto

package util

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type EventMessage struct {
	Module string `protobuf:"bytes,1,opt,name=module,proto3" json:"module,omitempty"`
	Action string `protobuf:"bytes,2,opt,name=action,proto3" json:"action,omitempty"`
}

func (m *EventMessage) Reset()         { *m = EventMessage{} }
func (m *EventMessage) String() string { return proto.CompactTextString(m) }
func (*EventMessage) ProtoMessage()    {}
func (*EventMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_12a349e2e2ef6f29, []int{0}
}
func (m *EventMessage) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EventMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EventMessage.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EventMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventMessage.Merge(m, src)
}
func (m *EventMessage) XXX_Size() int {
	return m.Size()
}
func (m *EventMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_EventMessage.DiscardUnknown(m)
}

var xxx_messageInfo_EventMessage proto.InternalMessageInfo

func (m *EventMessage) GetModule() string {
	if m != nil {
		return m.Module
	}
	return ""
}

func (m *EventMessage) GetAction() string {
	if m != nil {
		return m.Action
	}
	return ""
}

func init() {
	proto.RegisterType((*EventMessage)(nil), "axiome.v1beta1.EventMessage")
}

func init() { proto.RegisterFile("axiome/v1beta1/events.proto", fileDescriptor_12a349e2e2ef6f29) }

var fileDescriptor_12a349e2e2ef6f29 = []byte{
	// 171 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4e, 0xac, 0xc8, 0xcc,
	0xcf, 0x4d, 0xd5, 0x2f, 0x33, 0x4c, 0x4a, 0x2d, 0x49, 0x34, 0xd4, 0x4f, 0x2d, 0x4b, 0xcd, 0x2b,
	0x29, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x83, 0x48, 0xea, 0x41, 0x25, 0x95, 0xec,
	0xb8, 0x78, 0x5c, 0x41, 0xf2, 0xbe, 0xa9, 0xc5, 0xc5, 0x89, 0xe9, 0xa9, 0x42, 0x62, 0x5c, 0x6c,
	0xb9, 0xf9, 0x29, 0xa5, 0x39, 0xa9, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x50, 0x1e, 0x48,
	0x3c, 0x31, 0xb9, 0x24, 0x33, 0x3f, 0x4f, 0x82, 0x09, 0x22, 0x0e, 0xe1, 0x39, 0x59, 0x9c, 0x78,
	0x24, 0xc7, 0x78, 0xe1, 0x91, 0x1c, 0xe3, 0x83, 0x47, 0x72, 0x8c, 0x13, 0x1e, 0xcb, 0x31, 0x5c,
	0x78, 0x2c, 0xc7, 0x70, 0xe3, 0xb1, 0x1c, 0x43, 0x94, 0x5c, 0x52, 0x66, 0x49, 0x52, 0x69, 0x72,
	0x76, 0x6a, 0x89, 0x5e, 0x7e, 0x51, 0xba, 0x7e, 0x8a, 0x41, 0x9e, 0x49, 0x66, 0xb6, 0x7e, 0x62,
	0x45, 0x6e, 0x8a, 0x7e, 0x69, 0x49, 0x66, 0x4e, 0x12, 0x1b, 0xd8, 0x41, 0xc6, 0x80, 0x00, 0x00,
	0x00, 0xff, 0xff, 0xc4, 0x85, 0x56, 0x42, 0xaf, 0x00, 0x00, 0x00,
}

func (m *EventMessage) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EventMessage) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EventMessage) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Action) > 0 {
		i -= len(m.Action)
		copy(dAtA[i:], m.Action)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.Action)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Module) > 0 {
		i -= len(m.Module)
		copy(dAtA[i:], m.Module)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.Module)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintEvents(dAtA []byte, offset int, v uint64) int {
	offset -= sovEvents(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *EventMessage) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Module)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	l = len(m.Action)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	return n
}

func sovEvents(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozEvents(x uint64) (n int) {
	return sovEvents(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *EventMessage) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvents
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: EventMessage: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EventMessage: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Module", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Module = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Action", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Action = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipEvents(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvents
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipEvents(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowEvents
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowEvents
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowEvents
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthEvents
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupEvents
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthEvents
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthEvents        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowEvents          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupEvents = fmt.Errorf("proto: unexpected end of group")
)