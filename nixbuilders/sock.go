package nixbuilders

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

func reader(r io.Reader) {
	buf := make([]byte, 8)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			return
		}
		println("Client got:", string(buf[0:n]))
	}
}

func connectSocket() {
	socketPath := "/nix/var/nix/daemon-socket/socket"
	c, err := net.Dial("unix", socketPath)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	go reader(c)

	wr := bufio.NewWriter(c)

	buf := make([]byte, binary.MaxVarintLen64)
	magic := uint64(0x6e697863)
	n := binary.PutUvarint(buf, magic)
	buf = buf[:n]
	fmt.Println(wr.Write(buf))
	fmt.Println(wr.WriteString("\n"))
	wr.Flush()

	// for {
	// 	_, err := c.Write([]byte("hi"))
	// 	if err != nil {
	// 		log.Fatal("write error:", err)
	// 		break
	// 	}
	// 	time.Sleep(1e9)
	// }
	<-time.After(time.Hour * 24)
}
