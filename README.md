# turboin
Hotkey creator for X written in go
# Installation
First install go using your system's package manager, then clone the repository anywhere on your operating system

    git clone https://github.com/ViSaturn/turboin

Then **READ THE INSTALL.SH FILE THROUGHLY** before running it like so

    cd turboin
    chmod +x install.sh
    sudo ./install.sh
    
# Usage
To Run it:

    turboin run

And to Stop it:

    turboin stop

# Dependencies
jezek's fork of xgbutil

# Features
Turboin is still in development and does not have all of it's
planned features, but here are the current features,
you can create normal hotkeys and double click hotkeys where clicking a key
twice quickly will run a different command than just clicking it once

# Planned Features
- Mode bindings, similiar to i3's (almost complete)
- Allow turboin to be ran as a daemon
- Global keybindings that can be run regardless of mode
- One keybinding for multiple commands
- variables for your config (e.g. mod4: super)
- settings for your cocnfig (e.g. seperator: " + ")
- MacOS and Windows 7, 8 and 10 support
- Improved installation process

# Example Configuration
http://github.com/ViSaturn/turboin/blob/main/config.go
