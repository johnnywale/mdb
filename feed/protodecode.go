package feed

import (
	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	//	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
)

type ProtoDecoder struct {
	fc   *desc.FileDescriptor
	keys map[string]*desc.MessageDescriptor
}

func (p *ProtoDecoder) Init(bb []byte) {
	var fds dpb.FileDescriptorSet

	err := proto.Unmarshal(bb, &fds)
	if err != nil {
		panic(err)
	}
	fc, _ := desc.CreateFileDescriptorFromSet(&fds)
	p.fc = fc
	p.keys = make(map[string]*desc.MessageDescriptor)
	md := fc.GetMessageTypes()
	for _, val := range md {
		p.keys[val.GetFullyQualifiedName()] = val
	}
}

func (p *ProtoDecoder) GetInstance(messageType string) *dynamic.Message {
	return dynamic.NewMessage(p.keys[messageType])
}

func (p *ProtoDecoder) Decode(message *dynamic.Message, bs []byte) {
	proto.Unmarshal(bs, message)
}
