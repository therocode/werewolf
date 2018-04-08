package werewolf

import "fmt"

func (instance *Werewolf) playerJoin(nick string) {
	instance.participants[nick] = &Player{instance.irc, nick, &RoleVillager{}}
	instance.participants[nick].message("welcome. The main game channel is '%s'. Be sure to join", instance.mainChannel.name)
}

func (instance *Werewolf) getPlayer(nick string) *Player {
	if val, exists := instance.participants[nick]; exists {
		return val
	}
	return nil
}

func (player *Player) message(msg string, args ...interface{}) {
	player.irc.Privmsg(player.name, fmt.Sprintf(msg, args...))
}
