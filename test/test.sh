#!/bin/sh

[ -d test ] && cd test
[ -f freeze-test ] && rm freeze-test
go build -o freeze-test ..
[ ! -f freeze-test ] && echo "Failed to build freeze-test" && exit 1
[ ! -d output ] && mkdir output/

if [ ! -d golden ]; then
  echo "No golden files, generating..."
  mkdir golden/

  for f in configurations/*.json; do
    filename=$(basename -- "$f")
    ./freeze-test --config $f --output golden/"${filename%.*}".svg $f
  done

  echo "Generated golden files, verify and commit."
  rm freeze-test
  exit 0
fi


FAILURES=0
for f in configurations/*.json; do
  filename=$(basename -- "$f")
  if [ ! -f golden/"${filename%.*}".svg ]; then
    echo "Generating golden file for $filename, verify and commit."
    ./freeze-test --config $f --output golden/"${filename%.*}".svg $f
    continue
  fi
  ./freeze-test --config $f --output output/"${filename%.*}".svg $f
  diff --color output/"${filename%.*}".svg golden/"${filename%.*}".svg
  [ $? -ne 0 ] && echo "=== Test failed for $filename ===" &&  FAILURES=$((FAILURES + 1))
done

rm freeze-test

[ $FAILURES -ne 0 ] && echo "$FAILURES tests failed" && exit 1
