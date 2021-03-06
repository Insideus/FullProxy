package Interface

import (
	"FullProxy/FullProxy/Proxies/Basic"
	"FullProxy/FullProxy/Sockets"
	"bufio"
	"log"
	"net"
	"os"
	"os/signal"
)


func setupControlCSignal(server net.Listener, masterConnection net.Conn){
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(serverConnection  net.Listener, clientConnection net.Conn){
		<-c
		_ = clientConnection.Close()
		_ = masterConnection.Close()
		os.Exit(0)
	}(server, masterConnection)
}


func startProxying(clientConnection net.Conn, targetConnection net.Conn){
	clientConnectionReader := bufio.NewReader(clientConnection)
	clientConnectionWriter := bufio.NewWriter(clientConnection)
	targetConnectionReader := bufio.NewReader(targetConnection)
	targetConnectionWriter := bufio.NewWriter(targetConnection)
	Basic.Proxy(
		clientConnection, targetConnection,
		clientConnectionReader, clientConnectionWriter,
		targetConnectionReader, targetConnectionWriter)
}


func Server(address string, port string){
	log.Print("Starting Interface server")
	server, BindingError  := net.Listen("tcp", address + ":" + port)
	if BindingError != nil {
		log.Print("Something goes wrong: " + BindingError.Error())
		return
	}
	log.Printf("Bind successfully in %s:%s", address, port)
	log.Print("Waiting for proxy server connections...")
	masterConnection, connectionError := server.Accept()
	if connectionError != nil{
		_ = server.Close()
		return
	}
	log.Print("Reverse connection received from: ", masterConnection.RemoteAddr())
	setupControlCSignal(server, masterConnection)
	masterConnectionWriter := bufio.NewWriter(masterConnection)
	for {
		clientConnection, _ := server.Accept()
		_, connectionError := Sockets.Send(masterConnectionWriter, []byte{1})
		if connectionError != nil{
			break
		}
		targetConnection, _ := server.Accept()
		go startProxying(clientConnection, targetConnection)
	}
	_ = masterConnection.Close()
	_ = server.Close()
}