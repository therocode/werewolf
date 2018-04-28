package timeline

import (
	"reflect"
	"testing"
)

// Test roles
type Game struct{}

func (Game) Generate() []Event {
	return []Event{
		Event{
			Name:   "night_starts",
			Before: map[string]bool{},
			After:  map[string]bool{},
		},
		Event{
			Name:   "day_starts",
			Before: map[string]bool{"night_starts": true},
			After:  map[string]bool{},
		},
	}
}

type Werewolf struct{}

func (Werewolf) Generate() []Event {
	return []Event{
		Event{
			Name:   "werewolves_see_each_other",
			Before: map[string]bool{"night_starts": true},
			After:  map[string]bool{},
		},
		Event{
			Name:   "werewolves_kill",
			Before: map[string]bool{"werewolves_see_each_other": true},
			After:  map[string]bool{"day_starts": true},
		},
	}
}

type Villager struct{}

func (Villager) Generate() []Event {
	return []Event{
		Event{
			Name:   "lynch",
			Before: map[string]bool{"day_starts": true},
			After:  map[string]bool{},
		},
	}
}

type Doctor struct{}

func (Doctor) Generate() []Event {
	return []Event{
		Event{
			Name:   "doctor_heals",
			Before: map[string]bool{"werewolves_kill": true},
			After:  map[string]bool{"day_starts": true},
		},
	}
}

type Seer struct{}

func (Seer) Generate() []Event {
	return []Event{
		Event{
			Name:   "seer_identifies",
			Before: map[string]bool{"werewolves_kill": true},
			After:  map[string]bool{"day_starts": true},
		},
	}
}

func TestEmptyTimeline(t *testing.T) {
	// GIVEN an empty set of generators
	generators := map[Generator]bool{}

	// WHEN generating the timeline
	result := Generate(generators)

	// THEN an empty timeline pops out
	if len(result) > 0 {
		t.Fatalf("Expected empty result but got %v", result)
	}
}

func TestTimelineWithGame(t *testing.T) {
	// GIVEN a single game generator
	generators := map[Generator]bool{
		Game{}: true,
	}

	// WHEN generating the timeline
	result := Generate(generators)

	// THEN a timeline with only basic game events pops out
	expected := []string{
		"night_starts",
		"day_starts",
	}
	verifyTimeline(t, result, expected)
}

func TestTimelineWithGameAndWerewolf(t *testing.T) {
	// GIVEN a timeline with a game generator and a werewolf generator
	generators := map[Generator]bool{
		Werewolf{}: true,
		Game{}:     true,
	}

	// WHEN generating the timeline
	result := Generate(generators)

	// THEN a timeline with game events and werewolf events pops out
	expected := []string{
		"night_starts",
		"werewolves_see_each_other",
		"werewolves_kill",
		"day_starts",
	}
	verifyTimeline(t, result, expected)
}

func TestTimelineWithGameAndWerewolfAndDoctor(t *testing.T) {
	// GIVEN
	generators := map[Generator]bool{
		Game{}:     true,
		Werewolf{}: true,
		Doctor{}:   true,
	}

	// WHEN
	result := Generate(generators)

	// THEN
	expected := []string{
		"night_starts",
		"werewolves_see_each_other",
		"werewolves_kill",
		"doctor_heals",
		"day_starts",
	}
	verifyTimeline(t, result, expected)
}

func TestTimelineWithGameAndWerewolfAndVillager(t *testing.T) {
	// GIVEN
	generators := map[Generator]bool{
		Game{}:     true,
		Werewolf{}: true,
		Villager{}: true,
	}

	// WHEN
	result := Generate(generators)

	// THEN
	expected := []string{
		"night_starts",
		"werewolves_see_each_other",
		"werewolves_kill",
		"day_starts",
		"lynch",
	}
	verifyTimeline(t, result, expected)
}

func TestTimelineWithGameAndWerewolfAndSeer(t *testing.T) {
	// GIVEN
	generators := map[Generator]bool{
		Game{}:     true,
		Werewolf{}: true,
		Seer{}:     true,
	}

	// WHEN
	result := Generate(generators)

	// THEN
	expected := []string{
		"night_starts",
		"werewolves_see_each_other",
		"werewolves_kill",
		"seer_identifies",
		"day_starts",
	}
	verifyTimeline(t, result, expected)
}

func TestTimelineWithGameAndWerewolfAndDoctorAndSeer(t *testing.T) {
	// GIVEN
	generators := map[Generator]bool{
		Game{}:     true,
		Werewolf{}: true,
		Seer{}:     true,
		Doctor{}:   true,
	}

	// WHEN
	result := Generate(generators)

	// THEN
	expected := []string{
		"night_starts",
		"werewolves_see_each_other",
		"werewolves_kill",
		"doctor_heals",
		"seer_identifies",
		"day_starts",
	}
	verifyTimeline(t, result, expected)
}

func verifyTimeline(t *testing.T, actual []Event, expected []string) {
	actualEventNames := mapEventToString(actual, func(event Event) string { return event.Name })
	if !reflect.DeepEqual(actualEventNames, expected) {
		t.Fatalf("Expected %v but got %v", expected, actualEventNames)
	}
}

func mapEventToString(events []Event, f func(Event) string) []string {
	result := []string{}
	for _, event := range events {
		result = append(result, f(event))
	}

	return result
}
