package main

import (
    "github.com/bwmarrin/discordgo"
)

/* main handler for slash commands */
func SlashCommandHandler(s * discordgo.Session, i *discordgo.InteractionCreate) {
    if handler, ok := SlashCommandHandlers[i.ApplicationCommandData().Name]; ok {
        handler(s, i)
    }
}

func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

    if m.Author.ID == BotID { return }

    if m.Content == "ping" {
        _, _ = s.ChannelMessageSend(m.ChannelID, "pong")
    }
}

