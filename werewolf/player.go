package werewolf

func (instance *Werewolf) playerJoin(nick string) {
	instance.participants[nick] = Player{}
}

func (instance *Werewolf) getPlayer(nick string) *Player {
	if val, exists := instance.participants[nick]; exists {
		return &val
	}
	return nil
}
