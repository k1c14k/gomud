package main

import (
    "player"
)

func HandleLine(line string) {
    player.Send("\r\n\"+
        "                        Hello!\r\n"+
        "\r\n" +
        "         Let's play!\r\n")
}