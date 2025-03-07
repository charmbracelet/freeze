#!/bin/bash
set -ex

# Parse input values
input="$(apify actor get-input)"
code="$(jq -r .code <<< "$input")"

args=()

language="$(jq -r .language <<< "$input")"
if [[ "$language" != "" ]]; then
    args+=(--language $language)
fi

line_numbers="$(jq -r .showLineNumbers <<< "$input")"
if [[ $line_numbers == 'true' ]]; then
    args+=(--show-line-numbers)
fi

window="$(jq -r .window <<< "$input")"
if [[ $window == 'true' ]]; then
    args+=(--window)
fi

# Generate screenshot and save it to output.png
echo "$code" |
    freeze -o output.png "${args[@]}"

# Upload output.png to key-value store
apify actor set-value output.png --contentType image/png < output.png

# Construct an output object and push it to the dataset (Actor results)
echo '{}' |
    jq ".image = \"$(apify actor get-public-url output.png)\"" | 
    jq ".language = \"$language\"" |
    apify actor push-data
