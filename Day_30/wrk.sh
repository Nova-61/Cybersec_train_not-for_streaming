#!/bin/bash

# Пример запуска wrk для нагрузки на /users
wrk -t12 -c400 -d30s http://localhost:8080/users
