# Hypunch

## a completely open source and super light weight **UDP hole punching server** written in golang.

### How to use

First you should be familiar with go and have golang installed on your computer - [Go official website](https://go.dev/learn/)

To get started with Hypunch, you first need to clone this repo or simply copy and paste the source code from hypunch.go into your project.

`git clone https://github.com/Hyper-Reflector-Team/hypunch.git`

Likely you already have an external server to run this code on for users to connnect, but if you don't you can use a VPS service like a [Digital Ocean Droplet](https://www.digitalocean.com/products/droplets)

### Install external dependencies

We currently use google's unique id generator for generating a "match" id for the users connecting \n.
`go get github.com/google/uuid`

Alternatively you are free to modify the server to not use UUID and remain dependency free.

### Run the server

You can either use
`go run hypunch.go`

Or compile the binary and run it.
`go build hypunch.go`
Then
`./hypunch`

### External Clients example

You can check the external client for an example of how to connect and send keep-alive messages on the client side.
The example is in Node, but you can port it to any language as long as you can run a web socket server, the example uses dgram.

### How to contribute

This project is open to any contributions that help improve it! Maybe you'd like to just port this to a new language like Typescript or Python.
You can also submit translations, help clean up the code or improve the example file.
All you need to do is Fork this project and make a new pull request!
