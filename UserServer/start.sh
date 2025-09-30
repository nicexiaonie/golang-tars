#!/bin/bash
set -ex
make
./UserServer --config=config/config.conf
