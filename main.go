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
	"strconv"

	gomail "gopkg.in/gomail.v2"
)

// Constants for mails.
const (
	MailSubject string = "Automatic image submission from Signal, pl: My Playlist"
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
	MailPort      int
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
	cmd, stdout := StartSignal(MyPhone, TargetGroupID, writer)
	filenameChannel := make(chan string)
	defer close(filenameChannel)
	go processFile(filenameChannel)
	FilterMessages(stdout, TargetGroupID, filenameChannel, writer)
	err := cmd.Wait()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func parseFlags() {
	phonePtr := flag.String("u", "", "the recipient account's phone number")
	groupPtr := flag.String("g", "", "The Signal Group ID to monitor")
	mailPtr := flag.String("e", "", "The destination Nixplay email")
	mailUserPtr := flag.String("user", "", "The SMTP user")
	mailPassPtr := flag.String("pass", "", "The SMTP password")
	MPtr := flag.String("s", "", "The SMTP Server")
	mailFromPtr := flag.String("f", "", "The SMTP from address")
	mailPortPtr := flag.Int("p", 587, "The SMTP server port")

	flag.Usage = func() {
		fmt.Printf("Syntax:\n\tsignal-nixplay-bridge [flags]\nwhere flags are:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	MyPhone = *phonePtr
	TargetGroupID = *groupPtr
	NixplayEmail = *mailPtr
	MailServer = *MPtr
	MailUser = *mailUserPtr
	MailPass = *mailPassPtr
	MailFrom = *mailFromPtr
	MailPort = *mailPortPtr
	// Explicitly default MailPort because run.sh will pass an empty string.
	if MailPort == 0 {
		MailPort = 587
	}

	if MyPhone == "" || TargetGroupID == "" || NixplayEmail == "" || MailServer == "" || MailUser == "" || MailPass == "" || MailFrom == "" {
		flag.Usage()
		os.Exit(1)
	}
}

// StartSignal starts the Signal-cli.
func StartSignal(myPhone string, targetGroupID string, writer io.Writer) (*exec.Cmd, io.ReadCloser) {
	cmd := exec.Command("signal-cli", "-u", myPhone, "receive", "-t", "-1", "--json")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	return cmd, stdout
}

// FilterMessages from stdIn.
func FilterMessages(stdout io.Reader, targetGroupID string, filenameChannel chan string, writer io.Writer) {
	go func() {
		scanner := bufio.NewScanner(stdout)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			fmt.Fprintf(writer, "Starting a new iteration.")
			text := scanner.Bytes()
			if text != nil {
				var message Message
				json.Unmarshal(text, &message)
				attachments := message.Envelope.DataMessage.Attachments
				if message.Envelope.DataMessage.GroupInfo.GroupID == targetGroupID && attachments != nil {
					for i := 0; i < len(attachments); i++ {
						if attachments[i].ContentType == "image/jpeg" &&
							attachments[i].Size > 512000 {
							filenameChannel <- strconv.Itoa(attachments[i].ID)
							fmt.Fprintf(writer, "Send attachment id %v", attachments[i].ID)
						}
					}
				}
			}
		}
	}()
}

func processFile(filenameChannel chan string) {
	for fileName := range filenameChannel {
		originalFilePath := "/root/.local/share/signal-cli/attachments/" + fileName
		tmpFileName := "/tmp/" + fileName + ".jpg"
		copyFile(originalFilePath, tmpFileName)
		sendMail(tmpFileName)
		err := os.Remove(tmpFileName)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func sendMail(fileName string) {
	m := gomail.NewMessage()
	m.SetHeader("From", MailFrom)
	m.SetHeader("To", NixplayEmail)
	m.SetHeader("Subject", MailSubject)
	m.SetBody("text/html", MailBody)
	m.Attach(fileName)

	d := gomail.NewDialer(MailServer, MailPort, MailUser, MailPass)
	if err := d.DialAndSend(m); err != nil {
		log.Println(err)
	}
}

func copyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}
