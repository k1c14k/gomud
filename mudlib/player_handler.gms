package main

func HandleLine(line string) {
    if line == "north" {
        room.TryMove("north")
    } else {
        if line == "south" {
            room.TryMove("south")
        } else {
            player.Send("I don't understand that command.")
        }
    }
    room.SendDescription()
}
