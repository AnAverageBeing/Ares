package mcutils

import (
	"Ares/net/minecraft/packet"
)

type nextState int

const (
	Status nextState = 1
	Login  nextState = 2
)

func GetHandshakePacket(ip string, port int, protocol int, nextState nextState) packet.Packet {
	pk := packet.Marshal(
		0x00,
		packet.VarInt(protocol),
		packet.String(ip),
		packet.UnsignedShort(port),
		packet.VarInt(nextState),
	)
	return pk
}

func GetLoginPacket(name string, versionProtocol int) packet.Packet {
	if versionProtocol == 760 || versionProtocol == 759 {
		pk := packet.Marshal(
			0x00,
			packet.String(name),
			packet.Boolean(false),
			packet.Boolean(false),
		)
		return pk
	}

	if versionProtocol == 761 {
		pk := packet.Marshal(
			0x00,
			packet.String(name),
			packet.Boolean(false),
		)
		return pk
	}

	pk := packet.Marshal(
		0x00,
		packet.String(name),
	)

	return pk
}
