package main

import (
    "fmt"
    "strconv"
    "regexp"
    "time"
    "errors"
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
    {
        Name: "new-assignment",
        Type: discordgo.ChatApplicationCommand,
        Description: "Creates a new assignment",
        Options: []*discordgo.ApplicationCommandOption{
            {
                Type: discordgo.ApplicationCommandOptionString,
                Name: "assignment-name",
                Description: "name of the new assignment",
                Required: true,
            },
            {
                Type: discordgo.ApplicationCommandOptionString,
                Name: "due-date",
                Description: "date the assignment is due in YYYY-MM-DD HH:MM format",
                Required: true,
            },
        },
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
        if !adminCommandPrivledge(s, i) { return }

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

        interactionSuccess("sucessfully initialized server", s, i)
    },
    "purge": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
        if !adminCommandPrivledge(s, i) { return }

        guildID := i.Interaction.GuildID

        /* delete all channels */
        channels, _ := s.GuildChannels(guildID);
        for _, c := range channels {
            s.ChannelDelete(c.ID)
        }

        s.GuildChannelCreate(guildID, "general", discordgo.ChannelTypeGuildText)

        interactionSuccess("purged server", s, i)
    },
    "new-assignment": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
        /* do some validation on args */
        parseDateStringCommandArg(i.ApplicationCommandData().Options[1].StringValue())

        /* insert assignment into database */

    },
}

/* helpers */

/* restricts command to only be used by admins */
func adminCommandPrivledge(s *discordgo.Session, i *discordgo.InteractionCreate) bool {
    if (i.Interaction.Member.Permissions >> 3) & 0x1 == 0x1 { return true }
    interactionError("You do not have permission to execute this command.", s, i)
    return false
}

/* checks if string is in YYYY-MM-DD HH:MM format */
/* converts datestring option to date */
func parseDateStringCommandArg(dateString string) (*time.Time, error) {
    /* issues with the regex - it's crude since its only for user friendliness */
    /* doesnt check leap years */
    /* doesnt validate number of days in each month */
    /* doesnt check for zero month and dates ie 2021-00-00 */

    r := regexp.MustCompile(`(\d{4})-(0\d|1[012])-([012]\d|3[01]) ([01]\d|2[0123]):([012345]\d)`)
    matches := r.FindStringSubmatch(dateString)
    if matches == nil { return nil, errors.New("invalid datestring") }

    /* make this less horrendous */
    year, _ := strconv.Atoi(matches[1])
    month, _ := strconv.Atoi(matches[2])
    day, _ := strconv.Atoi(matches[3])
    hour, _ := strconv.Atoi(matches[4])
    minute, _ := strconv.Atoi(matches[5])
    t := time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC)
    return &t, nil
}

/* interaction responses */
func interactionSuccess(successMessage string, s *discordgo.Session, i *discordgo.InteractionCreate) {
    s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseChannelMessageWithSource,
        Data: &discordgo.InteractionResponseData{
            Content: successMessage,
        },
    })
}

func interactionError(errorMessage string, s *discordgo.Session, i *discordgo.InteractionCreate) {
    s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseChannelMessageWithSource,
        Data: &discordgo.InteractionResponseData{
            Content: fmt.Sprintf("**[ERROR]** %s", errorMessage),
        },
    })
}

