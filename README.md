# Casty API Server
This is a backend API server of Casty project written in go!

<a target="_blank" href="https://documenter.getpostman.com/view/471191/SzYT5246">
  <img src="https://img.shields.io/badge/Postman-api%20documentation-orange?logo=postman&style=for-the-badge" alt="Postman API Documentation">
</a>

You can find API Documentations on Postman

## Requirements
* Golang `(1.14)` Always be up to date!) [Install Golang!](https://golang.org/doc/install)
* gRPC.server **This project needs to connect to Casty gRPC server!**  [More info](https://github.com/CastyLab/grpc.server)

## Clone the project
```bash
$ git clone https://github.com/CastyLab/api.server.git
```

## Configuraition
There is a `config.example.yml` file that you should make a copy of, and call it `config.yml` in your work directory.

```bash
$ cp config.example.yml config.yml
```

The most important environments here are `grpc.host`, `grpc.port` and `secrets.hcaptcha_secret`
```yaml
grpc:
  host: "localhost"
  port: 55283
```

`grpc.host` and `grpc.port` are the gRPC.server ip address and port that you should have! [Casty gRPC.server](https://github.com/CastyLab/grpc.server)

`secrets.hcaptcha_secret` is hcaptcha secret key that you should set up on google admin console!

for more information about hcaptcha [click here](https://www.hcaptcha.com/)

You're ready to Go!

## Run project with go compiler
you can simply run the project with following command
* this project uses go mod file, You can run this project out of the $GOPAH file!
```bash
$ go run server.go
```

or if you're considering building the project use
```bash
$ go build -o server .
```

## or Build/Run the docker image
```bash
$ docker build . --tag=casty.api

$ docker run -dp --restart=always 9002:9002 casty.api
```

## Contributing
Thank you for considering contributing to this project!
