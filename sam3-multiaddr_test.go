package multiaddrsam

import (
	. "github.com/eyedeekay/sam3"
	"log"
	"testing"
)

func TestNTCP(t *testing.T) {
	k, e := createEepServiceKey()
	if e != nil {
		log.Println(e)
		t.Fatal(e.Error())
	}
	NewI2PMultiaddr("/ntcp/" + k.String())
	log.Println("Successfully ran the ntcp test", k.String())
}

func TestSSU(t *testing.T) {
	k, e := createEepServiceKey()
	if e != nil {
		log.Println(e)
		t.Fatal(e.Error())
	}
	NewI2PMultiaddr("/ssu/" + k.String())
	log.Println("Successfully ran the ssu test", k.String())
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
