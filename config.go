package main

/*
d888888P                   dP                oo
   88                      88
   88    dP    dP 88d888b. 88d888b. .d8888b. dP 88d888b.
   88    88    88 88'  `88 88'  `88 88'  `88 88 88'  `88
   88    88.  .88 88       88.  .88 88.  .88 88 88    88
   dP    `88888P' dP       88Y8888' `88888P' dP dP    dP

          Welcome to Turboin's example config!
*/

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
			// "mod4-z": []interface{} {"mocp --volume +10", "mocp --volume -10"},

			// Special keys example
			// Hold mod4 and z before also holding mod4-x
			// "mod4-x": []interface{} {"rofi -show drun", "dmenu_run"},
			"{mod4:z}-mod4-x": []interface{} {"mocp --unpause", "mocp --pause"},

			// Switch your mode (does not work yet)
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
