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
	samhost          string
	samport          string
	bytes            []byte

    I2POnly          bool
	baseMultiAddress ma.Multiaddr
	I2PAddr
}

// These were picked arbitrarily. They will probably change.
var P_GARLIC_NTCP = 445
var P_GARLIC_SSU = 890

var P_GARLIC_SAM = 765

var P_GARLIC_VSAM []byte

//binary.PutVarint(m.VCode, int64(m.Code))

func (addr I2PMultiaddr) Address() *I2PAddr {
	return &addr.I2PAddr
}

func (addr I2PMultiaddr) SAMAddress() string {
	return "/sam/" + addr.samhost + ":" + addr.samport
}

//
func (addr I2PMultiaddr) Bytes() []byte {
    if addr.I2POnly {
        return []byte(addr.SAMAddress() + "/ntcp/" + addr.Address().String())
    }else{
        if addr.baseMultiAddress != nil {
            return []byte(addr.SAMAddress() + "/ntcp/" + addr.Address().String() + addr.baseMultiAddress.String())
        }else{
            return []byte(addr.SAMAddress() + "/ntcp/" + addr.Address().String())
        }
    }
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
	baddr, _ := NewI2PMultiaddr(string(rb), addr.I2POnly, addr.samhost + ":" + addr.samport )
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
		Code:  addr.Code,
		Name:  addr.Name,
		VCode: addr.VCode,
		Size:  0,
		Path:  false,
	})
	p = append(p, ma.Protocol{
		Code:  P_GARLIC_SAM,
		Name:  "sam",
		VCode: P_GARLIC_VSAM,
		Size:  0,
		Path:  false,
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

func NewI2PMultiaddr(protocol string, i2ponly bool, samaddr ...string) (I2PMultiaddr, error) {
	var m I2PMultiaddr
	var err error
	m.VCode = make([]byte, 2)
	P_GARLIC_VSAM = make([]byte, 2)
	binary.PutVarint(P_GARLIC_VSAM, int64(P_GARLIC_SAM))
	if len(samaddr) == 1 {
		if i := strings.SplitN(samaddr[0], "/sam/", 2); len(i) == 2 {
			if j := strings.Split(i[1], ":"); len(j) >= 2 && len(j) <= 3 {
				m.samhost = j[0]
				m.samport = j[1]
			}
		}
	} else if len(samaddr) == 0 {
		m.samhost = "127.0.0.1"
		m.samport = "7656"
	} else {
		return m, fmt.Errorf("SAM address passed to multiaddr invalid %s", samaddr[0])
	}
	if i := strings.SplitN(protocol, "/ntcp/", 2); len(i) == 2 {
		s := strings.Split(i[1], "/")
		m.I2PAddr, err = NewI2PAddrFromString(s[0])
		if err != nil {
			return m, err
		}
		m.baseMultiAddress = m.Decapsulate(m)
		m.Name = "ntcp"
		m.Code = P_GARLIC_NTCP
		binary.PutVarint(m.VCode, int64(m.Code))
		m.bytes = m.Bytes()
		return m, nil
	}
	if i := strings.SplitN(protocol, "/ssu/", 2); len(i) == 2 {
		s := strings.Split(i[1], "/")
		m.I2PAddr, err = NewI2PAddrFromString(s[0])
		if err != nil {
			return m, err
		}
		m.baseMultiAddress = m.Decapsulate(m)
		m.Name = "ssu"
		m.Code = P_GARLIC_SSU
		binary.PutVarint(m.VCode, int64(m.Code))
		m.bytes = m.Bytes()
		return m, fmt.Errorf("sam3-multiaddr Error: %s, %s", "ssu isn't implemented yet. Come back later.", s[0])
	}

	return m, fmt.Errorf("sam3-multiaddr Error: %s", "Not an i2p Multiaddr")
}
