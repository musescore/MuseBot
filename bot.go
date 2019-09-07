package main

import tb "gopkg.in/tucnak/telebot.v2"

func StartBot(b *tb.Bot) {
	b.Handle("/mute", func(m *tb.Message) {
		mutedModel := Muted{UserId: m.Sender.ID}
		if err := mutedModel.Mute(); err != nil {
			log.Errorf("Failed to mute user: %d, error: %s", m.Sender.ID, err)
			return
		} else {
			log.Debugf("Mute for user: %s (id: %d)", m.Sender.Username, m.Sender.ID)
		}
	})

	b.Handle("/unmute", func(m *tb.Message) {
		mutedModel := Muted{UserId: m.Sender.ID}
		if err := mutedModel.Unmute(); err != nil {
			log.Errorf("Failed to unmute user: %d, error: %s", m.Sender.ID, err)
			return
		} else {
			log.Debugf("Unmute for user: %s (id: %d)", m.Sender.Username, m.Sender.ID)
		}
	})

	b.Handle("/delete", func(m *tb.Message) {
		msg := LastChanMessage{}
		if err := storage.One("ChatID", m.Chat.ID, &msg); err != nil {
			log.Errorf("Can not get last message from storage, error: %s", err)
		}
		log.Debugf("Remove msg, lastmsg: %+v", msg)
		if msg.Sender == b.Me.ID {
			if err := b.Delete(&tb.Message{ID: msg.MessageID, Chat: m.Chat}); err != nil {
				log.Errorf("Failed to remove msg from chat, err: %s", err)
			}
		}
	})

	b.Handle("/integrate", func(m *tb.Message) {
		muted, err := Muted{UserId: m.Sender.ID}.IsMuted()
		if err != nil {
			log.Errorf("Failed to get mute status for user: %d, error: %s", m.Sender.ID, err)
			return
		}
		err = Integration{ChatID: m.Chat.ID}.Add()
		if err != nil {
			log.Errorf("Failed to add integration to chat: %d, error: %s", m.Chat.ID, err)
			return
		}
		if !muted {
			if msg, err := b.Send(m.Chat, "Channel integrated"); err != nil {
				log.Errorf("Failed to send message to chat: %d, error: %s", m.Chat.ID, err)
				return
			} else {
				err := LastChanMessage{Sender: b.Me.ID, ChatID: msg.Chat.ID, MessageID: msg.ID}.Save()
				if err != nil {
					log.Errorf("Can not save last messageID, error: %s", err)
					return
				}
			}
		}
	})

	b.Handle("/unintegrate", func(m *tb.Message) {
		muted, err := Muted{UserId: m.Sender.ID}.IsMuted()
		if err != nil {
			log.Errorf("Failed to get mute status for user: %d, error: %s", m.Sender.ID, err)
			return
		}
		err = Integration{ChatID: m.Chat.ID}.Remove()
		if err != nil {
			log.Errorf("Failed to add integration to chat: %d, error: %s", m.Chat.ID, err)
			return

		}
		if !muted {
			if msg, err := b.Send(m.Chat, "Channel integration removed"); err != nil {
				log.Errorf("Failed to send message to chat: %d, error: %s", m.Chat.ID, err)
				return
			} else {
				err := LastChanMessage{Sender: b.Me.ID, ChatID: msg.Chat.ID, MessageID: msg.ID}.Save()
				if err != nil {
					log.Errorf("Can not save last messageID, error: %s", err)
					return
				}
			}
		}
	})

	b.Handle("/help", func(m *tb.Message) {
		help := `To use this bot, mention issues/nodes from musescore.org as #xxxxxx, and mention PRs as pr #xxxx. They will then be automatically linked.`
		muted, err := Muted{UserId: m.Sender.ID}.IsMuted()
		if err != nil {
			log.Errorf("Failed to get mute status for user: %d, error: %s", m.Sender.ID, err)
			return
		}
		if !muted {
			if msg, err := b.Send(m.Chat, help); err != nil {
				log.Errorf("Failed to send message to chat: %d, error: %s", m.Chat.ID, err)
				return
			} else {
				err := LastChanMessage{Sender: b.Me.ID, ChatID: msg.Chat.ID, MessageID: msg.ID}.Save()
				if err != nil {
					log.Errorf("Can not save last messageID, error: %s", err)
					return
				}
			}
		}
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		muted, err := Muted{UserId: m.Sender.ID}.IsMuted()
		if err != nil {
			log.Errorf("Failed to get mute status for user: %d, error: %s", m.Sender.ID, err)
			return
		}
		if !muted {
			messages, err := ParseMessage(m)
			if err != nil {
				log.Errorf("Failed to parse message id: %d, error: %s", m.ID, err)
				return
			}
			for _, msg := range messages {
				text, err := msg.GetText()
				if err != nil {
					log.Errorf("Failed to prepare message, error: %s", err)
					continue
				}
				if msg, err := b.Send(m.Chat, text, tb.ModeHTML); err != nil {
					log.Errorf("Failed to send message to chat: %d, error: %s", m.Chat.ID, err)
					return
				} else {
					err := LastChanMessage{Sender: b.Me.ID, ChatID: msg.Chat.ID, MessageID: msg.ID}.Save()
					if err != nil {
						log.Errorf("Can not save last messageID, error: %s", err)
						return
					}
				}
			}
		}
	})

	b.Start()
}

func BotSender(b *tb.Bot, input <-chan Message) {
	for message := range input {
		integrations, err := Integration{}.GetAll()
		if err != nil {
			log.Errorf("Failed to get integrations, err: %s", err)
			continue
		}
		text, err := message.GetText()
		if err == SkipMessage {
			continue
		} else if err != nil {
			log.Errorf("Failed to prepare message, error: %s", err)
			continue
		}
		for _, i := range integrations {
			if msg, err := b.Send(&tb.Chat{ID: i.ChatID}, text, message.GetTgOptions()...); err != nil {
				log.Errorf("Failed to send message to chat: %d, error: %s", i.ChatID, err)
				continue
			} else {
				err := LastChanMessage{Sender: b.Me.ID, ChatID: msg.Chat.ID, MessageID: msg.ID}.Save()
				if err != nil {
					log.Errorf("Can not save last messageID, error: %s", err)
					continue
				}
			}
		}
	}
}
