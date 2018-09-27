package multiaddrsam

import (
	"fmt"
	"strings"

	. "github.com/eyedeekay/sam3"
	ma "github.com/multiformats/go-multiaddr"
)

type I2PMultiaddr struct {
	Name             string
	Code             int
	baseMultiAddress ma.Multiaddr
	I2PAddr
}

// These were picked arbitrarily. They will probably change.
var P_GARLIC_NTCP = 445
var P_GARLIC_SSU = 890

func (addr I2PMultiaddr) Address() *I2PAddr {
	return &addr.I2PAddr
}

//
func (addr I2PMultiaddr) Encapsulate(multiaddr ma.Multiaddr) ma.Multiaddr {
	if strings.Contains(addr.String(), multiaddr.String()) {
		return addr
	}
	i2pAddrString := "/ntcp/" + addr.String()
	multiAddrString := "/" + multiaddr.Protocols()[0].Name + "/" + multiaddr.String()
	baddr, _ := NewI2PMultiaddr(i2pAddrString + multiAddrString)
	return baddr
}

func (addr I2PMultiaddr) Decapsulate(multiaddr ma.Multiaddr) ma.Multiaddr {
	i2pAddrString := "/ntcp/" + addr.String()
	var multiAddrString string
	for _, mp := range multiaddr.Protocols() {
		if mp.Name != "" {
			multiAddrString += "/" + mp.Name + "/" + multiaddr.String()
		}
	}
	tmp := strings.Replace(multiAddrString, i2pAddrString, "", -1)
	baddr, _ := ma.NewMultiaddr(tmp)
	return baddr
}

func (addr I2PMultiaddr) Protocols() []ma.Protocol {
	p := []ma.Protocol{}
	p = append(p, ma.Protocol{Code: P_GARLIC_NTCP, Name: addr.Name, Size: 31})
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
		return string(addr.I2PAddr.String()), nil
	}
	if code == P_GARLIC_SSU {
		return string(addr.I2PAddr.String()), nil
	}
	return addr.baseMultiAddress.ValueForProtocol(code)
}

func NewI2PMultiaddr(inputs string) (I2PMultiaddr, error) {
	var m I2PMultiaddr
	var err error
	if i := strings.SplitN(inputs, "/ntcp/", 2); len(i) == 2 {
		s := strings.Split(i[1], "/")
		m.I2PAddr, err = NewI2PAddrFromString(s[0])
		if err != nil {
			return m, err
		}
		m.baseMultiAddress = m.Decapsulate(m)
		m.Name = "ntcp"
		m.Code = P_GARLIC_NTCP
		return m, err
	}
	if i := strings.SplitN(inputs, "/ssu/", 2); len(i) == 2 {
		s := strings.Split(i[1], "/")
		m.I2PAddr, err = NewI2PAddrFromString(s[0])
		if err != nil {
			return m, err
		}
		m.baseMultiAddress = m.Decapsulate(m)
		m.Name = "ssu"
		m.Code = P_GARLIC_SSU
		return m, fmt.Errorf("sam3-multiaddr Error: %s, %s", "ssu isn't implemented yet. Come back later.", s[0])
	}
	return m, fmt.Errorf("sam3-multiaddr Error: %s", "Not an i2p Multiaddr")
}
