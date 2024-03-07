#!/bin/bash

freeze artichoke.hs --execute "eza -l" --output eza.svg
freeze artichoke.hs --language haskell --output language.svg
freeze artichoke.hs --theme dracula --output theme.svg
freeze artichoke.hs --border.radius 8 --output corner-radius.svg
freeze artichoke.hs --border.radius 8 --window --output window.svg
freeze artichoke.hs --border.radius 8 --window --border.width 1 --output border.svg
freeze artichoke.hs --border.radius 8 --window --border.width 1 --padding 30,50,30,30 --output padding.svg
freeze artichoke.hs --border.radius 8 --window --border.width 1 --padding 30,50,30,30 --margin 50,60,100,60 --output margin.svg
freeze artichoke.hs -c full --output shadow.svg

