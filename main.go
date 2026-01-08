package main

import (
	"github.com/IBM/sarama"
	"io"
	"log"
	"net/http"
)

func main() {

	healthcheck := func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "OK\n")
	}
	http.HandleFunc("/health", healthcheck)

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	producer, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		log.Fatal("Failed to create producer", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	//topics := func(w http.ResponseWriter, r *http.Request) {
	//	topicName := strings.TrimPrefix(r.URL.String(), "/topics/")
	//	io.WriteString(w, topicName)
	//}
	//http.HandleFunc("/topics/", topics)

	topics := func(w http.ResponseWriter, r *http.Request) {
		topicName := r.PathValue("topic")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read body", http.StatusBadRequest)
			return
		}

		partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
			Topic: topicName,
			Value: sarama.ByteEncoder(body),
		})
		if err != nil {
			log.Printf("FAILED to send message: %s\n", err)
		} else {
			log.Printf("> message sent to partition %d at offset %d\n", partition, offset)
			w.WriteHeader(http.StatusOK)
		}
	}

	http.HandleFunc("/topics/{topic}", topics)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
