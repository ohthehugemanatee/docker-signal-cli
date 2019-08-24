package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
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
	Args := os.Args[1:]
	if len(Args) < 2 {
		log.Fatal("you must provide two arguments: your phone number (already registered on this device with signal-cli, and the target Group ID")
	}
	CollectMessages(Args[0], Args[1], os.Stdout)
}

// CollectMessages from Signal-cli.
func CollectMessages(myPhone string, targetGroupID string, writer io.Writer) {
	cmd := exec.Command("signal-cli", "-u", myPhone, "receive", "-t", "-1", "--json")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Fatal error")
		log.Fatal(err)
	}
	FilterMessages(stdout, targetGroupID, writer)

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(writer, "Fatal error")
		log.Fatal(err)
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
