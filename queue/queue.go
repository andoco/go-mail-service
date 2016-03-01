package queue

import (
	"log"
	"time"

	"bitbucket.org/andoco/gomailservice/delivery"
)

var queueChannel chan *delivery.MailMessage
var done chan bool
var enqueuer MailEnqueuer
var dequeuer MailDequeuer
var listeners []chan *delivery.MailMessage

func Enqueue(msg *delivery.MailMessage) {
	queueChannel <- msg
}

func Listen() chan *delivery.MailMessage {
	c := make(chan *delivery.MailMessage)
	listeners = append(listeners, c)

	return c
}

func Start() {
	queueChannel = make(chan *delivery.MailMessage)
	done = make(chan bool)
	enqueuer = newEnqueuer()
	dequeuer = newDequeuer()
	listeners = []chan *delivery.MailMessage{}
	go process(queueChannel, done, enqueuer)
	go processDequeue()
}

func Stop() {
	log.Print("Stopping queue")
	log.Print("Closing mail queue channel")
	close(queueChannel)
	<-done
	log.Print("Stop queue complete")
}

func process(c chan *delivery.MailMessage, done chan bool, enqueuer MailEnqueuer) {
	for msg := range c {
		enqueuer.Enqueue(msg)
	}

	log.Print("Finished processing mail queue channel")
	done <- true
}

func processDequeue() {
	log.Print("Processing queue")
	for {
		msg := dequeuer.Dequeue()

		if msg != nil {
			log.Printf("Processing message %s", msg.Id)

			for _, c := range listeners {
				c <- msg
			}
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func newEnqueuer() MailEnqueuer {
	return FileMailEnqueuer{}
}

func newDequeuer() MailDequeuer {
	return FileMailDequeuer{}
}

type MailEnqueuer interface {
	Enqueue(msg *delivery.MailMessage)
}

type MailDequeuer interface {
	Dequeue() *delivery.MailMessage
}
