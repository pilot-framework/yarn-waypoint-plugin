# Waypoint Plugin Template

This is a Build plugin for the Waypoint ecosystem. It's intended to be used along with Yarn to build static files for frontend frameworks.

## Install steps

1. Clone this repository

```shell
git clone https://github.com/pilot-framework/yarn-waypoint-plugin.git
```

2. (Optional) The compiled binaries should already be up to date, but if you want to ensure the binaries match the source you can run the Makefile to compile the plugin. The `Makefile` will build the plugin for all architectures.

```shell
cd yarn-waypoint-plugin

make
```

You should see the following output:

```shell
Build Protos
protoc -I . --go_out=plugins=grpc:. --go_opt=paths=source_relative ./builder/output.proto
protoc -I . --go_out=plugins=grpc:. --go_opt=paths=source_relative ./registry/output.proto
protoc -I . --go_out=plugins=grpc:. --go_opt=paths=source_relative ./platform/output.proto
protoc -I . --go_out=plugins=grpc:. --go_opt=paths=source_relative ./release/output.proto

Compile Plugin
# Clear the output
rm -rf ./bin
GOOS=linux GOARCH=amd64 go build -o ./bin/linux_amd64/waypoint-plugin-mytest ./main.go 
GOOS=darwin GOARCH=amd64 go build -o ./bin/darwin_amd64/waypoint-plugin-mytest ./main.go 
GOOS=windows GOARCH=amd64 go build -o ./bin/windows_amd64/waypoint-plugin-mytest.exe ./main.go 
GOOS=windows GOARCH=386 go build -o ./bin/windows_386/waypoint-plugin-mytest.exe ./main.go 
```

3. For the plugin to be able to work with your `waypoint.hcl` file, you must run the following command to copy the compiled binary into your Waypoint configuration directory.

```shell
user@machine:~/yarn-waypoint-plugin $ make install
```

## Building with Yarn

This plugin ultimately acts as an alias to the `yarn build` command, and currently requires Yarn to be installed on the machine where the build is occurring.

To utilize the plugin, you should specify the execution directory (i.e. the top-level directory for your front-end application), and the output directory (i.e. the build directory where your static files end up).

An example build stanza would look like the following:

```
# with a tree of
.
├── README.md
├── client
│   ├── build # where yarn outputs build files
│   ├── ...

build {
   use "yarn" {
      directory = "client"
      output = "client/build"
   }
}
```