package SOCKS5

import (
	"bufio"
	"net"
)


func HandleCommandExecution(
	clientConnection net.Conn,
	clientConnectionReader *bufio.Reader, clientConnectionWriter *bufio.Writer,
	targetRequestedCommand *byte, targetAddressType *byte,
	targetAddress *string, targetPort *string,
	rawTargetAddress []byte, rawTargetPort []byte){


	switch *targetRequestedCommand {
	case Connect:
		PrepareConnect(clientConnection, clientConnectionReader, clientConnectionWriter, targetAddress, targetPort, rawTargetAddress, rawTargetPort, targetAddressType)
	case Bind:
		break
	case UDPAssociate:
		break
	}
}
