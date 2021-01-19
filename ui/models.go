package ui

import (
	"encoding/json"
)

type MessageType string

const (
	MessageTypeStartTunnelHTTP      MessageType = "MT_RequestTunnelStart_HTTP"
	MessageTypeStartTunnelDirectory MessageType = "MT_RequestTunnelStart_Directory"
	MessageTypeStartTunnelWebDav    MessageType = "MT_RequestTunnelStart_WebDav"
	MessageTypeStopTunnel           MessageType = "MT_RequestTunnelStop"
	MessageTypeAuthorization        MessageType = "MT_RequestLogin"
	MessageTypeLogout               MessageType = "MT_RequestLogout"
	MessageTypeOpenBrowser          MessageType = "MT_OpenInBrowser"
)

type Message struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type StopTunnelMessage struct {
	TunnelID string `json:"tunnelId"`
}

type OpenInBrowserMessage struct {
	URL string `json:"url"`
}
