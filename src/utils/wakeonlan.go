package utils

import (
	"encoding/hex"
	"fmt"
	"net"
	"os/exec"
	"strings"
)

func SendWakeOnLANPacket(mac string) error {

	if (!IsValidMacAddress(mac)) {
		return fmt.Errorf("Invalid MAC address format")
	}

	mac = strings.ReplaceAll(mac, ":", "")
	mac = strings.ReplaceAll(mac, "-", "")

	macBytes, err := hex.DecodeString(mac)
	if err != nil {
		return fmt.Errorf("Failed to decode MAC address: %v", err)
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
		return fmt.Errorf("Failed to dial UDP: %v", err)
	}
	defer conn.Close()

	_, err = conn.Write(packet)
	if err != nil {
		return fmt.Errorf("Failed to send packet: %v", err)
	}

	fmt.Println("Sending Wake-on-LAN packet to:", mac)
	return nil
}

func SendWakeOnLANPacketViaCommand(mac string) error {
	cmd := exec.Command("awake", mac)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Failed to execute wakeonlan command: %v, output: %s", err, output)
	}
	
	fmt.Println("Wake-on-LAN packet sent via command:", string(output))
	return nil

}