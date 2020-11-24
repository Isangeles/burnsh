## Introduction
Burn Shell is command line interface for [Flame engine](https://github.com/isangeles/flame).

CLI uses [Burn](https://github.com/Isangeles/burn) to handle user input and communicate with engine.

All commands must be prefixed with '$' character.
## Build & Run
Get sources from git:
```
$ go get -u github.com/isangeles/burnsh
```
Build shell:
```
$ go build github.com/isangeles/burnsh
```
Now, specify the ID of a valid Flame module in the configuration file:

Create file `.burnsh` in the Burn Shell executable directory(or run Burn Shell to create it
automatically) and add the following line:
```
module:[module ID]
```
Burn Shell will search the default modules directory(`data/modules`) for a module with the specified ID.

Flame modules are available for download [here](http://flame.isangeles.pl/mods).

Run shell:
```
$ ./burnsh
```
## Module directory
All UI-related files must be stored in the `data/modules/[module name]/burnsh` directory.

Translations for the UI needs to be stored in the `burnsh/lang` sub-directory of the module directory.

You can find default translations for the GUI in the `res/lang` directory of this repository.

For example check [Arena](https://github.com/Isangeles/arena) module.
## Multiplayer
It's possible to join an online game hosted on the [Fire](https://github.com/isangeles/fire) server.

To connect to the remote server set the `fire` value in `.burnsh` config file to `true` and specify server host and port in `server` config value.

After that Burn Shell will try to establish a connection with the game server on startup.

If the connection was successful you can use the `login` command to log in to the server.
## Commands
To run Burn or Burn Shell command use '$' character as prefix.
Without prefix, command will be treated as text and printed to out or sent to active player
chat channel if game was started.
### Burn Shell build-in commands:

Create module:
```
  $newmod
```
Description: Starts new module creation dialog. New module will be created in 'data/modules' directory. New module contains one chapter and start area.

Create new character:
```
  $newchar
```
Description: starts new character creation dialog.

Start new game:
```
  $newgame
```
Description: starts new game dialog.

Load game:
```
  $loadgame
```
Description: starts load game dialog.

Import exported characters:
```
  $importchars
```
Description: imports all characters from XML files in
data/modules/[module]/characters directory.

Login:
```
$login
```
Description: starts dialog for the authorization with remote game server.

Set target:
```
  $target
```
Description: searches current area for nearby targets to set for active PC.

Target information:
```
  $tarinfo
```
Description: prints informations about active PC target.

Loot target:
```
  $loot
```
Description: transfers all items from current dead target to active PC.

Talk with with target:
```
  $talk
```
Description: starts dialog with current PC target.

Show quests journal:
```
  $quests
```
Description: shows active PC quests.

Use character skill:
```
  $useskill
```
Description: starts dialog to use one of active PC skills.

Crafting dialog:
```
  $crafting
```
Description: starts items crafting dialog.

Trade with target:
```
  $trade
```
Description: starts trade dialog with current PC target.

Train with target:
```
  $train
```
Description: starts training dialog with current PC target.

Exit program:
```
  $close
```
Description: terminates program.
## Scripts
Burn Shell supports [Ash](https://github.com/Isangeles/burn/tree/master/ash) scripting language.

To run Ash script use '%' prefix, scripts are executed from 'data/scripts' directory.
Use '&' suffix to run script in background.
## Contributing
You are welcome to contribute to project development.

If you looking for things to do, then check TODO file or contact maintainer(dev@isangeles.pl).

When you find something to do, create new branch for your feature.
After you finish, open pull request to merge your changes with master branch.
## Documentation
Source code documentation can be easily browsed with `go doc` command.

Documentation for config files in form of Troff pages is available under `doc` directory.

You can easily view documentation pages with `man` command.

For example to display documentation page for guiset command:
```
$ man doc/file/.burnsh
```
Note that documentation is still incomplete.
## Contact
* Isangeles <<dev@isangeles.pl>>
## License
Copyright 2018-2020 Dariusz Sikora <<dev@isangeles.pl>>

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston,
MA 02110-1301, USA.
