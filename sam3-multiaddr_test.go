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
	x, e := NewI2PMultiaddr("/ntcp/"+k.String(), true, "/sam/127.0.0.1:7657")
	if e != nil {
		log.Println(e)
		t.Fatal(e.Error())
	}
	log.Printf("Successfully ran the ntcp test\n  %s\n  %s\n", k.String(), x.String())
	log.Println("  ", x.Protocols())
}

func TestSSU(t *testing.T) {
	k, e := createEepServiceKey()
	if e != nil {
		log.Println(e)
		t.Fatal(e.Error())
	}
	x, e := NewI2PMultiaddr("/ssu/"+k.String(), true, "/sam/127.0.0.1:7657")
	if e == nil {
		log.Println(e)
		t.Fatal(e.Error())
	}
	log.Printf("Successfully ran the ntcp test\n  %s\n  %s\n", k.String(), x.String())
	log.Println("  ", x.Protocols())
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
