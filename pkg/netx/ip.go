package netx

import "net"

func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8")
	if err != nil {
		return ""
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.String()

}
