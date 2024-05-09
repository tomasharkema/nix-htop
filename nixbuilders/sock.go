package nixbuilders

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/nix-community/go-nix/pkg/wire"
)

func reader(r io.Reader) {
	for {
		rd := wire.NewBytesReader(r, 64)
		buf := make([]byte, 64)
		for {
			n, err := rd.Read(buf[:])
			if err != nil {
				log.Fatalln("reader", err)
				return
			}
			println("Client got:", string(buf[0:n]))
		}
	}
}

func ConnectSocket() {
	socketPath := "/nix/var/nix/daemon-socket/socket"
	c, err := net.Dial("unix", socketPath)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	// wr, err := wire.NewBytesWriter(c, 64)
	// if err != nil {
	// 	panic(err)
	// }

	// go reader(c)
	fmt.Println(wire.WriteUint64(c, 0x6e697863))
	magix2, _ := wire.ReadUint64(c)
	fmt.Println(fmt.Sprintf("%02X", magix2), "0x6478696f", magix2 == 0x6478696f)

	fmt.Println(wire.WriteUint64(c, 0x10b))
	// for {
	// 	_, err := c.Write([]byte("hi"))
	// 	if err != nil {
	// 		log.Fatal("write error:", err)
	// 		break
	// 	}
	// 	time.Sleep(1e9)
	// }
}
