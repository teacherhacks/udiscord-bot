package main

import (
    "log"
    "fmt"
    "github.com/bwmarrin/discordgo"
)

/* main handler for slash commands */
func SlashCommandHandler(s *discordgo.Session, m *discordgo.InteractionCreate) {
    if handler, ok := SlashCommandHandlers[m.ApplicationCommandData().Name]; ok {
        handler(s, m)
    }
}

func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

    if m.Author.ID == BotID { return }

    if m.Content == "ping" {
        _, _ = s.ChannelMessageSend(m.ChannelID, "pong")
    }
}

func JoinHandler(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
    guildID := m.Member.GuildID

    /* create dm with user and ask them register */
    userID := m.Member.User.ID

    userChannel, err := s.UserChannelCreate(userID)
    if err != nil { log.Printf("Cannot open private channel with user %v", userID) }

    /* send user a message to tell them to verify */
    guildInfo, err := s.Guild(guildID)
    if err != nil { log.Printf("Error finding guild information") }

    s.ChannelMessageSend(userChannel.ID, fmt.Sprintf("Hello! Welcome to %v", guildInfo.Name))

}
