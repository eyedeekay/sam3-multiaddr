package multiaddrsam

import (
	. "github.com/eyedeekay/sam3"
	"log"
	"testing"
)

func TestMain(t *testing.T) {
	k, e := createEepServiceKey()
	if e != nil {
		log.Println(e)
		t.Fatal(e.Error())
	}
	NewI2PMultiaddr("/ntcp/" + k.String())
}

func createEepServiceKey() (*I2PKeys, error) {
	sam, err := NewSAM("127.0.0.1:7656")
	if err != nil {
		return nil, err
	}
	defer sam.Close()
	k, err := sam.NewKeys()
	return &k, err
}
