package main

import (
    "encoding/binary"
    "io"
    "log"
    "net"
    "os"
)

type FPMHeader struct {
    Version     uint8
    MessageType uint8
    MessageLen  uint16
}

func handleConnection(conn net.Conn) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Recovered from panic: %v. Closing connection.", r)
            conn.Close() // Close the current connection
        }
    }()
    
    for {
        h := FPMHeader{}
        err := binary.Read(conn, binary.BigEndian, &h)
        if err != nil {
            log.Println("Failed to read header:", err)
            return
        }

        if h.Version != 1 {
            panic("Unsupported FPM frame version")
        }

        if h.MessageType != 1 {
            panic("Unsupported FPM frame type")
        }

        n, err := io.CopyN(os.Stdout, conn, int64(h.MessageLen-4))
        if err != nil {
            panic(err)
        }

        if n != int64(h.MessageLen-4) {
            panic("Couldn't read entire message")
        }
    }
}

func main() {
    ln, err := net.Listen("tcp", ":2620")
    if err != nil {
        panic(err)
    }
    defer ln.Close()

    log.Println("Listening on :2620")
    for {
        conn, err := ln.Accept()
        if err != nil {
            log.Println("Failed to accept connection:", err)
            continue
        }
        go handleConnection(conn) // handle each connection in a new goroutine
    }
}
