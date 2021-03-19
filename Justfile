@unit:
  go test ./...

@build:
  go build .

@run path dry:
  ./filededupe {{path}} {{dry}}

@br: build (run "-p /mnt/c/Users/gemini/Downloads/dedupe-test-files" "-d")