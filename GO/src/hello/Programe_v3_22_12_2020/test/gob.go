package main

import (
    "encoding/gob"
    "fmt"
    "net"
    "strconv"
    "time"
)

type Message struct {
    Msg string
}

func main() {
    gob.Register(new(Message))

    clientAddr, err := net.ResolveTCPAddr("tcp", "localhost:12346")
    if err != nil {
        fmt.Println(err)
    }
    serverAddr, err := net.ResolveTCPAddr("tcp", "localhost:12345")
    if err != nil {
        fmt.Println(err)
    }

    serverListener, err := net.ListenTCP("tcp", serverAddr)
    if err != nil {
        fmt.Println(err)
    }
    conn, err := net.DialTCP("tcp", clientAddr, serverAddr)
    if err != nil {
        fmt.Println(err)
    }
    serverConn, err := serverListener.AcceptTCP()
    if err != nil {
        fmt.Println(err)
    }
	dec := gob.NewDecoder(conn) // Will read from network.
    enc := gob.NewEncoder(serverConn)
    
    done := false
    go func() {
        for !done {
            recieveMessage(dec)
        }
    }()

    for i := 1; i < 1000; i++ {
        sent := Message{strconv.Itoa(i)}
        sendMessage(sent, enc)
    }
    time.Sleep(time.Second)
    done = true
}

func sendMessage(msg Message, enc *gob.Encoder) {
    err := enc.Encode(msg)
    if err != nil {
        fmt.Println(err)
    }
}

func recieveMessage(dec *gob.Decoder) {
    msg := new(Message)
    err := dec.Decode(msg)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("Client recieved:", msg.Msg)
}
