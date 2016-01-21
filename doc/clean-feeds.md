# Clean Feed representations

With increased usage of our Feed feature we observe unintended representations. Those include events which reference entities (users, objects, connections) which don't exist anymore or duplicated events for interactions like follow-unfollow in short succession.

Ideally we represent a clean list of entities from a multitude of sources (social graph, global, second degree social, potentially interest graph). The challenge presents itself in that we are operating on isolated buckets of documents without referential integrity. While this is a delibrate design decision we need to account for it.

A side-effect of this is that we potentially transport all entities referenced in feed entries for easy "joining" in the consuming client.

constraints:
* combine entities from multiple sources
* deduplicate entities
* filter out entities which reference deleted entities
* sort entities

entities:
* events
* posts

sources:
* global
* social graph
* second degree social (connection updates of my connections)

### API

No visible API changes are planned.

### Internals

In the following code examples we use events as the basis to illustrate how the interactions could work out.

Gather and filter:

``` go
type source func() ([]Event, error)

func gather(sources ...source) ([]Event, error) {
	events := []Event{}

	for _, source := range sources {
		es, err := s(events)
		if err != nil {
			return nil, err
		}

		events = append(events, es)
	}

	return events, nil
}

type fiter func(Event) bool

func sift(events []Event, fs ...filter) []Event {
	es := events[:0]

	for _, event := range events {
		keep := true

		for _, f := range filters {
			if f(event) {
				keep = false
				break
			}
		}

		if keep {
			es = append(es, event)
		}
	}

	return es
}

// ...

	es, err := gather(
		globalEvents,
		socialGraphEvents,
		secondDegreeEvents,
	)
	if err != nil {
		return err
	}

	// Fetch dependent information like objects, origins and targets.

	es = sift(
		es,
		objectExists(objects),
		originExists(origins),
		targetExists(targets),
	)

// ...
```
