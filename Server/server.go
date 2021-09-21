package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

func process(channel chan net.Conn, id, counter uint32) {
	i := counter
	for {
		select {
			case c := <-channel:
				io.WriteString(c, fmt.Sprintf("%d,%d",id ,i))
				return
			default:
				fmt.Printf("%d : %d\n", id, i)
		}
		i++
		time.Sleep(time.Millisecond * 500)
	}
}

func server() {
	// se crea el servidor
	s, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer s.Close()

	// administrador de procesos
	channel := make(chan net.Conn)
	procs := []uint32{0,1,2,3,4}
	for _,v := range(procs) {
		go process(channel, v, 0)
	}

	for {
		// peticiones del cliente
		c, err := s.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		handleClient(c, channel)
		
		c.Close()
	}
}

func handleClient(c net.Conn, channel chan net.Conn) {
	var msg string
	err := gob.NewDecoder(c).Decode(&msg)
	if err != nil {
		return
	} else {
		if msg == "Get" {
			fmt.Println("Proceso enviado")
			channel <- c
		} else {
			// id y contador del proceso (cliente)
			id,_ := strconv.Atoi(msg[0:1])
			count,_ := strconv.Atoi(msg[2:])
			fmt.Printf("Proceso (%d,%d) retornado\n", id, count)
			go process(channel, uint32(id), uint32(count))
		}
	}
}

func main() {
	go server()
	var input string
	fmt.Scanln(&input)
}