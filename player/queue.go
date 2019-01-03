package player

import (
	"github.com/bwmarrin/discordgo"
)

var messageQueue MessageQueue

// MessageQueue keeps state of the track queue
type MessageQueue struct {
	messages []*discordgo.MessageCreate
}

// push adds a message to the end of the message queue
func (q MessageQueue) push(message *discordgo.MessageCreate) {
	q.messages = append(q.messages, message)
}

// shift removes the message queue's first message and returns it
func (q MessageQueue) shift() *discordgo.MessageCreate {
	first, rest := q.messages[0], q.messages[1:]
	q.messages = rest
	return first
}

// Queue track
func Queue(message *discordgo.MessageCreate) {
	if len(messageQueue.messages) > 0 {
		messageQueue.push(message)
		// TODO: fetch video metadata and add it to queue
		// add url too.
		// [] struct {
		// 	discordgo.VideoInfo
		// 	url string
		// }
		return
	}
}

// Shift queue. Removes first element from queue and returns it
func Shift() *discordgo.MessageCreate {
	message := messageQueue.shift()
	return message
}
