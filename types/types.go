package types

import "encoding/gob"

type Packet struct {
	Type string
	Data interface{}
}

type CreateLobbyPacketData struct {
	LobbyID     string
	PrivateAddr string
}

type CreatedLobbyPacketData struct {
	LobbyID string
}

type JoinLobbyPacketData struct {
	PrivateAddr string
	LobbyID     string
}

type ClientJoinedPacketData struct {
	ClientPrivateAddr string
	ClientPublicAddr  string
}

type JoinedLobbyPacketData struct {
	LobbyID         string
	HostPrivateAddr string
	HostPublicAddr  string
}

func RegisterTypes() {
	gob.Register(Packet{})
	gob.Register(CreateLobbyPacketData{})
	gob.Register(CreatedLobbyPacketData{})
	gob.Register(JoinLobbyPacketData{})
	gob.Register(ClientJoinedPacketData{})
	gob.Register(JoinedLobbyPacketData{})
}
