package lavinia

import (
	"github.com/golang/protobuf/proto"
	"log"
)

func nullMsg() []byte {
	msg := new(NetworkMessage)
	msg.Proto = proto.Uint32(2)

	data, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	return data
}
