package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"

	gomail "gopkg.in/gomail.v2"
)

// Constants for mails.
const (
	MailSubject string = "Automatic image submission from Signal"
	MailBody    string = "Dear Nixplay. Add this, please."
)

// Variables from CLI arguments
var (
	MyPhone       string
	TargetGroupID string
	NixplayEmail  string
	MailServer    string
	MailUser      string
	MailPass      string
	MailFrom      string
	MailPort      int = 587
)

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
	parseFlags()
	fmt.Printf("Monitoring...")
	writer := os.Stdout
	stdout := CollectMessages(MyPhone, TargetGroupID, writer)
	FilterMessages(stdout, TargetGroupID, writer)

}

func parseFlags() {
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
}

// CollectMessages from Signal-cli.
func CollectMessages(myPhone string, targetGroupID string, writer io.Writer) io.ReadCloser {
	cmd := exec.Command("signal-cli", "-u", myPhone, "receive", "-t", "-1", "--json")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return stdout
}

// FilterMessages from stdIn.
func FilterMessages(stdout io.Reader, targetGroupID string, writer io.Writer) {
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	wg := new(sync.WaitGroup)
	for scanner.Scan() {
		wg.Add(1)
		text := scanner.Bytes()
		go func(t []byte) {
			defer wg.Done()
			if t != nil {
				var message Message
				json.Unmarshal(t, &message)
				attachments := message.Envelope.DataMessage.Attachments
				if message.Envelope.DataMessage.GroupInfo.GroupID == targetGroupID && attachments != nil {
					for i := 0; i < len(attachments); i++ {
						if attachments[i].ContentType == "image/jpeg" &&
							attachments[i].Size > 512000 {
							// SendMail(strconv.Itoa(attachments[i].ID))
							fmt.Fprintf(writer, "Send attachment id %v", attachments[i].ID)
						}
					}
				}
			}
		}(text)
	}
	wg.Wait()

}

// SendMail using the CLI-defined params.
func SendMail(fileName string) {
	m := gomail.NewMessage()
	m.SetHeader("From", MailFrom)
	m.SetHeader("To", NixplayEmail)
	m.SetHeader("Subject", MailSubject)
	m.SetBody("text/html", MailBody)
	m.Attach("/root/.local/share/signal-cli/attachments/" + fileName)

	d := gomail.NewDialer(MailServer, MailPort, MailUser, MailPass)
	if err := d.DialAndSend(m); err != nil {
		log.Fatal(err)
	}
}
