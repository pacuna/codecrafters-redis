package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var data = make(map[string]string)

func main() {

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go func(c net.Conn) {
			defer c.Close()
			for {
				buf := make([]byte, 256)
				n, err := c.Read(buf)
				if err == io.EOF {
					break
				}
				buf = buf[:n]
				ary := ParseArray(buf)
				cmd := strings.ToLower(ary.Items[0].Val)
				if cmd == "ping" {
					fmt.Fprintf(c, "+PONG\r\n")
				} else if cmd == "echo" {
					fmt.Fprintf(c, "+%s\r\n", ary.Items[1].Val)
				} else if cmd == "set" {
					if len(ary.Items) <= 3 {
						data[ary.Items[1].Val] = ary.Items[2].Val
					} else {
						if strings.ToLower(ary.Items[3].Val) == "ex" {
							exVal, _ := strconv.Atoi(ary.Items[4].Val)
							data[ary.Items[1].Val] = ary.Items[2].Val
							time.AfterFunc(time.Second*time.Duration(exVal), func() {
								delete(data, ary.Items[1].Val)
							})
						} else if strings.ToLower(ary.Items[3].Val) == "px" {
							exVal, _ := strconv.Atoi(ary.Items[4].Val)
							data[ary.Items[1].Val] = ary.Items[2].Val
							time.AfterFunc(time.Millisecond*time.Duration(exVal), func() {
								delete(data, ary.Items[1].Val)
							})
						}
					}
					fmt.Fprintf(c, "+OK\r\n")
				} else if cmd == "get" {
					if res, ok := data[ary.Items[1].Val]; ok {
						fmt.Fprintf(c, "$%d\r\n%s\r\n", len(res), res)
					} else {
						fmt.Fprintf(c, "$%d\r\n", -1)
					}
				} else {
					fmt.Fprintf(c, "+PONG\r\n")
				}

			}
		}(conn)
	}
}

type Item struct {
	Type string
	Val  string
}

type Array struct {
	Items []*Item
}

//*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n
func ParseArray(input []byte) *Array {
	var i int

	// skip *
	i += 1

	// read array len
	aryLen := 0
	for ; input[i] != '\r'; i++ {
		aryLen = aryLen*10 + int(input[i]-'0')
	}

	ary := &Array{
		Items: make([]*Item, aryLen),
	}

	// read items
	for j := 0; j < aryLen; j++ {
		itemLen := 0
		for ; input[i] != '$'; i++ {
		}
		// skip $
		i += 1
		for ; input[i] != '\r'; i++ {
			itemLen = itemLen*10 + int(input[i]-'0')
		}
		// skip \r and \n
		i += 2
		item := &Item{
			Type: "bulkstring",
			Val:  string(input[i : i+itemLen]),
		}
		ary.Items[j] = item
	}

	return ary
}
