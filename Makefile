APP_NAME := invi
OUTPUT_DIR := build
# VERSION := $(shell git describe --tags --always)
VERSION := $(shell date +%Y.%m.%d.%H%M)-$(shell git describe --tags --always)

PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 windows/amd64 android/arm64

LDFLAGS := "-s -w -X 'main.version=$(VERSION)'"

.PHONY: all clean

build: $(PLATFORMS)
	@echo "Build completed successfully!"

$(PLATFORMS):
	@mkdir -p $(OUTPUT_DIR)
	$(eval GOOS := $(word 1, $(subst /, ,$@)))
	$(eval GOARCH := $(word 2, $(subst /, ,$@)))
	$(eval OUTPUT_NAME := $(OUTPUT_DIR)/$(APP_NAME)-$(GOOS)-$(GOARCH)-$(VERSION)$(if $(filter windows, $(GOOS)),.exe,))

	@echo "Building for $(GOOS)/$(GOARCH)..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags=$(LDFLAGS) -o $(OUTPUT_NAME)

	@if [ $$? -ne 0 ]; then \
		echo "An error occurred while building for $(GOOS)/$(GOARCH)"; \
		exit 1; \
	fi

compress:
	@which upx > /dev/null 2>&1 || { echo "UPX not found. Please install it to use compression."; exit 1; }
	@echo "Compressing binaries with UPX..."
	@upx $(OUTPUT_DIR)/*
	# @upx --best --lzma $(OUTPUT_DIR)/*
	# @upx --ultra-brute -9 $(OUTPUT_DIR)/*

clean:
	@echo "Cleaning build directory..."
	rm -rf $(OUTPUT_DIR)
	@echo "Clean completed!"

archivate:
	@echo "Compressing binaries..."
	for file in $(OUTPUT_DIR)/*; do \
		if echo $$file | grep -q "windows" || echo $$file | grep -q "darwin"; then \
			zip -j $$file.zip $$file; \
		else \
			tar -czf $$file.tar.gz -C $(OUTPUT_DIR) $$(basename $$file); \
		fi; \
		rm $$file; \
	done

localdeploy: clean build
	@echo "Copying binaries to /usr/local/bin..."
	cp $(OUTPUT_DIR)/$(APP_NAME)-linux-amd64-* /usr/local/bin/$(APP_NAME)
