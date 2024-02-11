package net

import (
	"bytes"
	"compress/zlib"
	"encoding/hex"
	"fmt"
	"net"
)

type ConnectionCommandType int

const (
	SendData ConnectionCommandType = iota
	StartCompression
	StopCompression
	CloseConnection
)

type ConnectionCommand struct {
	command ConnectionCommandType
	data    []byte
}

type ConnectionRead struct {
	length int
	data   []byte
}

type ZLibContext struct {
	zLibWriter *zlib.Writer
	payload    *bytes.Buffer
}

func NewZLibContext() ZLibContext {
	buf := &bytes.Buffer{}
	return ZLibContext{
		zLibWriter: zlib.NewWriter(buf),
		payload:    buf,
	}
}

func (z ZLibContext) Compress(data []byte) []byte {
	_, err := z.zLibWriter.Write(data)
	if err != nil {
		fmt.Println("Error compressing data:", err)
		return nil
	}
	err = z.zLibWriter.Flush()
	if err != nil {
		fmt.Println("Error flushing data:", err)
		return nil
	}
	result := z.payload.Bytes()
	z.payload.Reset()
	return result
}

type Connection struct {
	conn        net.Conn
	command     chan ConnectionCommand
	read        chan ConnectionRead
	compressed  bool
	zlibContext ZLibContext
}

func (c Connection) HandleConnection() {
	go c.handleAsyncWrite()
	go c.handleAsyncRead()
}

func (c Connection) handleAsyncWrite() {
	for {
		select {
		case command := <-c.command:
			switch command.command {
			case SendData:
				c.sendData(command.data)
			case StartCompression:
				c.zlibContext = NewZLibContext()
				c.compressed = true
			case StopCompression:
				c.compressed = false
			case CloseConnection:
				err := c.conn.Close()
				if err != nil {
					fmt.Println("Error closing connection:", err)
				}
				return
			}
		}
	}
}

func (c Connection) handleAsyncRead() {
	for {
		buf := make([]byte, 1024)
		length, err := c.conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}
		fmt.Println("Data received: " + hex.EncodeToString(buf[:length]))
		c.read <- ConnectionRead{length: length, data: buf[:length]}
	}
}

func (c Connection) sendData(data []byte) {
	fmt.Println("Data to send: " + hex.EncodeToString(data))
	var payload bytes.Buffer

	if c.compressed == true {
		compressedData := c.zlibContext.Compress(data)
		payload.Write(compressedData)
	} else {
		payload.Write(data)
	}

	fmt.Println("Data sent: " + hex.EncodeToString(payload.Bytes()))
	fmt.Println(hex.EncodeToString(payload.Bytes()))
	_, err := c.conn.Write(payload.Bytes())
	if err != nil {
		fmt.Println("Error writing data:", err)
		return
	}
}

func NewConnection(conn net.Conn, command chan ConnectionCommand, read chan ConnectionRead) *Connection {
	return &Connection{conn: conn, command: command, compressed: false, read: read}
}
