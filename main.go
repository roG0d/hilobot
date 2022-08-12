package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
)

func main() {
	// Load environment variables
	err := godotenv.Load("environment.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/receive", slashCommandHandler)

	fmt.Println("[INFO] Server listening")
	http.ListenAndServe(":8080", nil)
}

func slashCommandHandler(w http.ResponseWriter, r *http.Request) {
	s, err := slack.SlashCommandParse(r)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch s.Command {
	case "/hilobot":

		//Get the link to the message
		params := &slack.Msg{Text: s.Text}
		link := params.Text

		//Creating a confirmation response
		response := fmt.Sprintf("Altoke, enviando vergüenza al mensaje: %v", link)
		w.Header().Set("Content-Type", "application/json")

		resp := make(map[string]string)
		resp["response_type"] = "in_channel"
		resp["text"] = response
		JsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatal("Error happened in JSON Marshal. Error %s", err)
		}
		w.Write(JsonResp)

		//Link string treatment
		var splits []string = strings.Split(link, "/")

		var channelIdDone string = splits[len(splits)-2]
		fmt.Printf("ChannelId: %s\n", channelIdDone)

		var threadTs1 string = splits[len(splits)-1]
		var threadTs2 string = strings.Replace(threadTs1, "p", "", -1)
		index := 10
		threadTsDone := threadTs2[:index] + "." + threadTs2[index:]
		fmt.Printf("ThreadTs: %s", threadTsDone)

		//Creating the message requested
		api := slack.New(os.Getenv("SLACK_VERIFICATION_TOKEN"))
		channelId := channelIdDone
		threadTs := threadTsDone
		text := "Maquin, contesta dentro del hilo anda ;)"

		api.PostMessage(channelId, slack.MsgOptionTS(threadTs), slack.MsgOptionText(text, false))

	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
