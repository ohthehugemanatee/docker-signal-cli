package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

}

// CollectMessages from Signal-cli.
func CollectMessages(writer io.Writer) {
	myPhone := "+4915146621809"
	cmd := exec.Command("signal-cli", "-u", myPhone, "receive", "-t", "-1", "--json")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Fatal error")
		log.Fatal(err)
	}
	FilterMessages(stdout, writer)

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(writer, "Fatal error")
		log.Fatal(err)
	}
}

// FilterMessages from stdIn.
func FilterMessages(stdout io.Reader, writer io.Writer) {
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	go func() {
		for scanner.Scan() {
			text := scanner.Bytes()
			if text != nil {
				var message Message
				json.Unmarshal(text, &message)
				attachments := message.Envelope.DataMessage.Attachments
				// if the message meets criteria, write into a channel
				if (message.Envelope.DataMessage.GroupInfo.GroupID == "DsFSSsmOQH2yx6UTGlgj3A==") && (attachments != nil) {
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
