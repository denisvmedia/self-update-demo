Executable File Update Demo Written in Go
=========================================

This application demonstrates how to update a running program (Windows only).

## Prerequisites

- Windows 10.
- Go 1.16 or newer.
- TCP ports 8080 and 8081 free and available.

## How to use

1. run `build.bat`
1. cd to `playground` folder
1. run `app.exe` (accept firewall rules if asked)
1. navigate to [http://localhost:8080](http://localhost:8080)
1. make sure the version you see is 1.0.0
1. check the update (click [Check for new version](http://localhost:8080/check))
1. you should see version 1.1.0 available
1. click [Upgrade](http://localhost:8080/install)
1. you should notice application restart and a temporary page saying that it will reload in 5 seconds
1. after reload you should see version 1.1.0 used
1. check updates (no more updates should be available now)

