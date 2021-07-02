#!/usr/bin/env bash

# ---- Starting ISO Websim .. ----
# ---- Setting ENV variables -----

TLS_ENABLED=false
TLS_CERT_FILE=/path/to/isosim/certs/cert.pem
TLS_KEY_FILE=/path/to/isosim/certs/key.pem

# ---- Starting App -----

go run . -http-port 8080 --log-level TRACE \
         -specs-dir ../../test/testdata/specs -html-dir ../../web \
         -data-dir ../../test/testdata/appdata
