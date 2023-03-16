package utils

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/segmentio/kafka-go"
	gomail "gopkg.in/mail.v2"
	"log"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func Diff(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}

// ex: SendMail("127.0.0.1:25", (&mail.Address{"from name", "from@example.com"}).String(), "Email Subject", "message body", []string{(&mail.Address{"to name", "to@example.com"}).String()})
func SendMail(addr, from, subject, body string, to []string) error {
	r := strings.NewReplacer("\r\n", "", "\r", "", "\n", "", "%0a", "", "%0d", "")

	c, err := smtp.Dial(addr)

	if err != nil {
		fmt.Print(err)
		return err
	}
	defer c.Close()
	if err = c.Mail(r.Replace(from)); err != nil {
		fmt.Print(err)
		return err
	}
	for i := range to {
		to[i] = r.Replace(to[i])
		if err = c.Rcpt(to[i]); err != nil {
			fmt.Print(err)
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		fmt.Print(err)
		return err
	}

	msg := "To: " + strings.Join(to, ",") + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"Content-Transfer-Encoding: base64\r\n" +
		"\r\n" + base64.StdEncoding.EncodeToString([]byte(body))
	fmt.Println(msg)
	_, err = w.Write([]byte(msg))
	if err != nil {
		fmt.Print(err)
		return err
	}
	err = w.Close()
	if err != nil {
		fmt.Print(err)
		return err
	}
	return c.Quit()
}

func SendMailGmail(to string, subject string, body string) error {
	fmt.Printf("SendMailGmail to(%s)/subject(%s)/body(%s)\n", to, subject, body)
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", "chartdrug@gmail.com")

	// Set E-Mail receivers
	m.SetHeader("To", to)

	// Set E-Mail subject
	m.SetHeader("Subject", subject)

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/html", body)

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, "chartdrug@gmail.com", "ynoltzxpnezlyrto")

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		//panic(err)
		return err
	}
	return nil
}

func SendMailError(subject string, body string) {
	err := SendMailGmail("chartdrug@gmail.com", "Ошибка сервера: "+subject, body)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func SendMsgInjection(uids []string) {
	w := &kafka.Writer{
		Addr:         kafka.TCP("45.8.191.170:9092"),
		Topic:        "calc_injection",
		RequiredAcks: kafka.RequireAll,
	}

	var msgs []kafka.Message
	for _, uid := range uids {

		msgs = append(msgs, kafka.Message{
			Key: []byte(uid),
			// create an arbitrary message payload for the value
			Value: []byte(uid),
		})
	}

	// Passing a context can prevent the operation from blocking indefinitely.
	switch err := w.WriteMessages(context.Background(), msgs...).(type) {
	case nil:
	case kafka.WriteErrors:
		for i := range msgs {
			if err[i] != nil {
				// handle the error writing msgs[i]
				log.Printf("could not write message default" + err[i].Error())
			}
		}
	default:
		// handle other errors
		log.Printf("could not write message default" + err.Error())
	}
}

/*
func SendMsgInjection(uids []string) {

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "calc_injection",
	})
	var msgs []kafka.Message
	for _, uid := range uids {

		msgs = append(msgs, kafka.Message{
			Key: []byte(uid),
			// create an arbitrary message payload for the value
			Value: []byte(uid),
		})
	}
	err := w.WriteMessages(context.Background(), msgs...)
	if err != nil {
		log.Printf("could not write message " + err.Error())
	}
}

*/

/*
func SendMsgInjection(uids []string) {

	var (
		brokerList = kingpin.Flag("brokerList", "List of brokers to connect").Default("89.208.219.91:9092").Strings()
		topic      = kingpin.Flag("topic", "Topic name").Default("calc_injection").String()
		maxRetry   = kingpin.Flag("maxRetry", "Retry limit").Default("5").Int()
	)

	kingpin.Parse()
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = *maxRetry
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(*brokerList, config)
	if err != nil {
		SendMailError("sarama.NewSyncProducer", err.Error())
		log.Printf(err.Error())
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Printf(err.Error())
		}
	}()
	for _, uid := range uids {
		msg := &sarama.ProducerMessage{
			Topic: *topic,
			Value: sarama.StringEncoder(uid),
		}
		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			SendMailError("producer.SendMessage", err.Error())
			log.Printf(err.Error())
		}

		log.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", *topic, partition, offset)
	}


}
*/
