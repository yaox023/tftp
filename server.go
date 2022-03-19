package tftp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
)

func Serve() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	conn, err := net.ListenPacket("udp", "127.0.0.1:69")
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, DatagramSize)
	for {
		n, clientAddr, err := conn.ReadFrom(buf)
		if err != nil {
			log.Println(err)
			continue
		}
		req := RequestPacket{}
		err = req.Unmarshal(buf[:n])
		if err != nil {
			log.Println(err)
			continue
		}
		if req.Opcode == OpcodeReadRequest {
			go handleReadRequest(clientAddr, req)
		} else {
			go handleWriteRequest(clientAddr, req)
		}
	}
}

func handleReadRequest(addr net.Addr, req RequestPacket) {
	log.Printf("get read request, addr: %s, filename: %s\n", addr.String(), req.FileName)

	conn, err := net.Dial("udp", addr.String())
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	if req.Mode != ModeOctet {
		ep := ErrorPacket{ErrorCodeInvalidMode, "only mode octet is supported"}
		epBytes, err := ep.Marshal()
		if err != nil {
			log.Println(err)
			return
		}
		_, err = conn.Write(epBytes)
		if err != nil {
			log.Println(err)
			return
		}
		return
	}

	file, err := os.Open(req.FileName)
	if err != nil {
		ep := ErrorPacket{ErrorCodeFileNotFound, fmt.Sprintf("file not found, %s\n", err.Error())}
		epBytes, err := ep.Marshal()
		if err != nil {
			log.Println(err)
			return
		}
		_, err = conn.Write(epBytes)
		if err != nil {
			log.Println(err)
			return
		}
	}
	defer file.Close()

	dp := DataPacket{Payload: file}
	ep := ErrorPacket{}
	ap := AckPacket{}

	buf := make([]byte, DatagramSize)
	for n := DatagramSize; n == DatagramSize; {

		// send data packet
		dpBytes, err := dp.Marshal()
		if err != nil {
			log.Println(err)
			return
		}
		n, err = conn.Write(dpBytes)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("sent data packet with block: %d, size: %d\n", dp.Block, len(dpBytes))

		// wait for ack packet
		rn, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}

		packet := buf[:rn]
		switch {
		case ap.Unmarshal(packet) == nil:
			if ap.Block != dp.Block {
				log.Printf("error block, expect: %d, got: %d\n", dp.Block, ap.Block)
				return
			}
			log.Printf("got ack packet with block %d\n", ap.Block)
		case ep.Unmarshal(packet) == nil:
			log.Println(ep.Code, ep.Msg)
			return
		default:
			log.Println("invalid packet, ", buf[:n])
			return
		}

	}
}

func handleWriteRequest(addr net.Addr, req RequestPacket) {
	log.Printf("get write request, addr: %s, filename: %s\n", addr.String(), req.FileName)

	conn, err := net.Dial("udp", addr.String())
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	if req.Mode != ModeOctet {
		ep := ErrorPacket{ErrorCodeInvalidMode, "only mode octet is supported"}
		epBytes, err := ep.Marshal()
		if err != nil {
			log.Println(err)
			return
		}
		_, err = conn.Write(epBytes)
		if err != nil {
			log.Println(err)
			return
		}
		return
	}

	_, err = os.Stat(req.FileName)

	if err != nil && !os.IsNotExist(err) {
		ep := ErrorPacket{ErrorCodeFileExists, "file already exists"}
		epBytes, err := ep.Marshal()
		if err != nil {
			log.Println(err)
			return
		}
		_, err = conn.Write(epBytes)
		if err != nil {
			log.Println(err)
			return
		}
		return
	}

	file, err := os.Create(req.FileName)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	ap := AckPacket{0}
	apBytes, err := ap.Marshal()
	if err != nil {
		log.Println(err)
		return
	}
	_, err = conn.Write(apBytes)
	if err != nil {
		log.Println(err)
		return
	}

	buf := make([]byte, DatagramSize)
	dp := DataPacket{Payload: new(bytes.Buffer)}
	ep := ErrorPacket{}
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}
		packet := buf[:n]

		blockSize, err := dp.Unmarshal(packet)
		if err == nil {
			log.Printf("got data packet with block %d, blockSize %d\n", dp.Block, blockSize)

			ap.Block = dp.Block
			apBytes, err = ap.Marshal()
			if err != nil {
				log.Println(err)
				return
			}
			_, err = conn.Write(apBytes)
			if err != nil {
				log.Println(err)
				return
			}
			log.Printf("sent ack packet with block %d", ap.Block)
			if blockSize != BlockSize {
				bs, err := ioutil.ReadAll(dp.Payload)
				if err != nil {
					log.Println(err)
					return
				}
				n, err = file.Write(bs)
				if err != nil {
					log.Println(err)
					return
				}
				log.Printf("write file success, filename: %s, filesize: %d\n", req.FileName, n)
				return
			}
		} else {
			err = ep.Unmarshal(packet)
			if err != nil {
				log.Println("invalid packet, ", packet, err)
				return
			}
			log.Println(ep.Code, ep.Msg)
		}
	}
}
