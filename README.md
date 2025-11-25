# My Game Library

![](media/example.png)

Hi, this is a tool to make it easy to download and play my games. I currently only support MacOS and Windows, but will, one day, support Linux. I recommend that you download from the releases if you are interesting.

## Building Locally

```
go run main.go
```

## Building a Release

```
go install gioui.org/cmd/gogio@latest
```

Then you can build:
GOOS=windows GOARCH=amd64 go build -o gamelibrary.exe .

```
mkdir release
~/go/bin/gogio -target windows .
mv main.exe release/win-gamelibrary.exe
rm *.syso

~/go/bin/gogio -target macos .
mv main_amd64.app	 release/mac-amd-gamelibrary.app
mv main_arm64.app release/mac-arm-gamelibrary.app
```
