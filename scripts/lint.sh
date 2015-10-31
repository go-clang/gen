 #!/bin/bash

echo "gofmt:"
export GOFMT=$(gofmt -l . 2>&1 | grep --invert-match -P "(_gen.go|testdata/)")
echo "$GOFMT"
$(exit $(echo -n "$GOFMT" | wc -l))

echo "goerrcheck:"
export GOERRCHECK=$(errcheck ./... 2>&1 | grep --invert-match -P "(_gen.go|testdata/)")
echo "$GOERRCHECK"
$(exit $(echo -n "$GOERRCHECK" | wc -l))

echo "govet:"
export GOVET=$(go tool vet -all=true -v=true . 2>&1 | grep --invert-match -P "(_gen.go|testdata/|Checking file|\%p of wrong type|can't check non-constant format)")
echo "$GOVET"
$(exit $(echo -n "$GOVET" | wc -l))

echo "golint:"
export GOLINT=$(golint ./... 2>&1 | grep --invert-match -P "(_gen.go|testdata/|_string.go:)")
echo "$GOLINT"
$(exit $(echo -n "$GOLINT" | wc -l))
