package main

import (
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

    /* create channels for student */
    studentChannelPermissions := []*discordgo.PermissionOverwrite{
        {
            ID: guildID,
            Type: discordgo.PermissionOverwriteTypeRole,
            Allow: 0,
            Deny: 1024, // disable @everyone from viewing channel
        },
        {
            ID: m.User.ID,
            Type: discordgo.PermissionOverwriteTypeMember,
            Allow: 1024, // allow viewing
            Deny: 0,
        },
    }

    studentCategory, _ := s.GuildChannelCreate(guildID, fmt.Sprintf("%s's channels"), discordgo.ChannelTypeGuildCategory);
    s.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
        Name: "questions",
        Type: discordgo.ChannelTypeGuildText,
        ParentID: studentCategory.ID,
        Topic: "Use this channel to ask instructors any private questions",
        PermissionOverwrites: studentChannelPermissions,
    })

}
