package main

import (
    "player"
)

func HandleLine(line string) {
    player.Send("\r\n"+
        "\t\tHello!\r\n"+
        "\r\n" +
        "\t\t\tLet's play!\r\n")
}