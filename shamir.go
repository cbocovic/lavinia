package lavinia

//#include "galois.h"
import "C"

import (
	"crypto/rand"
	"io"
)

//split takes a document and splits it into shares using a
//k out of n shamir secret sharing scheme. It returns the shares
//in a 2D array, or nil and a non-nil error
func split(document []byte, n int, k int) ([][]byte, error) {
	length := len(document)
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
			//TODO: This probably does not work as intended...
			C.galois_region_xor(C.CString(string(document[i:i+2])), C.CString(string(ybytes)), C.CString(string(out)), 2)
		}
	}

	return shares, nil

}

func bytes2int(bytes []byte) int {
	return int(bytes[0]) + int(bytes[1])*256

}

func int2bytes(myint int) []byte {
	return []byte{byte(myint % 256), byte(myint >> 8)}

}
