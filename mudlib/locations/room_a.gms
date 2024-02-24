package main

func SendDescription() {
    player.Send("You are in a room. There is a door to the north.")
}

func TryMove(direction string) {
    some_var := 3
    int_var := 2
    if direction == "north" {
        player.Send("You move north.")
        player.Send(some_var/int_var)
        player.MoveTo("locations/room_b")
    } else {
        player.Send(some_var*int_var)
        player.Send("You can't go that way.")
    }
}