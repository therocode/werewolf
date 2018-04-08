package werewolf

type RoleVillager struct {
}

func (role *RoleVillager) name() string {
	return "villager"
}

func (role *RoleVillager) dayCyclePhase() DayCyclePhase {
	return nil
}
