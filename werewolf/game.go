package werewolf

func (instance *Werewolf) NewGame(owner string) {
	instance.state = GameStateInvite
	instance.owner = owner
	instance.playerJoin(owner)
}

func (instance *Werewolf) startGame() {
	instance.state = GameStateDay
	//instance.irc.Privmsgf(channel, "Unknown command '%s'", command)
}
