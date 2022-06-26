# Copyright (c) 2022 alwin
# 
# This software is released under the MIT License.
# https://opensource.org/licenses/MIT
build:
	go build && chmod a+x github-cleaner

dev:
	go run main.go