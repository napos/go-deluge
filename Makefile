# Run me to verify that all tests pass and all binaries are buildable before pushing!
# If you do not, then Travis will be sad.

BUILD_TYPE?=build

# Everything; this is the default behavior
all: format library

# go fmt ftw
format:
	go list ./... | grep -v vendor | xargs go fmt

library:
	go $(BUILD_TYPE)

# Building Examples
example: examples
examples:
	for example in $$(ls examples); do \
		go $(BUILD_TYPE) ./examples/$$example; \
	done

clean:
	rm $$(ls examples)

fixmes: fixme
fixme:
	@grep -rn FIXME * | grep -v vendor/ | grep -v README.md | grep --color FIXME || echo "No FIXMES!  YAY!"

.PHONY: examples clean
