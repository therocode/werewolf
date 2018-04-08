package werewolf

type DayCyclePhaseWerewolf struct {
}

func (phase *DayCyclePhaseWerewolf) start() {
}

func (phase *DayCyclePhaseWerewolf) end() {
}

func (phase *DayCyclePhaseWerewolf) handleCommand(player *Player, command string, args []string) {
}

type RoleWerewolf struct {
}

func (role *RoleWerewolf) name() string {
	return "werewolf"
}

func (role *RoleWerewolf) dayCyclePhase() DayCyclePhase {
	return &DayCyclePhaseWerewolf{}
}
