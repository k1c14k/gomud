package main

func GetDescription() string {
    return "You are in a small room. There is a door to the south."
}

func TryMove(direction string) {
    if direction == "south" {
        player.Send("You move to the south.")
        player.Send(7%2)
        player.Send(7-5)
        player.MoveTo("locations/room_a")
    } else {
        player.Send("You can't go that way.")
    }
}