package main

import (
	"bytes"
	"log"
	"net"
	"testing"
	"time"
)

func TestMain_Echo(t *testing.T) {
	conn, err := net.Dial("tcp", "0.0.0.0:6379")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	_, err = conn.Write([]byte("*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n"))
	if err != nil {
		log.Fatal(err)
	}
	want := []byte("+hey\r\n")
	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	buf = buf[:n]
	if bytes.Compare(want, buf) != 0 {
		t.Errorf("want %v, got %v instead", want, buf)
	}
}

func TestMain_SetGet_NewVal(t *testing.T) {
	conn, err := net.Dial("tcp", "0.0.0.0:6379")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	_, err = conn.Write([]byte("*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"))
	if err != nil {
		log.Fatal(err)
	}
	want := []byte("+OK\r\n")
	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	buf = buf[:n]
	if bytes.Compare(want, buf) != 0 {
		t.Errorf("want %v, got %v instead", want, buf)
	}

	_, err = conn.Write([]byte("*2\r\n$3\r\nGET\r\n$3\r\nfoo\r\n"))
	if err != nil {
		log.Fatal(err)
	}
	want = []byte("$3\r\nbar\r\n")
	buf = make([]byte, 256)
	n, err = conn.Read(buf)
	buf = buf[:n]
	if bytes.Compare(want, buf) != 0 {
		t.Errorf("want %v, got %v instead", want, buf)
	}
}

func TestMain_SetGet_NewVal_EX(t *testing.T) {
	conn, err := net.Dial("tcp", "0.0.0.0:6379")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	_, err = conn.Write([]byte("*5\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n$2\r\nEX\r\n$1\r\n3\r\n"))
	if err != nil {
		log.Fatal(err)
	}
	want := []byte("+OK\r\n")
	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	buf = buf[:n]
	if bytes.Compare(want, buf) != 0 {
		t.Errorf("want %v, got %v instead", want, buf)
	}

	_, err = conn.Write([]byte("*2\r\n$3\r\nGET\r\n$3\r\nfoo\r\n"))
	if err != nil {
		log.Fatal(err)
	}
	want = []byte("$3\r\nbar\r\n")
	buf = make([]byte, 256)
	n, err = conn.Read(buf)
	buf = buf[:n]
	if bytes.Compare(want, buf) != 0 {
		t.Errorf("want %v, got %v instead", want, buf)
	}

	//after 3 seconds, key should be gone
	time.Sleep(time.Duration(3) * time.Second)
	_, err = conn.Write([]byte("*2\r\n$3\r\nGET\r\n$3\r\nfoo\r\n"))
	if err != nil {
		log.Fatal(err)
	}
	want = []byte("$-1\r\n")
	buf = make([]byte, 256)
	n, err = conn.Read(buf)
	buf = buf[:n]
	if bytes.Compare(want, buf) != 0 {
		t.Errorf("want %v, got %v instead", want, buf)
	}
}

func TestMain_SetGet_NewVal_PX(t *testing.T) {
	conn, err := net.Dial("tcp", "0.0.0.0:6379")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	_, err = conn.Write([]byte("*5\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n$2\r\nPX\r\n$3\r\n300\r\n"))
	if err != nil {
		log.Fatal(err)
	}
	want := []byte("+OK\r\n")
	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	buf = buf[:n]
	if bytes.Compare(want, buf) != 0 {
		t.Errorf("want %v, got %v instead", want, buf)
	}

	_, err = conn.Write([]byte("*2\r\n$3\r\nGET\r\n$3\r\nfoo\r\n"))
	if err != nil {
		log.Fatal(err)
	}
	want = []byte("$3\r\nbar\r\n")
	buf = make([]byte, 256)
	n, err = conn.Read(buf)
	buf = buf[:n]
	if bytes.Compare(want, buf) != 0 {
		t.Errorf("want %v, got %v instead", want, buf)
	}

	//after 300 ms, key should be gone
	time.Sleep(time.Duration(300) * time.Millisecond)
	_, err = conn.Write([]byte("*2\r\n$3\r\nGET\r\n$3\r\nfoo\r\n"))
	if err != nil {
		log.Fatal(err)
	}
	want = []byte("$-1\r\n")
	buf = make([]byte, 256)
	n, err = conn.Read(buf)
	buf = buf[:n]
	if bytes.Compare(want, buf) != 0 {
		t.Errorf("want %v, got %v instead", want, buf)
	}
}
