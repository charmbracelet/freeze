#!/bin/bash

freeze artichoke.hs --radius 8 --output corner-radius.svg
freeze artichoke.hs --radius 8 --window --output window.svg
freeze artichoke.hs --radius 8 --window --border --output border.svg
freeze artichoke.hs --radius 8 --window --border --padding 30,50,30,30 --output padding.svg
freeze artichoke.hs --radius 8 --window --border --padding 30,50,30,30 --margin 20 --output margin.svg
freeze artichoke.hs --radius 8 --window --border --padding 20,40,20,20 --shadow --output shadow.svg

