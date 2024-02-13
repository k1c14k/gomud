package main

import (
    "player"
)

func HandleLine(line string) {
    player.Send("\r\n"+
        "\t\tHello!\r\n"+
        "\r\n" +
        "\t\t\tLet's play!\r\n")
    player.Send(player.String())
    if line == "enter" {
        player.Send("Entered")
    }
    if line == "exit" {
        player.Send("Exited")
    } else {
        player.Send("Not exited")
    }
}