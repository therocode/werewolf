package timeline

import (
	"reflect"
	"testing"
)

func TestEmptyTimeline(t *testing.T) {
	// GIVEN a timeline with no generators
	timeline := Timeline{}

	// WHEN generating the timeline
	result := timeline.Generate()

	// THEN an empty timeline pops out
	if len(result) > 0 {
		t.Fatalf("Expected empty result but got %v", result)
	}
}

func TestTimelineWithGame(t *testing.T) {
	// GIVEN a timeline with a game generator
	timeline := Timeline{[]Generator{Game{}}}

	// WHEN generating the timeline
	result := timeline.Generate()

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
	timeline := Timeline{[]Generator{Werewolf{}, Game{}}}

	// WHEN generating the timeline
	result := timeline.Generate()

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
	// GIVEN a timeline with a game generator and a werewolf generator
	timeline := Timeline{[]Generator{
		Werewolf{},
		Game{},
		Doctor{},
	}}

	// WHEN generating the timeline
	result := timeline.Generate()

	// THEN a timeline with game events and werewolf events pops out
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
	// GIVEN a timeline with a game generator and a werewolf generator
	timeline := Timeline{[]Generator{Game{}, Werewolf{}, Seer{}}}

	// WHEN generating the timeline
	result := timeline.Generate()

	// THEN a timeline with game events and werewolf events pops out
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
	// GIVEN a timeline with a game generator and a werewolf generator
	timeline := Timeline{[]Generator{
		Seer{},
		Doctor{},
		Werewolf{},
		Game{},
	}}

	// WHEN generating the timeline
	result := timeline.Generate()

	// THEN a timeline with game events and werewolf events pops out
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
