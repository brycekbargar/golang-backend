// Package articledomain represents the Model for Articles. This includes Authorship and commentary.
// There is some amount of reliance on the User domain, notable the favoriting system.
package articledomain

import "time"

var periodicNow = time.Now().UTC()
var fearfulNow = time.Now().UTC()
var complexNow = time.Now().UTC()

// Fixture is a slice of valid article objects (mostly for test data purposes).
var Fixture = []*Article{
	{
		"periodic-slug",
		"periodic title",
		"periodic description",
		"periodic body",
		[]string{
			"shared tag 1",
			"periodic tag 1",
			"periodic tag 2",
			"shared tag 2",
		},
		periodicNow,
		periodicNow,
		"periodic author",
		[]*Comment{
			{
				1,
				"periodic brawny comment body",
				time.Now(),
				"periodic brawny comment author",
			},
		},
	},
	{
		"complex-slug",
		"complex title",
		"complex description",
		"complex body",
		make([]string, 0),
		complexNow,
		complexNow,
		"complex author",
		[]*Comment{
			{
				1,
				"complex hysterical comment body",
				time.Now(),
				"complex hysterical comment author",
			},
			{
				7,
				"complex blue comment body",
				time.Now(),
				"complex blue comment author",
			},
			{
				4,
				"complex unbiased comment body",
				time.Now(),
				"complex unbiased comment author",
			},
		},
	},
	{
		"fearful-slug",
		"fearful title",
		"fearful description",
		"fearful body",
		[]string{
			"shared tag 1",
			"shared tag 2",
		},
		fearfulNow,
		fearfulNow,
		"fearful author",
		make([]*Comment, 0),
	},
}
