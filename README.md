# Franz

## Introduction

Franz is heavily inspired by Apache Kafka, which is a distributed event streaming platform.

### Why am I building this?

I first encountered Apache Kafka from Fireship.io's video, I did not fully understanding what Kafka is, what it does, and why was it built. Hence, I am building my own version of Kafka with some twist to obtain a better understanding of Kafka and what it means to be a distributed event streaming platform.

**Note: Below are from codecrafters.io which I am using to get a structured path to building the finished product. Also note that only the starter code and codecrafters config/build file/folders is from codecrafters.io, everything else is my own code.**

## Run the code

### Pre-requisite

- You should have latest version of [Go](https://go.dev/doc/install) installed
- You should be using an Unix-based machine to run the code (I will support Docker image in the future)

```bash
chmod u+x your_program.sh # if needed
./your_program.sh
```


[![progress-banner](https://backend.codecrafters.io/progress/kafka/dc080a07-9161-4c66-94bd-bdb7078cfc0a)](https://app.codecrafters.io/users/codecrafters-bot?r=2qF)

This is a starting point for Go solutions to the
["Build Your Own Kafka" Challenge](https://codecrafters.io/challenges/kafka).

In this challenge, you'll build a toy Kafka clone that's capable of accepting
and responding to APIVersions & Fetch API requests. You'll also learn about
encoding and decoding messages using the Kafka wire protocol. You'll also learn
about handling the network protocol, event loops, TCP sockets and more.

**Note**: If you're viewing this repo on GitHub, head over to
[codecrafters.io](https://codecrafters.io) to try the challenge.

# Passing the first stage

The entry point for your Kafka implementation is in `app/server.go`. Study and
uncomment the relevant code, and push your changes to pass the first stage:

```sh
git commit -am "pass 1st stage" # any msg
git push origin master
```

That's all!

# Stage 2 & beyond

Note: This section is for stages 2 and beyond.

1. Ensure you have `go (1.19)` installed locally
1. Run `./your_program.sh` to run your Kafka broker, which is implemented in
   `app/server.go`.
1. Commit your changes and run `git push origin master` to submit your solution
   to CodeCrafters. Test output will be streamed to your terminal.
