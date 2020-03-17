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
	message := ""
	if err != nil {
	  output.VerbosePrint(fmt.Sprintf("[-] %s", err))
	} else {
	   message = handshake(conn, profile)
	   output.VerbosePrint(fmt.Sprintf("[+] TCP established for %s", profile["paw"]))
	   listen(conn, profile, server)
	}
	return message
 }

func listen(conn net.Conn, profile map[string]interface{}, server string) {
    scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
        message := scanner.Text()
		bites, status := commands.RunCommand(strings.TrimSpace(message), server, profile)
		pwd, _ := os.Getwd()
		response := make(map[string]interface{})
		response["response"] = string(bites)
		response["status"] = status
		response["pwd"] = pwd
		jdata, _ := json.Marshal(response)
		conn.Write(jdata)
    }
}

func handshake(conn net.Conn, profile map[string]interface{}) string {
	//write the profile
	jdata, _ := json.Marshal(profile)
    conn.Write(jdata)
	conn.Write([]byte("\n"))

	//read back the paw and contact
    data := make([]byte, 512)
    n, _ := conn.Read(data)
    //extract tokens
    s := strings.Split(string(data[:n]),"\t")
    contact := strings.TrimSpace(s[0])
    paw := strings.TrimSpace(s[1])
    conn.Write([]byte("\n"))
	profile["paw"] = paw
	return contact
}
