package mcutils

import (
	"Ares/net/minecraft/packet"
	"fmt"
)

type nextState int

const (
	Status nextState = 1
	Login  nextState = 2
)

func GetHandshakePacket(ip string, port int, protocol int, nextState nextState) (pk packet.Packet) {
	pk = packet.Marshal(
		0x00,
		packet.VarInt(protocol),
		packet.String(ip),
		packet.UnsignedShort(port),
		packet.VarInt(nextState),
	)
	return
}

func GetLoginPacket(name string, versionProtocol int) (pk packet.Packet) {

	if versionProtocol == 760 || versionProtocol == 759 {
		pk = packet.Marshal(
			0x00,
			packet.String(name),
			packet.Boolean(false),
			packet.Boolean(false),
		)
		fmt.Println("WOW")
		return
	}

	if versionProtocol == 761 {
		pk = packet.Marshal(
			0x00,
			packet.String(name),
			packet.Boolean(false),
			packet.UUID{},
		)
		fmt.Println("HMM")
		return
	}

	pk = packet.Marshal(
		0x00,
		packet.String(name),
	)
	fmt.Println("XD")

	return
}
