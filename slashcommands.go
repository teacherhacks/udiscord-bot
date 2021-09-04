package main

import (
    _ "fmt"
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

        /* create channels =-=-=-=-=-=-= */
        /* general channels */
        s.GuildChannelCreate(guildID, "announcements", discordgo.ChannelTypeGuildText)
        s.GuildChannelCreate(guildID, "rules", discordgo.ChannelTypeGuildText)

        /* admin channels */
        adminCategory, _ := s.GuildChannelCreate(guildID, "admin", discordgo.ChannelTypeGuildCategory);
        s.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
            Name: "admin-commands",
            Type: discordgo.ChannelTypeGuildText,
            ParentID: adminCategory.ID,
        })

        /* create roles =-=-=-=-=-=-=-= */
        role, _ := s.GuildRoleCreate(guildID);
        s.GuildRoleEdit(guildID, role.ID, "Instructor", 16718213, true, 8, true);

        role, _ = s.GuildRoleCreate(guildID);
        s.GuildRoleEdit(guildID, role.ID, "TA", 16760604, true, 8, true);

        role, _ = s.GuildRoleCreate(guildID);
        s.GuildRoleEdit(guildID, role.ID, "Student", 1889791, true, 174016814656, true);

        s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
            Type: discordgo.InteractionResponseChannelMessageWithSource,
            Data: &discordgo.InteractionResponseData{
                Content: "sucessfully initialized server",
            },
        })
    },
    "purge": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
        guildID := i.Interaction.GuildID

        /* delete all channels */
        channels, _ := s.GuildChannels(guildID);
        for _, c := range channels {
            s.ChannelDelete(c.ID)
        }

        s.GuildChannelCreate(guildID, "general", discordgo.ChannelTypeGuildText)

        s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
            Type: discordgo.InteractionResponseChannelMessageWithSource,
            Data: &discordgo.InteractionResponseData{
                Content: "purged server",
            },
        })
    },
}

