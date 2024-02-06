#!/bin/sh

cd test

if [ -f freeze-test ]; then
  rm freeze-test
fi

go build -o freeze-test ..

if [ ! -f freeze-test ]; then
  echo "Failed to build freeze-test"
  exit 1
fi

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

if [ ! -d output ]; then
  mkdir output/
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
  if [ $? -ne 0 ]; then
    echo "=== Test failed for $filename ==="
    FAILURES=$((FAILURES + 1))
  fi
done

if [ $FAILURES -ne 0 ]; then
  echo "$FAILURES tests failed"
  rm freeze-test
  exit 1
fi

rm freeze-test
