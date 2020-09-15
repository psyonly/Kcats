package tools

import (
	"strings"
)

// Message SendMessage
type Message struct {
	Sender   string
	Receiver string
	Text     string
}

// ErrMessage is a templete
var ErrMessage = Message{
	Sender:   "Server",
	Receiver: "Local",
	Text:     "E 你输入了不存在的ID，请检查是否存在拼写错误或者其他问题。",
}

// NewMSG is to get a new one
func NewMSG(sr, rc string, text string) (msg *Message) {
	msg = &Message{
		Sender:   sr,
		Receiver: rc,
		Text:     text,
	}
	return
}

// ParseMSG is a func which could make a []byte to MSG
// it will return a msg
func ParseMSG(raw []byte) (msg *Message) {
	msg = &Message{}
	str := string(raw)
	tem := strings.Split(str, "\n")
	SR := strings.Split(tem[0], " ")
	msg.Sender = SR[0]
	msg.Receiver = SR[1]
	msg.Text = tem[1]
	return
}

// Segment is to make a msg to segment
func (msg *Message) Segment() []byte {
	str := msg.Sender + " " + msg.Receiver + "\n" + msg.Text + "\r"
	return []byte(str)
}

// DeCodeMSG a msg string to readable Text
func DeCodeMSG(msg []byte) (ans string) {
	tem := string(msg)
	all := strings.Split(tem, "\n")
	sr := strings.Split(all[0], " ")
	who, text := sr[0], all[1]
	return "[" + who + "]:" + text
}
