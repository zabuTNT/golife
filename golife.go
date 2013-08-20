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
	commands, _ := c.GetString("parameters", "commands")
	results, _ := c.GetString("parameters", "results")
	com = strings.Split(commands, ";")
	res = strings.Split(results, ";")

	SERVER, _ := c.GetString("parameters", "server")
	NICK, _ := c.GetString("parameters", "nick")
	CHANNEL, _ := c.GetString("parameters", "channel")
	USER, _ := c.GetString("parameters", "user")

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

				fmt.Printf("Got line: %s \n", inputServer)

				if len(arrayIRCparse) >= 3 {
					//fmt.Printf("\nDisplaying Array Parsing: \n" + arrayIRCparse[2] + "\n") //Display the first portion of the parsing
					if strings.Contains(arrayIRCparse[2], "You have not registered") {
						writer.WriteString("PRIVMSG NICKSERV : REGISTER " + password + " " + mail + "\n")
						writer.Flush()

					} else if strings.Contains(arrayIRCparse[2], "This nickname is registered and protected") {
						writer.WriteString("PRIVMSG NICKSERV : IDENTIFY " + password + "\n")
						writer.Flush()

					}
          //TODO: Do a switch for the commands
          
					if strings.HasPrefix(arrayIRCparse[2], "!") {
						if strings.Contains(arrayIRCparse[2], com[0]) {
							writer.WriteString("privmsg " + CHANNEL + " :" + res[0] + "\n")
							writer.Flush()
						} else if strings.Contains(arrayIRCparse[2], com[1]) {
							writer.WriteString("privmsg " + CHANNEL + " :" + res[1] + "\n")
							writer.Flush()
						} else {
							writer.WriteString("privmsg " + CHANNEL + " :Available Commands: " + strings.Trim(strings.Replace(commands, ";", " ", -1), " ") + "\n")
							writer.Flush()
						}

					}

				}

				time.Sleep(50)
			}
		}
	}

}
