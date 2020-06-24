SILENT = @

TARGET = cmd/app.go

RUN_CMD = go run $(TARGET)

RUN_ARGS = --resolution 480x320 \
	--fps 24 \
	--framesrc ipc://../gipan/imageframes.ipc  \
	--soundsrc ipc://../gipan/soundframes.ipc  \
	--keyevtqueue ipc://../gipan/keys.ipc

run:
	$(RUN_CMD) $(RUN_ARGS)
