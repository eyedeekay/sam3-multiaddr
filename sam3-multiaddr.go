package multiaddrsam

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/eyedeekay/sam3"
	ma "github.com/multiformats/go-multiaddr"
)

// I2PMultiaddr is an i2p-enabled multiaddr, with both types(TCP-like and UDP-like)
// capabilities
type I2PMultiaddr struct {
	Name    string
	Code    int
	VCode   []byte
	samhost string
	samport string
	bytes   []byte

	I2POnly          bool
	baseMultiAddress ma.Multiaddr
	sam3.I2PAddr
}

// These were picked arbitrarily. They will probably change.

// GarlicNTCP indicates an i2p streaming connection
const GarlicNTCP = 445

// GarlicSSU indicates an i2p datagram connection
const GarlicSSU = 890

// GarlicSAM indicates an i2p SAM connection
const GarlicSAM = 765

// GarlicVSAM stores the VCode of the connection
var GarlicVSAM []byte

//binary.PutVarint(m.VCode, int64(m.Code))

// Address converts the I2PMultiaddr into a sam3.I2PAddr which implements net.Addr
func (addr I2PMultiaddr) Address() *sam3.I2PAddr {
	return &addr.I2PAddr
}

// SAMAddress returns "/sam/SAMq Host:SAM Port" of the I2PMultiaddr as a string
func (addr I2PMultiaddr) SAMAddress() string {
	return "/sam/" + addr.samhost + ":" + addr.samport
}

// SAMAddressString returns the SAM address as "SAM host:SAM port"
func (addr I2PMultiaddr) SAMAddressString() string {
	return addr.samhost + ":" + addr.samport
}

// Bytes returns the whole address as a slice of bytes
func (addr I2PMultiaddr) Bytes() []byte {
	if addr.I2POnly {
		return []byte(addr.SAMAddress() + "/ntcp/" + addr.Address().String())
	}
	if addr.baseMultiAddress != nil {
		return []byte(addr.SAMAddress() + "/ntcp/" + addr.Address().String() + addr.baseMultiAddress.String())
	}
	return []byte(addr.SAMAddress() + "/ntcp/" + addr.Address().String())
}

// String returns the I2PMultiaddr as a string
func (addr I2PMultiaddr) String() string {
	return string(addr.Bytes())
}

// Encapsulate implements Multiaddr
func (addr I2PMultiaddr) Encapsulate(multiaddr ma.Multiaddr) ma.Multiaddr {
	mb := addr.Bytes()
	ob := multiaddr.Bytes()

	rb := make([]byte, len(mb)+len(ob))

	copy(rb, mb)
	copy(rb[len(mb):], ob)
	baddr, _ := NewI2PMultiaddr(string(rb), addr.I2POnly, addr.samhost+":"+addr.samport)
	return baddr
}

// Decapsulate implements Multiaddr
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

// Protocols implements Multiaddr
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
		Code:  GarlicSAM,
		Name:  "sam",
		VCode: GarlicVSAM,
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

// Equal determines if two I2PMultiaddr's are equal by comparing the strings.
func (addr I2PMultiaddr) Equal(multiaddr ma.Multiaddr) bool {
	if multiaddr.String() == addr.String() {
		return true
	}
	return false
}

// ValueForProtocol implements I2PMultiaddr
func (addr I2PMultiaddr) ValueForProtocol(code int) (string, error) {
	if code == addr.Code {
		return string(addr.I2PAddr.String()), nil
	}
	return addr.baseMultiAddress.ValueForProtocol(code)
}

// NewI2PMultiaddr creates a new i2p multiaddr, with a protocol string(ntcp, ssu)
// a switch for only being an i2p address, and an optional SAM Host/Port pair
func NewI2PMultiaddr(protocol string, i2ponly bool, samaddr ...string) (I2PMultiaddr, error) {
	var m I2PMultiaddr
	var err error
	m.VCode = make([]byte, 2)
	m.I2POnly = i2ponly
	GarlicVSAM = make([]byte, 2)
	binary.PutVarint(GarlicVSAM, int64(GarlicSAM))
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
		m.I2PAddr, err = sam3.NewI2PAddrFromString(s[0])
		if err != nil {
			return m, err
		}
		m.baseMultiAddress = m.Decapsulate(m)
		m.Name = "ntcp"
		m.Code = GarlicNTCP
		binary.PutVarint(m.VCode, int64(m.Code))
		m.bytes = m.Bytes()
		return m, nil
	}
	if i := strings.SplitN(protocol, "/ssu/", 2); len(i) == 2 {
		s := strings.Split(i[1], "/")
		m.I2PAddr, err = sam3.NewI2PAddrFromString(s[0])
		if err != nil {
			return m, err
		}
		m.baseMultiAddress = m.Decapsulate(m)
		m.Name = "ssu"
		m.Code = GarlicSSU
		binary.PutVarint(m.VCode, int64(m.Code))
		m.bytes = m.Bytes()
		return m, fmt.Errorf("sam3-multiaddr Error: %s, %s", "ssu isn't implemented yet. Come back later.", s[0])
	}

	return m, fmt.Errorf("sam3-multiaddr Error: %s", "Not an i2p Multiaddr")
}
