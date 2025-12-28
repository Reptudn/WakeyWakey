package utils

import (
	"encoding/hex"
	"fmt"
	"net"
	"strings"
)

func SendWakeOnLANPacket(mac string) error {
	mac = strings.ReplaceAll(mac, ":", "")
	mac = strings.ReplaceAll(mac, "-", "")

	if len(mac) != 12 {
		return fmt.Errorf("invalid MAC address format")
	}

	macBytes, err := hex.DecodeString(mac)
	if err != nil {
		return fmt.Errorf("failed to decode MAC address: %v", err)
	}

	packet := make([]byte, 102)

	for i := 0; i < 6; i++ {
		packet[i] = 0xFF
	}

	for i := 1; i <= 16; i++ {
		copy(packet[i*6:], macBytes)
	}

	conn, err := net.Dial("udp", "255.255.255.255:9")
	if err != nil {
		return fmt.Errorf("failed to dial UDP: %v", err)
	}
	defer conn.Close()

	_, err = conn.Write(packet)
	if err != nil {
		return fmt.Errorf("failed to send packet: %v", err)
	}

	fmt.Println("Sending Wake-on-LAN packet to:", mac)
	return nil
}