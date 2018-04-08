package werewolf

type Role interface {
	name() string
	dayCyclePhase() DayCyclePhase
}
