# The Knight's Tour Problem

Program implemented in go language to solve the knights tour problem.

## Pre-rquisites

* VSC Visual Studio Code
* Docker (If you want to use devContainer)

## devContainer

For using devConainer, push key-combo **command-shift-P** and choose re-open in container. 
There are two files into *.devContainer folder*:
* **devcontainer.json** (devContainer definition)
* **Dockerfile** (file to build docker image)

## Build

```
go mod tidy

go build
```

## Run

```
go run main.go
```

## Solution example

```
This is an example generated with this program.

-----------------------------------------
| 01 | 38 | 59 | 36 | 43 | 48 | 57 | 52 | 
-----------------------------------------
| 60 | 35 | 02 | 49 | 58 | 51 | 44 | 47 | 
-----------------------------------------
| 39 | 32 | 37 | 42 | 03 | 46 | 53 | 56 | 
-----------------------------------------
| 34 | 61 | 40 | 27 | 50 | 55 | 04 | 45 | 
-----------------------------------------
| 31 | 10 | 33 | 62 | 41 | 26 | 23 | 54 | 
-----------------------------------------
| 18 | 63 | 28 | 11 | 24 | 21 | 14 | 05 | 
-----------------------------------------
| 09 | 30 | 19 | 16 | 07 | 12 | 25 | 22 | 
-----------------------------------------
| 64 | 17 | 08 | 29 | 20 | 15 | 06 | 13 | 
-----------------------------------------

```
