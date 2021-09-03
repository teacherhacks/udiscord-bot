package main

import (
    "log"
    "fmt"
    "github.com/bwmarrin/discordgo"
)

var BotID string
var DGSession *discordgo.Session

func BotRun() {
    DGSession, err := discordgo.New(BotToken)
    if err != nil { panic(err.Error()) }

    DGSession.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
    DGSession.AddHandler(MessageHandler);
    DGSession.AddHandler(SlashCommandHandler);

    /* grab the bot user's id */
    botUser, err := DGSession.User("@me")
    if err != nil { panic(err.Error()) }
    BotID = botUser.ID

    err = DGSession.Open()
    if err != nil { panic(err.Error()) }
    fmt.Println("Bot is now running...")

    /* initialize slash commands */
    for _, c := range SlashCommands {
        _, err := DGSession.ApplicationCommandCreate(DGSession.State.User.ID, "", c)
        if err != nil {
            log.Printf("Failed to create slash command: %v", c.Name, err)
        }
    }

}

func BotStop() {
    DGSession.Close()
}

