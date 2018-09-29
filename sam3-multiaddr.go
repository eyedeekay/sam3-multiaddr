package multiaddrsam

import (
    "encoding/binary"
	"fmt"
	"strings"

	. "github.com/eyedeekay/sam3"
	ma "github.com/multiformats/go-multiaddr"
)

type I2PMultiaddr struct {
	Name             string
	Code             int
    VCode            []byte
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
func (addr I2PMultiaddr) Bytes() []byte {
	return []byte("/ntcp/" + addr.Address().String())
}

func (addr I2PMultiaddr) String() string {
    return string(addr.Bytes())
}

//
func (addr I2PMultiaddr) Encapsulate(multiaddr ma.Multiaddr) ma.Multiaddr {
	mb := addr.Bytes()
	ob := multiaddr.Bytes()

	rb := make([]byte, len(mb)+len(ob))

	copy(rb, mb)
	copy(rb[len(mb):], ob)
	baddr, _ := NewI2PMultiaddr(string(rb))
	return baddr
}

func (addr I2PMultiaddr) Decapsulate(multiaddr ma.Multiaddr) ma.Multiaddr {
	ms := string(addr.Bytes())
	os := string(multiaddr.Bytes())

	if i := strings.LastIndex(ms, os); i > 0 {
		baddr, _ := ma.NewMultiaddr(ms[:i])
		return baddr
	}

	baddr, _ := ma.NewMultiaddr(addr.String())
	return baddr
}

func (addr I2PMultiaddr) Protocols() []ma.Protocol {
	p := []ma.Protocol{}
	p = append(p, ma.Protocol{
        Code: addr.Code,
        Name: addr.Name,
        Size: 0,
        Path: false,
    })
	if addr.baseMultiAddress != nil {
		for _, mp := range addr.baseMultiAddress.Protocols() {
			if mp.Name != "" {
				p = append(p, mp)
			}
		}
	}
	return p
}

func (addr I2PMultiaddr) Equal(multiaddr ma.Multiaddr) bool {
	if multiaddr.String() == addr.String() {
		return true
	}
	return false
}

func (addr I2PMultiaddr) ValueForProtocol(code int) (string, error) {
	if code == addr.Code {
		return string(addr.I2PAddr.String()), nil
	}
	return addr.baseMultiAddress.ValueForProtocol(code)
}

func NewI2PMultiaddr(inputs string) (I2PMultiaddr, error) {
	var m I2PMultiaddr
	var err error
    m.VCode = make([]byte, 2)
	if i := strings.SplitN(inputs, "/ntcp/", 2); len(i) == 2 {
		s := strings.Split(i[1], "/")
		m.I2PAddr, err = NewI2PAddrFromString(s[0])
		if err != nil {
			return m, err
		}
		m.baseMultiAddress = m.Decapsulate(m)
		m.Name = "ntcp"
		m.Code = P_GARLIC_NTCP
        binary.PutVarint(m.VCode, int64(m.Code))
		return m, nil
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
        binary.PutVarint(m.VCode, int64(m.Code))
		return m, fmt.Errorf("sam3-multiaddr Error: %s, %s", "ssu isn't implemented yet. Come back later.", s[0])
	}
	return m, fmt.Errorf("sam3-multiaddr Error: %s", "Not an i2p Multiaddr")
}
