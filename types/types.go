package types

import (
	"bytes"
	"encoding/gob"
	"net"
)

var (
	// Packet types
	PktCreateLobby  = "create_lobby"
	PktCreatedLobby = "created_lobby"
	PktJoinLobby    = "join_lobby"
	PktClientJoined = "client_joined"
	PktJoinedLobby  = "joined_lobby"
)

type Packet struct {
	Type string
	Data interface{}
}

type CreateLobbyPacketData struct {
	LobbyID     string
	PrivateAddr net.UDPAddr
}

type CreatedLobbyPacketData struct {
	LobbyID     string
	PublicAddr  net.UDPAddr
	PrivateAddr net.UDPAddr
}

type JoinLobbyPacketData struct {
	PublicAddr  net.UDPAddr
	PrivateAddr net.UDPAddr
	LobbyID     string
}

type ClientJoinedPacketData struct {
	ClientPrivateAddr net.UDPAddr
	ClientPublicAddr  net.UDPAddr
}

type JoinedLobbyPacketData struct {
	LobbyID         string
	HostPrivateAddr net.UDPAddr
	HostPublicAddr  net.UDPAddr
}

func RegisterTypes() {
	gob.Register(Packet{})
	gob.Register(CreateLobbyPacketData{})
	gob.Register(CreatedLobbyPacketData{})
	gob.Register(JoinLobbyPacketData{})
	gob.Register(ClientJoinedPacketData{})
	gob.Register(JoinedLobbyPacketData{})
}

func DecodePacket(data []byte) (*Packet, error) {
	var packet Packet
	err := gob.NewDecoder(bytes.NewBuffer(data)).Decode(&packet)
	if err != nil {
		return nil, err
	}
	return &packet, nil
}

func EncodePacket(packet *Packet) ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(packet)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
