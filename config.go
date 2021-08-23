package main

// Welcome to Turboin's example config!

func setupModes() map[string]interface{} {
	modes := map[string]interface{} {
		// The default mode is the most important mode,
		// without it, Turboin wont run
		"default": map[string]interface{} {
			// Normal Click Examples
			"mod4-d": "dmenu_run",
			"mod4-return": "termite",

			// Double Click Examples
			// First command is the command for one click
			// Second command is the command for the double click
			"mod4-x": []interface{} {"mocp --unpause", "mocp --pause"},
			"mod4-z": []interface{} {"mocp --volume +10", "mocp --volume -10"},

			// Switch your mode
			"mod4-s": "has mode mode_2",
		}, // Make sure you have a comma at the end of every mode

		"mode_2": map[string]interface{} {
			// Return to the previous mode
			"mod4-z": "has mode default",

			// Do the stuff
			"mod4-d": "doom run",
		}, // Make sure you have a comma at the end of every mode
	}

	return modes
}
