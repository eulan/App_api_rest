#!/bin/bash

cockroach start-single-node --insecure --store=json-test --listen-addr=localhost:26257 --http-addr=localhost:8080 --background

cockroach sql --insecure --host=localhost:26257
