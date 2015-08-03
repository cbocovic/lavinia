package lavinia

import (
	"crypto/sha256"
	"fmt"
	"github.com/cbocovic/chord"
	"github.com/cbocovic/chordFS"
	"os"
)

const (
	code byte = 3
)

type LaviniaServer struct {
	addr string
	node *chord.ChordNode
	fs   *fs.FileSystem

	storage         string
	pendingAudits   []AuditBlock
	pendingPayments []ServerPayment
}

//error checking function
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
	}
}

func Create(home string, addr string) *LaviniaServer {
	me := new(LaviniaServer)
	me.node = chord.Create(addr)

	if me.node == nil {
		return nil
	}

	me.fs = fs.Extend(home, addr, me.node)
	if me.fs == nil {
		return nil
	}

	me.storage = fmt.Sprintf("%s/pmts", home)
	err := os.MkdirAll(me.storage, 0755)

	if err != nil {
		checkError(err)
		return nil
	}

	me.node.Register(code, me)
	return me
}

func Join(home string, myaddr string, addr string) *LaviniaServer {
	me := new(LaviniaServer)
	me.node = chord.Join(myaddr, addr)
	if me.node == nil {
		return nil
	}

	me.fs = fs.Extend(home, addr, me.node)
	if me.fs == nil {
		return nil
	}

	me.storage = fmt.Sprintf("%s/pmts", home)
	err := os.MkdirAll(me.storage, 0755)
	if err != nil {
		checkError(err)
		return nil
	}

	me.node.Register(code, me)
	return me
}

func (me *LaviniaServer) Notify(id [sha256.Size]byte, myid [sha256.Size]byte) {
	//TODO: forward payment info to new node

}

func (me *LaviniaServer) Message(data []byte) []byte {
	//TODO: parse appropriately
	return nullMsg()

}

func Retrieve() {

}

func (me *LaviniaServer) Finalize() {
	me.fs.Finalize()
}

/** Printouts of information **/

func (me *LaviniaServer) Info() string {
	return me.node.Info()
}

func (me *LaviniaServer) ShowFingers() string {
	return me.node.ShowFingers()
}

func (me *LaviniaServer) ShowSucc() string {
	return me.node.ShowSucc()
}
