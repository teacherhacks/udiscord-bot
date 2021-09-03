package main

import (
)

func main() {
    BotRun()
    <-make(chan struct{})
    BotStop()
}
