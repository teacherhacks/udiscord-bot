package main

import (
    "fmt"
    "github.com/bwmarrin/discordgo"
)

var botID string

func main() {
    dg, err := discordgo.New(BotToken)
    if err != nil { panic(err.Error()) }

    dg.Identify.Intents = discordgo.IntentsGuildMessages
    dg.AddHandler(messageHandler);

    /* grab the bot user's id */
    botUser, err := dg.User("@me")
    if err != nil { panic(err.Error()) }
    botID = botUser.ID

    err = dg.Open()
    if err != nil { panic(err.Error()) }

    fmt.Println("Bot is now running...")

    <-make(chan struct{})
    dg.Close()
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

    if m.Author.ID == botID { return }

    if m.Content == "ping" {
        _, _ = s.ChannelMessageSend(m.ChannelID, "pong")
    }
}

