package main

func SendDescription() {
    player.Send("You are in a room. There is a door to the north.")
}

func TryMove(direction string) {
    if direction == "north" {
        player.Send("You move north.")
        player.MoveTo("locations/room_b")
    } else {
        player.Send("You can't go that way.")
    }
}