package main

import (
    "github.com/bwmarrin/discordgo"
)

var SlashCommands = []*discordgo.ApplicationCommand{
    {
        Name: "hello",
        Type: discordgo.ChatApplicationCommand,
        Description: "basic test",
    },
    {
        Name: "init",
        Type: discordgo.ChatApplicationCommand,
        Description: "Initializes an empty server",
    },
    {
        Name: "purge",
        Type: discordgo.ChatApplicationCommand,
        Description: "Nukes the server; used for development",
    },
}

var SlashCommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
    "hello": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
        s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
            Type: discordgo.InteractionResponseChannelMessageWithSource,
            Data: &discordgo.InteractionResponseData{
                Content: "yay slash command",
            },
        })
    },
    "init": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
        guildID := i.Interaction.GuildID

        /* create general channels */
        s.GuildChannelCreate(guildID, "announcements", discordgo.ChannelTypeGuildText)
        s.GuildChannelCreate(guildID, "rules", discordgo.ChannelTypeGuildText)

        /* create admin channels */
        adminCategory, _ := s.GuildChannelCreate(guildID, "admin", discordgo.ChannelTypeGuildCategory);
        s.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
            Name: "admin-commands",
            Type: discordgo.ChannelTypeGuildText,
            ParentID: adminCategory.ID,
        })

        s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
            Type: discordgo.InteractionResponseChannelMessageWithSource,
            Data: &discordgo.InteractionResponseData{
                Content: "sucessfully initialized server",
            },
        })
    },
    "purge": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

        /* delete all channels */

        s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
            Type: discordgo.InteractionResponseChannelMessageWithSource,
            Data: &discordgo.InteractionResponseData{
                Content: "purged server",
            },
        })
    },
}

