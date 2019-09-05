// TcpProxy project main.go
package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("missing message!")
		return
	}
	ip := os.Args[1]
	port, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("error happened ,exit")
		return
	}
	addr := os.Args[3]
	host := "Host: " + addr

	Service(ip, port, addr, host)
}

func Service(ip string, port int, dstaddr string, dsthost string) {
	// listen and accept
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(ip), port, ""})

	if err != nil {
		fmt.Println("listen error: ", err.Error())
		return
	}
	fmt.Println("init done...")

	for {
		client, err := listen.AcceptTCP()
		if err != nil {
			fmt.Println("accept error: ", err.Error())
			continue
		}
		go Channal(client, dstaddr, dsthost)   //build a channal for every coming TCP
	}
}

func Channal(client *net.TCPConn, addr string, rhost string) {

	tcpAddr, _ := net.ResolveTCPAddr("tcp4", addr)   //get the source ip
	conn, err := net.DialTCP("tcp", nil, tcpAddr)    //dial a tcp to dest
	if err != nil {
		fmt.Println("connection error: ", err.Error())
		client.Close()
		return
	}

	go ReadRequest(client, conn, rhost)   // client is from source, conn is to dest
	ReadResponse(conn, client)
}

func ReadRequest(lconn *net.TCPConn, rconn *net.TCPConn, dsthost string) {
	for {
		buf := make([]byte, 10240)
		n, err := lconn.Read(buf)
		if err != nil {
			break
		}

		mesg := changeHost(string(buf[:n]), dsthost)

		//mesgslice:=strings.Fields(mesg)
		//fmt.Println(mesgslice[1])
		//fmt.Println(mesgslice[3])
		//fmt.Println(strings.Fields(mesg))

		// print request
		fmt.Printf(mesg)
		//fmt.Println(mesg)
		rconn.Write([]byte(mesg))    // send mesg to dest
	}
	lconn.Close()                  //close the conn from source
}

   //do not deal with response
func ReadResponse(lconn *net.TCPConn, rconn *net.TCPConn) {
	for {
		buf := make([]byte, 10240)
		n, err := lconn.Read(buf)
		if err != nil {
			break
		}

		//fmt.Println(string(buf[:n]))
		rconn.Write(buf[:n])     //send mesg to source
	}
	lconn.Close()
}

func changeHost(request string, newhost string) string {
	reg := regexp.MustCompile(`Host[^\r\n]+`)
	return reg.ReplaceAllString(request, newhost)
}




