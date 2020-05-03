### Prerequisite
It also requires Nanomsg to communicate with Gipan. Please refer to https://github.com/nanomsg/nanomsg for installation.

### Run
```
$ go run main.go show \
    --resolution 640x480 \
    --fps 30 \
    --framesrc ipc://../gipan/imageframes.ipc  \
    --soundsrc ipc://../gipan/soundframes.ipc  \
    --keyevtqueue ipc://../gipan/keys.ipc
```
or just `make run`
