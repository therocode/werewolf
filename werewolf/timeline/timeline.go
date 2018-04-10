package timeline

import (
	"encoding/json"
	//"log"
	"reflect"
)

/* Event */
type Event struct {
	Name   string
	Before map[string]bool
	After  map[string]bool
}

/* Generator */
type Generator interface {
	generate() []Event
}

/* Roles */

// Game
type Game struct{}

func (instance Game) generate() []Event {
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

// Werewolf
type Werewolf struct{}

func (instance Werewolf) generate() []Event {
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

// Doctor
type Doctor struct{}

func (instance Doctor) generate() []Event {
	return []Event{
		Event{
			"doctor_heals",
			map[string]bool{"werewolves_kill": true},
			map[string]bool{"day_starts": true},
		},
	}
}

// Seer
type Seer struct{}

func (instance Seer) generate() []Event {
	return []Event{
		Event{
			"seer_identifies",
			map[string]bool{"werewolves_kill": true},
			map[string]bool{"day_starts": true},
		},
	}
}

/* Timeline */

type Timeline struct {
	generators []Generator
}

func FilterEvents(events []Event, predicate func(Event) bool) ([]Event, []Event) {
	trueResult := []Event{}
	falseResult := []Event{}
	for _, event := range events {
		if predicate(event) {
			trueResult = append(trueResult, event)
		} else {
			falseResult = append(falseResult, event)
		}
	}
	return trueResult, falseResult
}

func ContainsEvent(event Event, events []Event) bool {
	for _, otherEvent := range events {
		if event.Name == otherEvent.Name {
			return true
		}
	}

	return false
}

func MapEventToString(events []Event, f func(Event) string) []string {
	result := []string{}
	for _, event := range events {
		result = append(result, f(event))
	}

	return result
}

func str(v interface{}) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}

func CopyStringToBoolMap(m map[string]bool) map[string]bool {
	result := make(map[string]bool)
	for k, v := range m {
		result[k] = v
	}
	return result
}

func (instance *Timeline) Generate() []string {

	// Get a list of all events
	events := []Event{}
	for _, generator := range instance.generators {
		events = append(events, generator.generate()...)
	}

	// Iterate over all events and take any events with no preconditions
	// and put them in initial result list
	var results []Event
	results, events = FilterEvents(events, func(event Event) bool {
		return len(event.Before) == 0 && len(event.After) == 0
	})

	for {
		resultsBefore := results

		// Iterate over remaining event list, moving any events over that now have
		// their preconditions satisfied to the result list
		for _, event := range events {
			remainingBefores := CopyStringToBoolMap(event.Before)
			remainingAfters := CopyStringToBoolMap(event.After)
			var markedIndex int

			// Iterate over result list and first check off all befores, then all afters
			for index, result := range results {
				if len(remainingBefores) > 0 {
					if _, contains := remainingBefores[result.Name]; contains {
						delete(remainingBefores, result.Name)
					}

					// Mark index if all befores have been found
					if len(remainingBefores) == 0 {
						markedIndex = index + 1
					}
				}

				// Continue iterating over result list until all after conditions have been met
				if len(remainingBefores) == 0 && len(remainingAfters) > 0 {
					if _, contains := remainingAfters[result.Name]; contains {
						delete(remainingAfters, result.Name)
					}
				}
			}

			// If there are any remaining befores or afters, this event's conditions are not satisfied
			// and we move on to the next one. Otherwise we insert the event in the result list at the
			// specified index
			if len(remainingBefores) == 0 && len(remainingAfters) == 0 {
				results = append(results[:markedIndex], append([]Event{event}, results[markedIndex:]...)...)
			}
		}

		// Remove any events that were added to the result list from the event list
		_, events = FilterEvents(events, func(event Event) bool { return ContainsEvent(event, results) })

		// Keep doing this until result list has an unchanged length between two
		// iterations
		if reflect.DeepEqual(results, resultsBefore) {
			return MapEventToString(results, func(event Event) string { return event.Name })
		}
	}
}
