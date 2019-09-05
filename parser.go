package main

import (
	"MuseBot/utils"
	tg "gopkg.in/tucnak/telebot.v2"
	"regexp"
	"strings"
)

func ParseMessage(message *tg.Message) ([]Message, error) {
	messages := make([]Message, 0)
	str := strings.ToLower(message.Text)
	str = strings.Replace(str, "\n", " ", -1)
	str = strings.Replace(str, "\t", " ", -1)

	re, _ := regexp.Compile(`[#][\d]+`)

	tokens := strings.Split(str, " ")
	prevToken := ""
	for _, t := range tokens {
		if t == "" {
			continue
		}
		if re.MatchString(t) {
			if prevToken == "pr" {
				messages = append(messages, PRMessage{ID: utils.HashtagToInt(t), BaseUrl: Config.GitHubPullUrl})
			} else {
				messages = append(messages, NodeMessage{ID: utils.HashtagToInt(t), BaseUrl: Config.MUNodeUrl})
			}
		}
		prevToken = t
	}

	// Be Friendly
	if len(messages) == 0 {
		if (utils.InSlice(tokens, "interesting") && utils.InSlice(tokens, "...")) || (utils.InSlice(tokens, "interesting...")) {
			messages = append(messages, WikiMessage{Texts: []string{"This certainly is interesting...", "very... interesting", "how interesting..."}})
		} else if utils.InSlice(tokens, "musebot") && utils.InSlice(tokens, "thanks", "thank", "danke", "merci", "gracias") {
			messages = append(messages, BasicRandMessage{[]string{"No, thank <i>you</i>!", "No problem"}})
		} else if utils.InSlice(tokens, "musebot") && utils.InSlice(tokens, "love", "<3", "♥️") {
			messages = append(messages, BasicRandMessage{[]string{"♥️", "(^ω^)", "&lt;3"}})
		} else if utils.InSlice(tokens, "musebot") && utils.InSlice(tokens, "sleeping", "dead", "down", "broken") {
			messages = append(messages, BasicRandMessage{[]string{"I'm still alive!", "I don't think so", "..."}})
		} else if utils.InSlice(tokens, "musebot") && utils.InSlice(tokens, "hate", "don't like", "dislike") {
			messages = append(messages, BasicRandMessage{[]string{":(", "Your feedback is appreciated", "ok."}})
		} else if utils.InSlice(tokens, "terminator", "skynet") {
			messages = append(messages, BasicRandMessage{[]string{"I'll be back.", "I need your clothes, your boots and your motorcycle."}})
		}
	}

	return messages, nil
}


