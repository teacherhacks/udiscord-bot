package main

import (
    "strings"
    "log"
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
        Name: "register",
        Type: discordgo.ChatApplicationCommand,
        Description: "Registers user with role into class.",
        Options: []*discordgo.ApplicationCommandOption{
            {
                Name: "student",
                Description: "Register into class as a student.",
                Type: discordgo.ApplicationCommandOptionSubCommand,
            },
            {
                Name: "instructor",
                Description: "Register into class as a instructor; requires a password",
                Type: discordgo.ApplicationCommandOptionSubCommand,
                Options: []*discordgo.ApplicationCommandOption{
                    {
                        Type: discordgo.ApplicationCommandOptionString,
                        Name: "password",
                        Description: "Password to register as insturctor",
                        Required: true,
                    },
                },
            },
            {
                Name: "ta",
                Description: "Register into class as a TA; requires a password",
                Type: discordgo.ApplicationCommandOptionSubCommand,
                Options: []*discordgo.ApplicationCommandOption{
                    {
                        Type: discordgo.ApplicationCommandOptionString,
                        Name: "password",
                        Description: "Password to register as TA",
                        Required: true,
                    },
                },
            },
        },
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
        if !adminPrivledgeCommand(s, i) { return }

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
        if !adminPrivledgeCommand(s, i) { return }

        guildID := i.Interaction.GuildID

        /* delete all channels */
        channels, _ := s.GuildChannels(guildID)
        for _, c := range channels { s.ChannelDelete(c.ID) }

        /* delete all roles */
        roles, _ := s.GuildRoles(guildID)
        for _, r := range roles { s.GuildRoleDelete(guildID, r.ID) }

        s.GuildChannelCreate(guildID, "general", discordgo.ChannelTypeGuildText)

        interactionSuccess("purged server", s, i)
    },
    "register": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
        /* dm only is not good idea, since we don't have guildid */
        // if !dmCommand(s, i) { return }

        /* parse subcommands */
        switch i.ApplicationCommandData().Options[0].Name {
            case "student":
                registerStudent(s, i)
            case "instructor":
                registerInstructor(s, i)
            case "ta":
                registerTA(s, i)
        }

    },
    // TODO switch over to validation library
    "new-assignment": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
        assignmentNameArg := i.ApplicationCommandData().Options[0].StringValue();
        dueDateArg := i.ApplicationCommandData().Options[1].StringValue();

        guildID := i.Interaction.GuildID

        /* do some validation on args */
        dueDate, err := parseDateStringCommandArg(dueDateArg)

        /* insert assignment into database */
        _, err = DBNewAssignment(guildID, assignmentNameArg, dueDate.Unix());
        if err != nil {
            interactionError("Failed to create new assignment", s, i)
            return
        }
        interactionSuccess(fmt.Sprintf("Successfully created assignment %s", assignmentNameArg), s, i)

    },
}

/* sub commands for register */
func registerStudent(s *discordgo.Session, i *discordgo.InteractionCreate) {
    guildID := i.Interaction.GuildID

    /* create channels for student */
    studentChannelPermissions := []*discordgo.PermissionOverwrite{
        {
            ID: guildID,
            Type: discordgo.PermissionOverwriteTypeRole,
            Allow: 0,
            Deny: 1024, // disable @everyone from viewing channel
        },
        {
            ID: i.Member.User.ID,
            Type: discordgo.PermissionOverwriteTypeMember,
            Allow: 1024, // allow viewing
            Deny: 0,
        },
    }

    studentCategory, err := s.GuildChannelCreate(guildID, fmt.Sprintf("%s's channels", i.Member.User.Username), discordgo.ChannelTypeGuildCategory);
    if err != nil { log.Println(err) }

    s.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
        Name: "questions",
        Type: discordgo.ChannelTypeGuildText,
        ParentID: studentCategory.ID,
        Topic: "Use this channel to ask instructors any private questions",
        PermissionOverwrites: studentChannelPermissions,
    })

    studentRole, _ := findGuildRole(s, i.Interaction.GuildID, "Student")
    s.GuildMemberRoleAdd(guildID, i.Member.User.ID, studentRole)

    interactionSuccess("Sucessfully registered as student!", s, i)

}

func registerInstructor(s *discordgo.Session, i *discordgo.InteractionCreate) {
    guildID := i.Interaction.GuildID

    instructorRole, _ := findGuildRole(s, guildID, "Instructor")
    s.GuildMemberRoleAdd(guildID, i.Member.User.ID, instructorRole)

    interactionSuccess("Sucessfully registered as instructor!", s, i)
}

func registerTA(s *discordgo.Session, i *discordgo.InteractionCreate) {
    guildID := i.Interaction.GuildID

    taRole, _ := findGuildRole(s, guildID, "TA")
    s.GuildMemberRoleAdd(guildID, i.Member.User.ID, taRole)

    interactionSuccess("Sucessfully registered as TA!", s, i)
}

/* helpers */

/* restricts command to only be used by admins */
func adminPrivledgeCommand(s *discordgo.Session, i *discordgo.InteractionCreate) bool {
    if (i.Interaction.Member.Permissions >> 3) & 0x1 == 0x1 { return true }
    interactionError("You do not have permission to execute this command.", s, i)
    return false
}

/* restricts command to be used only in dm */
func dmCommand(s *discordgo.Session, i *discordgo.InteractionCreate) bool {
    channelInfo, err := s.Channel(i.Interaction.ChannelID)
    if err != nil { log.Println("Error retrieving channel information") }

    if channelInfo.Type == discordgo.ChannelTypeDM { return true }
    interactionError("Please run this command in a DM.", s, i)
    return false
}

/* finds a role based on name */
func findGuildRole(s *discordgo.Session, guildID string, roleName string) (string, error) {

    guildRoles, _ := s.GuildRoles(guildID)
    for _, r := range guildRoles {
        if strings.Compare(r.Name, roleName) == 0 { return r.ID, nil }
    }
    return "", errors.New(fmt.Sprintf("Cannot find role of name %s", roleName))
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
            Content: fmt.Sprintf("%s", successMessage),
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

