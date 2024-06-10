TARGETS = \
	linux/amd64 \
	linux/arm64 \
	windows/amd64 \
	windows/arm64 \
	darwin/amd64 \
	darwin/arm64

all: $(TARGETS)
	@echo "Build completed for all targets."

$(TARGETS):
	@echo "Building for $@..."
	@GOOS=$(word 1, $(subst /, ,$@)) GOARCH=$(word 2, $(subst /, ,$@)) go build -o ./hamon_$(word 1, $(subst /, ,$@))_$(word 2, $(subst /, ,$@)) main.go

clean:
	rm -rf ./dist

.PHONY: all clean $(TARGETS)
