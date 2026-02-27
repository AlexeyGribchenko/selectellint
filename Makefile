PLUGIN_NAME := selectellint
PLUGIN_PATH := /usr/local/lib/$(PLUGIN_NAME).so

.PHONY: all build install install-local clean

all: build

check-deps:
	@command -v go >/dev/null 2>&1 || { echo "Ошибка: go не установлен"; exit 1; }
	@command -v gcc >/dev/null 2>&1 || { echo "Ошибка: gcc не установлен"; exit 1; }

build: check-deps
	go mod download
	CGO_ENABLED=1 go build -buildmode=plugin -o $(PLUGIN_NAME).so ./plugin
	@echo "Плагин собран: $(PLUGIN_NAME).so"

install: build
	sudo cp $(PLUGIN_NAME).so $(PLUGIN_PATH)
	sudo chmod 644 $(PLUGIN_PATH)
	@echo "Плагин установлен в $(PLUGIN_PATH)"
	@echo "В .golangci.yml укажите: path: $(PLUGIN_PATH)"

install-local: build
	mkdir -p ./bin
	cp $(PLUGIN_NAME).so ./bin/
	@echo "Плагин установлен в ./bin/$(PLUGIN_NAME).so"
	@echo "В .golangci.yml укажите: path: ./bin/$(PLUGIN_NAME).so"

clean:
	rm -f $(PLUGIN_NAME).so