package lavinia

import (
	"crypto/sha256"
	"fmt"
	"github.com/cbocovic/chordFS"
	"os"
)

//Retrieve will look up a keyword in the DHT at addr and save it to path
func Retrieve(keyword string, path string, addr string) error {
	//make temprorary lavinia directory to store all files
	err := os.MkdirAll("lavinia2(tmp)", 0755)
	if err != nil {
		checkError(err)
		return nil
	}

	//look up manifests
	key := sha256.Sum256([]byte("00" + keyword))
	err = fs.Fetch(key, "lavinia2(tmp)/manifest", addr)
	checkError(err)
	if err != nil {
		return err
	}

	key = sha256.Sum256([]byte("01" + keyword))
	err = fs.Fetch(key, "lavinia2(tmp)/key", addr)
	checkError(err)
	if err != nil {
		return err
	}

	//loop through manifest and look up shares
	file, err := os.Open("lavinia2(tmp)/manifest")
	checkError(err)
	manifest := make([]byte, sha256.Size*7)
	_, err = file.Read(manifest)
	checkError(err)
	if err != nil {
		return err
	}
	file.Close()

	ctr := 0
	for i := 0; i < 7; i++ {
		copy(key[:], manifest[i*sha256.Size:(i+1)*sha256.Size])
		err = fs.Fetch(key, fmt.Sprintf("lavinia2(tmp)/share%d", i), addr)
		checkError(err)
		if err == nil {
			fmt.Printf("Retrieved share %d of 7.\n", i)
			ctr++
		}
	}

	//read in shares
	shares := make([][]byte, ctr)
	for i, _ := range shares {
		shares[i] = make([]byte, 4096)
		file, err = os.Open(fmt.Sprintf("lavinia2(tmp)/share%d", i))
		checkError(err)
		n, err := file.Read(shares[i])
		checkError(err)
		shares[i] = shares[i][:n]
		if n == 0 {
			fmt.Printf("Something weird happened.\n")
		} else {
			fmt.Printf("Length of share %d is %d. Read in %d bytes.\n", i, len(shares[i]), n)
		}
		file.Close()
	}

	//interpolate

	fmt.Printf("Length of share 0 is %d.\n", len(shares[0]))
	mended := mend(shares)
	file, err = os.Create("lavinia2(tmp)/mended")
	checkError(err)
	_, err = file.Write(mended)
	checkError(err)
	if err != nil {
		return err
	}
	fmt.Printf("saved mended.\n")
	file.Close()

	//TODO: decrypt file

	return nil

}
