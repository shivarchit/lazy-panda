<h1 align="center" id="title">Lazy Panda</h1>

<p align="center"><img src="https://socialify.git.ci/shivarchit/lazy-panda/image?description=1&amp;descriptionEditable=Lazy%20Panda&amp;font=Raleway&amp;forks=1&amp;issues=1&amp;language=1&amp;name=1&amp;owner=1&amp;pattern=Charlie%20Brown&amp;pulls=1&amp;stargazers=1&amp;theme=Dark" alt="project-image"></p>

<h2>🍰 Contribution Guidelines:</h2>

Open for contributions and suggestions

## Prerequisites

#### For MacOS:

Xcode Command Line Tools (And Privacy setting: [#277](https://github.com/go-vgo/robotgo/issues/277))

```
xcode-select --install
```

#### For Windows:

You can use mingw GCC directly from TDM Link
[TDM](https://jmeubank.github.io/tdm-gcc/img/dragon_logo1.gif)

OR use choco
```bash
  choco install mingw -y  
```

#### For everything else:

```
GCC

X11 with the XTest extension (the Xtst library)

"Clipboard": xsel xclip


"Bitmap": libpng (Just used by the "bitmap".)

"Event-Gohook": xcb, xkb, libxkbcommon (Just used by the "hook".)

```

##### Ubuntu:

```yml
# gcc
sudo apt install gcc libc6-dev

# x11
sudo apt install libx11-dev xorg-dev libxtst-dev

# Clipboard
sudo apt install xsel xclip

#
# Bitmap
sudo apt install libpng++-dev

# GoHook
sudo apt install xcb libxcb-xkb-dev x11-xkb-utils libx11-xcb-dev libxkbcommon-x11-dev libxkbcommon-dev

```

##### Fedora:

```yml
# x11
sudo dnf install libXtst-devel

# Clipboard
sudo dnf install xsel xclip

#
# Bitmap
sudo dnf install libpng-devel

# Hook
sudo dnf install libxkbcommon-devel libxkbcommon-x11-devel xorg-x11-xkb-utils-devel

```

    

##  Environment Variables

To run this project, you will need to add the following environment variables to your config.json file

`Port:` Default is 3010

`IP Address: ` Default is localhost

`Sys Tray Icon Path: ` Default is "panda.co" present in local server directory

`JWT Secret` 


  
<h2>💻 Built with</h2>

Technologies used in the project:

*   Kotlin
*   Golang

<h2>🛡️ License:</h2>

This project is licensed under the MIT
