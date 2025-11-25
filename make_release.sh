#!/bin/bash

mkdir release
~/go/bin/gogio -target windows .
mv main.exe release/colan-game-library.exe
rm *.syso

~/go/bin/gogio -target macos .
mv main_amd64.app colan-game-library.app
zip -r mac-amd64-colan-game-library.zip colan-game-library.app
mv mac-amd64-colan-game-library.zip release/
rm -r colan-game-library.app

mv main_arm64.app colan-game-library.app
zip -r mac-arm64-colan-game-library.zip colan-game-library.app
mv mac-arm64-colan-game-library.zip release/
rm -r colan-game-library.app