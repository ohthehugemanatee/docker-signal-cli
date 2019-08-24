package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/smtp"
	"os"
	"os/exec"
)

const mailSubject string = "Automatic image submission from Signal"
const mailBody string = "Dear Nixplay. Add this, please."

var NixplayEmail string
var MailServer string
var MailUser string
var MailPass string
var MailFrom string
var MailPort string = "587"

// Attachment from signal-cli
type Attachment struct {
	ContentType string
	ID          int
	Size        int
}

// Message from signal-cli
type Message struct {
	Envelope struct {
		Source       string
		SourceDevice int
		Relay        string
		Timestamp    int
		IsReceipt    bool
		DataMessage  struct {
			Timestamp        int
			Message          string
			ExpiresInSeconds int
			Attachments      []Attachment
			GroupInfo        struct {
				GroupID string
				Members string
				Name    string
				Type    string
			}
		}
		syncMessage string
		callMessage string
	}
}

func main() {
	phonePtr := flag.String("p", "", "the recipient account's phone number")
	groupPtr := flag.String("g", "", "The Signal Group ID to monitor")
	mailPtr := flag.String("e", "", "The destination Nixplay email")
	mailUserPtr := flag.String("e", "", "The SMTP user")
	mailPassPtr := flag.String("e", "", "The SMTP password")
	MPtr := flag.String("e", "", "The SMTP Server")
	mailFromPtr := flag.String("e", "", "The SMTP from address")

	flag.Usage = func() {
		fmt.Printf("Syntax:\n\tsignal-nixplay-bridge [flags]\nwhere flags are:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	MyPhone := *phonePtr
	TargetGroupID := *groupPtr
	NixplayEmail = *mailPtr
	MailServer = *MPtr
	MailUser = *mailUserPtr
	MailPass = *mailPassPtr
	MailFrom = *mailFromPtr

	if MyPhone == "" || TargetGroupID == "" || NixplayEmail == "" {
		flag.Usage()
		return
	}
	fmt.Printf("Monitoring...")
	CollectMessages(MyPhone, TargetGroupID, os.Stdout)
}

// CollectMessages from Signal-cli.
func CollectMessages(myPhone string, targetGroupID string, writer io.Writer) {
	cmd := exec.Command("signal-cli", "-u", myPhone, "receive", "-t", "-1", "--json")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	FilterMessages(stdout, targetGroupID, writer)

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

// FilterMessages from stdIn.
func FilterMessages(stdout io.Reader, targetGroupID string, writer io.Writer) {
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	go func() {
		for scanner.Scan() {
			text := scanner.Bytes()
			if text != nil {
				var message Message
				json.Unmarshal(text, &message)
				attachments := message.Envelope.DataMessage.Attachments
				if message.Envelope.DataMessage.GroupInfo.GroupID == targetGroupID && attachments != nil {
					for i := 0; i < len(attachments); i++ {
						if attachments[i].ContentType == "image/jpeg" &&
							attachments[i].Size > 512000 {
							fmt.Fprintf(writer, "Send attachment id %v", attachments[i].ID)
						}
					}
				}
			}
		}
	}()

}

func SendMail(file string) {
	body := "To: " + NixplayEmail + "\r\nSubject: " +
		mailSubject + "\r\n\r\n" + mailBody
	auth := smtp.PlainAuth("", MailUser, MailPass, MailServer)
	err := smtp.SendMail(MailServer+":"+MailPort, auth, MailFrom,
		[]string{NixplayEmail}, []byte(body))
	if err != nil {
		log.Fatal(err)
	}
}
