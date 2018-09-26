package multiaddrsam

import (
	. "github.com/eyedeekay/sam3"
	ma "github.com/multiformats/go-multiaddr"
	"strings"
)

type I2PMultiaddr struct {
	address          string
	baseMultiAddress ma.Multiaddr
	I2PAddr
}

var Garlic = 445

//
func (addr I2PMultiaddr) Encapsulate(multiaddr ma.Multiaddr) ma.Multiaddr {
	if !strings.Contains(addr.String(), multiaddr.String()) {
		addr, _ = NewI2PMultiaddr(addr.String() + multiaddr.String())
	}
	return addr
}

func (addr I2PMultiaddr) Decapsulate(multiaddr ma.Multiaddr) ma.Multiaddr {
	addr, _ = NewI2PMultiaddr(strings.Replace(addr.String(), multiaddr.String(), "", -1))
	return addr
}

func (addr I2PMultiaddr) Protocols() []ma.Protocol {
	p := []ma.Protocol{}
	return p
}

func (addr I2PMultiaddr) Equal(multiaddr ma.Multiaddr) bool {
	if multiaddr.String() == addr.String() {
		return true
	}
	return false
}

func (addr I2PMultiaddr) ValueForProtocol(code int) (string, error) {
	if code == Garlic {
		return addr.String(), nil
	}
	return addr.baseMultiAddress.ValueForProtocol(code)
}

func NewI2PMultiaddr(s string) (I2PMultiaddr, error) {
	multiAddress, err := ma.NewMultiaddr(s)
	if err != nil {
		return multiAddress.(I2PMultiaddr), err
	}
	m := multiAddress.(I2PMultiaddr)
	m.address, err = m.ValueForProtocol(Garlic)
	if err != nil {
		return m, err
	}
	addressAsMultiAddress, err := ma.NewMultiaddr("/garlic/" + m.address)
	if err != nil {
		return addressAsMultiAddress.(I2PMultiaddr), err
	}
	m.baseMultiAddress = multiAddress.Decapsulate(addressAsMultiAddress)
	return m, err
}
