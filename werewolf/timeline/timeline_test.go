package timeline

import (
	"reflect"
	"testing"
)

// Test roles
type Game struct{}

func (Game) generate() []Event {
	return []Event{
		Event{
			"night_starts",
			map[string]bool{},
			map[string]bool{},
		},
		Event{
			"day_starts",
			map[string]bool{"night_starts": true},
			map[string]bool{},
		},
	}
}

type Werewolf struct{}

func (Werewolf) generate() []Event {
	return []Event{
		Event{
			"werewolves_see_each_other",
			map[string]bool{"night_starts": true},
			map[string]bool{"day_starts": true},
		},
		Event{
			"werewolves_kill",
			map[string]bool{"werewolves_see_each_other": true},
			map[string]bool{"day_starts": true},
		},
	}
}

type Doctor struct{}

func (Doctor) generate() []Event {
	return []Event{
		Event{
			"doctor_heals",
			map[string]bool{"werewolves_kill": true},
			map[string]bool{"day_starts": true},
		},
	}
}

type Seer struct{}

func (Seer) generate() []Event {
	return []Event{
		Event{
			"seer_identifies",
			map[string]bool{"werewolves_kill": true},
			map[string]bool{"day_starts": true},
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
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %v but got %v", expected, result)
	}
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
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %v but got %v", expected, result)
	}
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
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %v but got %v", expected, result)
	}
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
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %v but got %v", expected, result)
	}
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
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %v but got %v", expected, result)
	}
}
