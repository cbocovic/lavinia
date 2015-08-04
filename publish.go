package lavinia

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/cbocovic/chordFS"
	"io"
	"os"
)

//Publish will split the document at path into shares (as stated
//in the lavinia protocol) and store these shares in the DHT by
//contacting the give address
func Publish(path string, addr string) error {

	//make temprorary lavinia directory to store all files
	err := os.MkdirAll("lavinia(tmp)", 0755)
	if err != nil {
		checkError(err)
		return nil
	}

	//TODO: encrypt file and save key and ciphertext
	//generate key
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		checkError(err)
		return nil
	}
	file, err := os.Create("lavinia(tmp)/key")
	checkError(err)
	_, err = file.Write(key)
	checkError(err)
	if err != nil {
		return err
	}
	fmt.Printf("saved key.\n")
	file.Close()

	//read in plaintext
	file, err = os.Open(path)
	checkError(err)
	plaintext := make([]byte, 4096)

	n, err := file.Read(plaintext)
	if err != nil {
		checkError(err)
		return err
	}
	file.Close()

	//pad plaintext to be mulitple of blocksize
	if n%aes.BlockSize != 0 {
		needed := aes.BlockSize - n%aes.BlockSize
		for i := n; i < needed; i++ {
			plaintext[i] = 0 //TODO:make more legit
		}
		n += needed
	}
	plaintext = plaintext[:n]

	//encrypt
	block, err := aes.NewCipher(key)
	if err != nil {
		checkError(err)
		return err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		checkError(err)
		return err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	//write ciphertext to file
	file, err = os.Create("lavinia(tmp)/ciphertext")
	checkError(err)
	_, err = file.Write(ciphertext)
	checkError(err)
	if err != nil {
		return err
	}
	fmt.Printf("saved ciphertext.\n")
	file.Close()

	//TODO: split encrypted file into shares
	shares, err := split(ciphertext, 7, 5)
	for i, share := range shares {
		file, err = os.Create(fmt.Sprintf("lavinia(tmp)/share%d", i))
		checkError(err)
		_, err = file.Write(share)
		file.Close()
		checkError(err)
	}

	/*TODO:tmp: mend file
	mended := mend(shares)
	file, err = os.Create("lavinia(tmp)/mended")
	checkError(err)
	_, err = file.Write(mended)
	checkError(err)
	if err != nil {
		return err
	}
	fmt.Printf("saved mended.\n")
	file.Close()
	*/

	//TODO: store all pieces in DHT
	//choose random keys and enter them into a manifest
	manifest := make([]byte, 7*sha256.Size)
	for i, _ := range shares {
		var tmpKey [sha256.Size]byte
		shareKey := manifest[sha256.Size*i : sha256.Size*(i+1)]
		if _, err = io.ReadFull(rand.Reader, shareKey); err != nil {
			checkError(err)
			return err
		}
		copy(tmpKey[:], shareKey)
		fs.Store(tmpKey, fmt.Sprintf("lavinia(tmp)/share%d", i), addr)
	}
	file, err = os.Create("lavinia(tmp)/manifest")
	checkError(err)
	file.Write(manifest)
	file.Close()

	//store manifest file
	fs.Store(sha256.Sum256([]byte("00test")), "lavinia(tmp)/manifest", addr)

	//store key file
	fs.Store(sha256.Sum256([]byte("01test")), "lavinia(tmp)/key", addr)

	//TODO: craft payments and store

	return nil
}

//encrypt saves a ciphertext version of the file and the key in the tmp directory
