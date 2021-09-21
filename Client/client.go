package main

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"time"
)

func process(channel chan net.Conn, id, counter uint32) {
	i := counter
	for {
		select {
			case <-channel:
				c, err := net.Dial("tcp", ":9999")
				if err != nil {
					fmt.Println(err)
					return
				}
				msg := strconv.Itoa(int(id)) + "," + strconv.Itoa(int(i))
				err = gob.NewEncoder(c).Encode(msg)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("mando mensaje")
				return
			default:
				fmt.Printf("%d : %d\n", id, i)
		}
		i++
		time.Sleep(time.Millisecond * 500)
	}
}

func client(status string, channel chan net.Conn) {
	c, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	
	if status == "Get" {
		// pedimos al servidor un proceso
		err = gob.NewEncoder(c).Encode("Get")
		if err != nil {
			fmt.Println(err)
		}
		// recibimos respuesta del servidor
		response, err := ioutil.ReadAll(c)
		if err != nil {
			fmt.Println(err)
		}
		// ejecutamos el proceso
		// id y contador del proceso (servidor)
		id,_ := strconv.Atoi(string(response)[0:1])
		count,_ := strconv.Atoi(string(response)[2:])
		go process(channel, uint32(id), uint32(count))
	} else {
		// regresamos al servidor el proceso
		channel <- c
	}
	c.Close()
}

func main() {
	channel := make(chan net.Conn)
	go client("Get", channel)

	var input string
	fmt.Scanln(&input)
	go client("Post", channel)
	time.Sleep(time.Millisecond * 500)
}