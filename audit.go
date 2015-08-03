package lavinia

import (
	"crypto/sha256"
	"time"
)

type ServerPayment struct {
	t      time.Time
	secret []byte
}

type AuditBlock struct {
	t      time.Time
	lookup [sha256.Size]byte //encrypted
	secret []byte            //encrypted
	rprev  []byte
	r      []byte
	rs     []byte
}
