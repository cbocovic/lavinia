package lavinia

//#include "galois.h"
import "C"

import (
	"crypto/rand"
	"fmt"
	"io"
	"unsafe"
)

//split takes a document and splits it into shares using a
//k out of n shamir secret sharing scheme. It returns the shares
//in a 2D array, or nil and a non-nil error
func split(document []byte, n int, k int) ([][]byte, error) {
	fmt.Printf("document: %x.\n", document)
	length := len(document)
	fmt.Printf("Document has length %d.\n", length)
	shares := make([][]byte, n)
	for i := 0; i < n; i++ {
		//the length of each share is a 2 byte x value plus
		//the length of the document
		shares[i] = make([]byte, len(document)+2)
		//populate shares with x values
		x := shares[i][:2]
		if _, err := io.ReadFull(rand.Reader, x); err != nil {
			checkError(err)
			return nil, err
		}
	}

	//choose k elements uniformly at random
	as := make([][]byte, k)
	for i := 0; i < k; i++ {
		as[i] = make([]byte, 2)
		if _, err := io.ReadFull(rand.Reader, as[i]); err != nil {
			checkError(err)
			return nil, err
		}
	}

	//operate on 2 bytes at a time
	for i := 0; i < length-1; i += 2 {
		//working on document[i:i+2].

		for j := 0; j < n; j++ {
			x := C.int(bytes2int(shares[j][:2]))
			y := C.int(0)
			for l := 0; l < k; l++ {
				a := C.int(bytes2int(as[l]))
				tmp := C.int(1)
				for m := 0; m <= l; m++ {
					tmp = C.galois_single_multiply(x, tmp, 16)
				}
				y = y ^ C.galois_single_multiply(a, tmp, 16)
			}

			ybytes := int2bytes(int(y))
			out := shares[j][i+2 : i+4]
			out[0] = document[i] ^ ybytes[0]
			out[1] = document[i+1] ^ ybytes[1]
		}
	}

	return shares, nil

}

func mend(shares [][]byte) []byte {

	length := len(shares[0])
	fmt.Printf("Shares have length = %d.\n", length)
	num := len(shares)
	document := make([]byte, length-2)
	for i := 2; i < length-1; i += 2 {
		//interpolate shares[][i:i+2]
		out := document[i-2 : i]
		out[0] = 0
		out[1] = 0
		cout := C.int(0)

		for j := 0; j < num; j++ {
			b := C.int(1)
			tmp := C.int(0)
			for k := 0; k < num; k++ {
				if k != j {
					xk := C.int(bytes2int(shares[k][:2]))
					xj := C.int(bytes2int(shares[j][:2]))
					b = C.galois_single_multiply(b, xk, 16)
					tmp = xk ^ xj
					b = C.galois_single_divide(b, tmp, 16)
				}
			}
			ttmp := make([]byte, 2)
			ctmp := C.CString(string(ttmp))
			C.galois_w16_region_multiply(C.CString(string(shares[j][i:i+2])), b, 2, ctmp, 0)

			wootmp := bytes2int([]byte(C.GoBytes(unsafe.Pointer(ctmp), C.int(2))))
			tmpint1 := C.int(wootmp)
			cout = tmpint1 ^ cout
		}
		out = int2bytes(int(cout))

		copy(document[i-2:i], out)
	}
	fmt.Printf("document: %x.\n", document)
	fmt.Printf("Document has length %d.\n", len(document))
	return document
}

func bytes2int(bytes []byte) int {
	if len(bytes) < 2 {
		if len(bytes) < 1 {
			return 0
		}
		return int(bytes[0])
	}
	return int(bytes[0]) + int(bytes[1])*256

}

func int2bytes(myint int) []byte {
	return []byte{byte(myint % 256), byte(myint >> 8)}

}
