package main

func HandleLine(line string) {
    var north_direction string
    north_direction = "north"
    south_direction := "south"
    if line == north_direction {
        room.TryMove("north")
    } else {
        if line == south_direction {
            room.TryMove("south")
        } else {
            player.Send("I don't understand that command.")
        }
    }
    room.SendDescription()
}
