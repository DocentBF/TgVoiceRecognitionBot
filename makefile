VOSK_BINARIES=vosk-linux-x86_64-0.3.45
VOSK_MODEL=vosk-model-small-ru-0.22

prepare:
	if [ ! -d "binaries/$(VOSK_BINARIES)" ]; then \
		echo "No binaries/$(VOSK_BINARIES)" \
		wget https://github.com/alphacep/vosk-api/releases/download/v0.3.45/$(VOSK_BINARIES).zip; \
		unzip -d ./binaries $(VOSK_BINARIES).zip; \
	fi
	if [ ! -d "models/$(VOSK_MODEL)" ]; then \
		wget https://alphacephei.com/vosk/models/$(VOSK_MODEL).zip; \
		unzip -d ./models $(VOSK_MODEL).zip; \
	fi

run-linux: prepare
	VOSK_PATH=$$(pwd)/binaries/$(VOSK_BINARIES) \
	LD_LIBRARY_PATH=$$VOSK_PATH \
	CGO_CPPFLAGS="-I $$VOSK_PATH" \
	CGO_LDFLAGS="-L $$VOSK_PATH" \
	go run ./internal

build-linux: prepare
	VOSK_PATH=$$(pwd)/binaries/$(VOSK_BINARIES) \
	LD_LIBRARY_PATH=$$VOSK_PATH \
	CGO_CPPFLAGS="-I $$VOSK_PATH" \
	CGO_LDFLAGS="-L $$VOSK_PATH -Wl,-rpath,./binaries/$(VOSK_BINARIES)" \
	GOARCH=amd64 \
	GOOS=linux \
	go build -ldflags="-w -s" -o ./TgVoiceRecognitionBot ./internal