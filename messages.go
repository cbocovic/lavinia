package lavinia

import (
	"github.com/golang/protobuf/proto"
	"log"
)

func storepmtMsg() []byte {
	msg := new(NetworkMessage)
	msg.Proto = proto.Uint32(3)

	data, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	return data

}

func auditMsg() []byte {
	msg := new(NetworkMessage)
	msg.Proto = proto.Uint32(3)

	data, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	return data

}

func payMsg() []byte {
	msg := new(NetworkMessage)
	msg.Proto = proto.Uint32(3)

	data, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	return data

}

func nullMsg() []byte {
	msg := new(NetworkMessage)
	msg.Proto = proto.Uint32(3)

	data, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	return data
}
