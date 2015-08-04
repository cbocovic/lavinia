package lavinia

//#include "galois.h"
import "C"

import (
	"crypto/rand"
	"fmt"
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
			cout := C.CString(string(out))
			//TODO: This probably does not work as intended...
			C.galois_region_xor(C.CString(string(document[i:i+2])), C.CString(string(ybytes)), cout, 2)
			copy(out, []byte(C.GoString(cout)))
		}
	}

	return shares, nil

}

func mend(shares [][]byte) []byte {
	length := len(shares[0])
	fmt.Printf("length = %d.\n", length)
	num := len(shares)
	document := make([]byte, length-2)
	for i := 2; i < length-2; i += 2 {
		//interpolate shares[][i:i+2]
		out := document[i-2 : i]
		out[0] = 0
		out[1] = 0
		//cout := C.CString(string(out))
		cout := C.int(0)
		fmt.Printf("Reconstructing part %d of document.\n", i)

		for j := 0; j < num; j++ {
			b := C.int(1)
			tmp := C.int(0)
			//fmt.Printf("share j=%d.\n", j)
			for k := 0; k < num; k++ {
				//fmt.Printf("looping k=%d.\n", k)
				if k != j {
					xk := C.int(bytes2int(shares[k][:2]))
					//fmt.Printf("xk was %d.\n", int(xk))
					xj := C.int(bytes2int(shares[j][:2]))
					//fmt.Printf("xj was %d.\n", int(xj))
					b = C.galois_single_multiply(b, xk, 16)
					//fmt.Printf("b was %d.\n", int(b))
					tmp = xk ^ xj
					b = C.galois_single_divide(b, tmp, 16)
					//fmt.Printf("b was %d.\n", int(b))
				}
			}
			fmt.Printf("share piece was %x.\n", shares[j][i:i+2])
			fmt.Printf("b was %x.\n", int2bytes(int(b)))
			fmt.Printf("cout was: %x.\n", int2bytes(int(cout)))
			ttmp := make([]byte, 2)
			ctmp := C.CString(string(ttmp))
			C.galois_w16_region_multiply(C.CString(string(shares[j][i:i+2])), b, 2, ctmp, 0)
			fmt.Printf("ctemp is: %x.\n", []byte(C.GoString(ctmp)))
			//C.galois_w16_region_multiply(C.CString(string(shares[j][i:i+2])), b, 2, cout, 1)
			//fmt.Printf("cout is: %x.\n", []byte(C.GoString(cout)))
			wootmp := bytes2int([]byte(C.GoString(ctmp)))
			tmpint1 := C.int(wootmp)
			cout = tmpint1 ^ cout
			fmt.Printf("cout is: %x.\n", int2bytes(int(cout)))
		}
		//copy(out, []byte(C.GoString(cout)))
		out = int2bytes(int(cout))
		fmt.Printf("out: %x.\n", out)
		copy(document[i-1:i], out)
	}
	fmt.Printf("document: %x.\n", document)
	return document
}

func bytes2int(bytes []byte) int {
	if len(bytes) < 2 {
		if len(bytes) < 1 {
			return 0
		}
		fmt.Printf("lol\n")
		return int(bytes[0])
	}
	return int(bytes[0]) + int(bytes[1])*256

}

func int2bytes(myint int) []byte {
	return []byte{byte(myint % 256), byte(myint >> 8)}

}
