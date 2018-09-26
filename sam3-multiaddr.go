package multiaddrsam

import (
	. "github.com/eyedeekay/sam3"
	ma "github.com/multiformats/go-multiaddr"
	"strings"
)

type I2PMultiaddr struct {
	//address          string
	baseMultiAddress ma.Multiaddr
	I2PAddr
}

var P_GARLIC_NTCP = 445

func (addr I2PMultiaddr) Address() *I2PAddr {
	return &addr.I2PAddr
}

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
    p = append(p, ma.Protocol{ Code: P_GARLIC_NTCP, Name: "ntcp", Size: 31})
    p = append(p, addr.baseMultiAddress.Protocols()...)
	return p
}

func (addr I2PMultiaddr) Equal(multiaddr ma.Multiaddr) bool {
	if multiaddr.String() == addr.String() {
		return true
	}
	return false
}

func (addr I2PMultiaddr) ValueForProtocol(code int) (string, error) {
	if code == P_GARLIC_NTCP {
		return addr.String(), nil
	}
	return addr.baseMultiAddress.ValueForProtocol(code)
}

func NewI2PMultiaddr(s string) (I2PMultiaddr, error) {
	var multiAddress ma.Multiaddr
	m := multiAddress.(I2PMultiaddr)
	var err error
	m.I2PAddr, err = NewI2PAddrFromString(s)
	if err != nil {
		return m, err
	}
	address, err := m.ValueForProtocol(P_GARLIC_NTCP)
	if err != nil {
		return m, err
	}
	addressAsMultiAddress, err := ma.NewMultiaddr("/garlic/" + address)
	if err != nil {
		return addressAsMultiAddress.(I2PMultiaddr), err
	}
	m.baseMultiAddress = multiAddress.Decapsulate(addressAsMultiAddress)
	return m, err
}
