APP       := Demo
TARGET    := UserServer
MFLAGS    :=
DFLAGS    :=
CONFIG    := config
STRIP_FLAG:= N
J2GO_FLAG:= 


GO      = CGO_ENABLED=0 GOOS=linux GOARCH=amd64 /Users/nieyuanpei/.asdf/installs/golang/1.25.0/go/bin/go

-include scripts/makefile.tars.gomod.mk


TARS_WEB_HOST   := http://192.168.8.245:9300
TARS_WEB_TOKEN  := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiJhZG1pbiIsImlhdCI6MTc3MDE5MjU0MCwiZXhwIjoxODA1NTMwMTQwfQ.rzoqqeRWY6gmEalrUuIff1wlBkuuF23R7I0c_rjE2rQ"
UPLOAD_USER   := $(shell whoami)
UPLOAD_OS     := linux

