//------------------------------------------------------//
//------------------------------------------------------//
//-------- GOLANG IRC BOT ------------------------------//
//-----------------------------------------------------//
//-----------------------------------------------------//

package main

import (
	"bufio" //Buffer to TCP stream to generate bytes
	"fmt"
	"github.com/msbranco/goconfig" //configuration parser
	"net"                          //TCP networking protocols in order to connect
	"os"
	"strings"
	"time" //To sleep

)

type Com struct {
	commands string
	output   string
}

var PING string = "PING :"
var inputServer string = ""
var pongdata string
var arraypong []string
var arrayIRCparse []string
var protection bool = false
var b byte = 2
var com []string
var res []string

func main() {

	//config parameters
	c, err := goconfig.ReadConfigFile("config.cfg")
	if err != nil {
		fmt.Println("Config not found!!!")
	}

	password, _ := c.GetString("parameters", "password")
	mail, _ := c.GetString("parameters", "mail")
	commands, _ := c.GetString("parameters", "speak")
	//map commands
	com = strings.Split(commands, "\n")
	var list_com string
	l := len(com)
	m := make(map[string]string)
	for i := 0; i < l; i++ {
		var res = strings.Split(com[i], "%")
		m[res[0]] = res[1]
		list_com += " " + res[0] //here I store the commands
	}

	SERVER, _ := c.GetString("parameters", "server")
	NICK, _ := c.GetString("parameters", "nick")
	CHANNEL, _ := c.GetString("parameters", "channel")
	USER, _ := c.GetString("parameters", "user")

	//welcome message
	welcome, _ := c.GetBool("parameters", "welcome")
	welcome_msg, _ := c.GetString("parameters", "message")
	fmt.Printf("Welcome Message: %s \n", welcome_msg)

	fmt.Printf("Starting IRC Bot...\n")
	fmt.Printf("Connecting to: %s \n", SERVER)
	fmt.Printf("Nickname is: %s \n", NICK)
	fmt.Printf("Joining Channel: %s \n", CHANNEL)
	//Begin connection to server

	irc, err := net.Dial("tcp", SERVER)

	//Initalize byte stream components
	reader := bufio.NewReader(irc)
	writer := bufio.NewWriter(irc)

	if err != nil {
		fmt.Printf("Connection failed: %s\n", err)
		os.Exit(1)
	} else {
		writer.WriteString(USER + "\r\n")
		writer.Flush()
		writer.WriteString("NICK " + NICK + "\r\n")
		writer.Flush()
		fmt.Printf("Got line: %s \n", inputServer)
		writer.Flush()
		for {
			inputServer, err = reader.ReadString('\n')
			if inputServer == "" {
				time.Sleep(50)
			} else {
				arrayIRCparse = strings.SplitAfterN(inputServer, ":", 3)
				if strings.HasPrefix(inputServer, "PING :") {
					fmt.Printf(inputServer)
					arraypong = strings.SplitAfterN(inputServer, ":", 2)
					fmt.Printf(arraypong[1])
					writer.WriteString("PONG " + arraypong[1] + "\r\n")
					writer.Flush()
					fmt.Printf("PONG :%s\n", arraypong[1])
					writer.WriteString("JOIN " + CHANNEL + "\r\n")
					writer.Flush()
				}

				if strings.Contains(inputServer, "JOIN :") && strings.Contains(inputServer, "!") && welcome {
					tmp := strings.Split(inputServer, "!")
					nickname := strings.Trim(tmp[0], ":")
					tmp2 := strings.Replace(welcome_msg, "$C", CHANNEL, -1)
					tmp2 = strings.Replace(tmp2, "$U", nickname, -1)
					if nickname != NICK {
						writer.WriteString("privmsg " + CHANNEL + " :" + tmp2 + "\n")
						writer.Flush()
					}
				}

				fmt.Printf("Got line: %s", inputServer)

				if len(arrayIRCparse) >= 3 {
					//Here it register or identify the bot
					if strings.Contains(arrayIRCparse[2], "You have not registered") {
						writer.WriteString("PRIVMSG NICKSERV : REGISTER " + password + " " + mail + "\n")
						writer.Flush()

					} else if strings.Contains(arrayIRCparse[2], "This nickname is registered and protected") {
						writer.WriteString("PRIVMSG NICKSERV : IDENTIFY " + password + "\n")
						writer.Flush()

					}
					//parse commands and output
					if strings.HasPrefix(arrayIRCparse[2], "!") {
						arrayIRCparse[2] = strings.Trim(arrayIRCparse[2], "\r\n")
						output := m[arrayIRCparse[2]]
						if output != "" {
							writer.WriteString("privmsg " + CHANNEL + " :" + output + "\n")
							writer.Flush()
						} else {
							writer.WriteString("privmsg " + CHANNEL + " :Available Commands:" + list_com + " \n")
							writer.Flush()
						}

					}

				}

				time.Sleep(50)
			}
		}
	}

}
