# go-log

> A log module base of log and file-rotatelogs

- `v0.0.1`->`v0.1.2`

  > remove needless path by config ThisFile

- `v0.2.0+`

  > remove needless path by go build flags:
  >
  > **UNIX:** 
  >
  > ```shell
  > go build -gcflags="all=-trimpath=${PWD}" -asmflags="all=-trimpath=${PWD}" \
  >   -o main main.go
  > ```
  >
  > **WIN:**
  >
  > ```shell
  > go build -gcflags="all=-trimpath=%cd%" -asmflags="all=-trimpath=%cd%" ^
  >   -o main.exe main.go
  > ```
  >
  > 

