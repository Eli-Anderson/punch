package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"log/slog"
	"net"
	"time"

	"github.com/Eli-Anderson/punch/types"
)

type Client struct {
	PrivateAddr *net.UDPAddr
	PublicAddr  *net.UDPAddr
}

type Lobby struct {
	ID        string
	CreatedAt time.Time
	Host      *Client
}

func main() {
	types.RegisterTypes()

	addr := net.UDPAddr{
		Port: 8003,
		IP:   net.ParseIP("127.0.0.1"),
	}
	server, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatal(err)
		return
	}

	lobbies := []*Lobby{}
	slog.Info("Listening for connections", "addr", addr.String())

	for {
		p := make([]byte, 2048)
		_, addr, err := server.ReadFromUDP(p)
		if err != nil {
			log.Fatal(err)
		}

		packet := types.Packet{}
		gob.NewDecoder(bytes.NewBuffer(p)).Decode(&packet)

		switch packet.Type {
		case types.PktCreateLobby:
			data := packet.Data.(types.CreateLobbyPacketData)
			slog.Info("Creating lobby", "addr", addr.String())
			client := &Client{
				PrivateAddr: &data.PrivateAddr,
				PublicAddr:  addr,
			}
			lobby := &Lobby{
				ID:        data.LobbyID,
				CreatedAt: time.Now(),
				Host:      client,
			}
			lobbies = append(lobbies, lobby)
			packet := types.Packet{
				Type: types.PktCreatedLobby,
				Data: types.CreatedLobbyPacketData{
					LobbyID:     lobby.ID,
					PublicAddr:  *addr,
					PrivateAddr: data.PrivateAddr,
				},
			}
			var buf bytes.Buffer
			gob.NewEncoder(&buf).Encode(packet)
			server.WriteToUDP(buf.Bytes(), addr)
		case types.PktJoinLobby:
			data := packet.Data.(types.JoinLobbyPacketData)
			slog.Info("Client joining lobby", "addr", addr.String(), "lobbyID", data.LobbyID)
			var lobby *Lobby
			for _, m := range lobbies {
				if m.ID == data.LobbyID {
					lobby = m
					break
				}
			}
			if lobby == nil {
				slog.Info("Client sent invalid lobby ID", "addr", addr.String())
				continue
			}
			packet := types.Packet{
				Type: types.PktClientJoined,
				Data: types.ClientJoinedPacketData{
					ClientPrivateAddr: data.PrivateAddr,
					ClientPublicAddr:  *addr,
				},
			}
			var buf bytes.Buffer
			gob.NewEncoder(&buf).Encode(packet)
			server.WriteToUDP(buf.Bytes(), lobby.Host.PublicAddr)

			packet = types.Packet{
				Type: types.PktJoinedLobby,
				Data: types.JoinedLobbyPacketData{
					LobbyID:         lobby.ID,
					HostPrivateAddr: *lobby.Host.PrivateAddr,
					HostPublicAddr:  *lobby.Host.PublicAddr,
				},
			}
			buf.Reset()
			gob.NewEncoder(&buf).Encode(packet)
			server.WriteToUDP(buf.Bytes(), addr)
		}
	}
}
