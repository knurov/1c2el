#!/bin/bash
/usr/lib/postgresql/12/bin/pg_dump  -h vm-001 -U electrolab -d electrolab --schema-only \
-t coil \
-t transformer \
-t  serial_number \
-f schema.sql

/usr/lib/postgresql/12/bin/pg_dump  -h vm-001 -U electrolab -d electrolab --data-only \
-t coil \
-t transformer \
-t  serial_number \
-f data.sql

