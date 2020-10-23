package rest

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/doublegrey/microlog/db"
	"github.com/doublegrey/microlog/utils"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	topic := strings.TrimSpace(r.Header.Get("Topic"))
	if len(strings.TrimSpace(topic)) < 1 {
		fmt.Fprintf(w, "'Topic' field does not exist or is empty")
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Failed to read request body: %v", err)
		log.Printf("Failed to read request body: %v\n", err)
		return
	}
	if len(topic) > 0 && utils.Config.Kafka.CustomTopics {
		topic = fmt.Sprintf("%s_%s", utils.Config.Kafka.TopicPrefix, topic)
	} else {
		topic = utils.Config.Kafka.Topic
	}
	err = db.Produce(b, topic)
	if err != nil {
		fmt.Fprintf(w, "Failed to produce message: %v", err)
		log.Printf("Failed to produce message: %v\n", err)
		return
	}
}

func ValidateToken(handler func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Token")
		if len(strings.TrimSpace(token)) < 1 {
			w.WriteHeader(401)
			fmt.Fprintf(w, "Access Token does not exists or is invalid")
		} else {
			handler(w, r)
		}
	})
}
