# gorc
gorc is a WIP modern terminal IRC client written in golang.

## Building
In order to build `gorc`, make sure you have the [Go toolchain](https://go.dev/) installed.

#### Steps
1. clone the gorc repo
```
git clone https://git.illusionman1212.tech/illusion/gorc.git
```
2. clone the bubbles repo
```
git clone https://github.com/charmbracelet/bubbles.git
```
3. apply this [patch](https://github.com/charmbracelet/bubbles/pull/87) to your local `bubbles`.
4. tell golang to replace the online `bubbles` dependency with the local version you just cloned by modifying **line 5** in `go.mod` and replacing the path after `=>` with the path of your local `bubbles` repo.
5. run `go build` in the root directory of the project.
This will download the dependencies and create an executable named "gorc".

### Note: the patch is a temporary measure until it's merged to master in `bubbles`
#### Note Note: prebuilt binaries will be provided once an initial release is done.

## Screenshots
TODO

## Keybindings
- Login Screen Bindings
	- `Tab` -> Move input focus down.
	- `Shift+Tab` -> Move input focus up.
	- `Space` -> Toggle TLS checkbox.
	- `Enter` -> 
		--- Move input focus down 
		--- Toggle TLS checkbox.
		--- Confirm connect button.
- Main Screen Bindings
	- All Panes:
		- `Tab` -> Move between panes forwards.
		- `Shift+Tab` -> Move between panes backwards.
		- `Left Arrow` -> Move to next channel.
		- `Right Arrow` -> Move to previous channel.
	- Main Pane:
		- `J,K` -> Scrolls viewport up, down one line.
		- `Up Arrow, Down Arrow` -> Scrolls viewport up, down one line.
		- `D,U` -> Scrolls viewport up, down half a page.
		- `F,B` -> Scrolls viewport up, down a full page.
		- `G` -> Scrolls viewport to the top.
		- `Shift+G` -> Scrolls viewport to the bottom.
	- Side Pane:
		- Same bindings as the Main Pane.

## Resources Used
### [IRCDocs](https://modern.ircdocs.horse/about.html), [IRCDocs Github](https://github.com/ircdocs/modern-irc)
### [RFC1459](https://datatracker.ietf.org/doc/html/rfc1459)
### [RFC2812](https://datatracker.ietf.org/doc/html/rfc2812)
### [IRCv3](https://ircv3.net/)

## License
gorc is licensed under the GPL-3 license.

