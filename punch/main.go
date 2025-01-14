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
		case "create_lobby":
			data := packet.Data.(types.CreateLobbyPacketData)
			slog.Info("Creating lobby", "addr", addr.String())
			privateAddr, err := net.ResolveUDPAddr("udp", data.PrivateAddr)
			if err != nil {
				slog.Info("Client sent invalid private address", "addr", addr.String())
				continue
			}
			client := &Client{
				PrivateAddr: privateAddr,
				PublicAddr:  addr,
			}
			lobby := &Lobby{
				ID:        data.LobbyID,
				CreatedAt: time.Now(),
				Host:      client,
			}
			lobbies = append(lobbies, lobby)
			packet := types.Packet{
				Type: "created_lobby",
				Data: types.CreatedLobbyPacketData{
					LobbyID: lobby.ID,
				},
			}
			var buf bytes.Buffer
			gob.NewEncoder(&buf).Encode(packet)
			server.WriteToUDP(buf.Bytes(), addr)
		case "join_lobby":
			data := packet.Data.(types.JoinLobbyPacketData)
			slog.Info("Client joining lobby", "addr", addr.String(), "lobbyID", data.LobbyID)
			privateAddr, err := net.ResolveUDPAddr("udp", data.PrivateAddr)
			if err != nil {
				slog.Info("Client sent invalid private address", "addr", addr.String())
				continue
			}
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
				Type: "client_joined",
				Data: types.ClientJoinedPacketData{
					ClientPrivateAddr: privateAddr.String(),
					ClientPublicAddr:  addr.String(),
				},
			}
			var buf bytes.Buffer
			gob.NewEncoder(&buf).Encode(packet)
			server.WriteToUDP(buf.Bytes(), lobby.Host.PublicAddr)

			packet = types.Packet{
				Type: "joined_lobby",
				Data: types.JoinedLobbyPacketData{
					LobbyID:         lobby.ID,
					HostPrivateAddr: lobby.Host.PrivateAddr.String(),
					HostPublicAddr:  lobby.Host.PublicAddr.String(),
				},
			}
			buf.Reset()
			gob.NewEncoder(&buf).Encode(packet)
			server.WriteToUDP(buf.Bytes(), addr)
		}
	}
}
