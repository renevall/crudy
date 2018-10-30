# CRUDY

CRUDY is a program to quickly generate CRUD applications following a DDD pattern.

Instead of providing convenient objects and methods, CRUDY aims to generate all the code required and give you access to it so you can edit it as you need to.

# Installing

Using Crudy aims to be easy, just use go get to install the executable

    go get -u github.com/renevall/crudy

# Getting Started

In order to get your code up and running you can write:

    crudy init <project name>

This will create the project folder under the $GOPATH/src/\<project name>

For Example:

    crudy init github.com/renevall/culv

Your code will be generated under $GOPATH/src/github.com/renevall/culv

After your code was generated make sure to run

    go get ./...

In order to get all dependencies

# Code Structure

Init will create the following project structure:

```
  ▾ appName/
    ▾ model/
        env.go
        config.go
        resource.go
    ▾ router/
        router.go
        resource-handler.go
    ▾ store/
        resource-store.go
      main.go
      config.go
      db.go
```

# TODO

- Resource Generator
- Testing Code
- More Flags and Parameters for customization
- Compatibility with mysql and sqlserver
