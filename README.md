# Casty API Server
This is a backend API server of Casty project written in go!

<a target="_blank" href="https://documenter.getpostman.com/view/471191/SzYT5246">
  <img src="https://img.shields.io/badge/Postman-api%20documentation-orange?logo=postman&style=for-the-badge" alt="Postman API Documentation">
</a>

You can find API Documentations on Postman

## Run Docker Container
```bash
$ docker run -p 9002:9002 castyapp/api:latest
```

## Requirements
* Golang `(1.15)` Always be up to date!) [Install Golang!](https://golang.org/doc/install)
* gRPC.server **This project needs to connect to Casty gRPC server!**  [More info](https://github.com/castyapp/grpc.server)

## Clone the project
```bash
$ git clone https://github.com/castyapp/api.server.git
```

## Configuraition
There is a `example.config.hcl` file that you should make a copy of, and call it `config.hcl` in your work directory.
```bash
$ cp example.config.hcl config.hcl
```

### Configure grpc client
You can find more information about how to run grpc server 
is available here [https://github.com/castyapp/grpc.server#readme]
```hcl
grpc {
  host = "localhost"
  port = 55283
}
```

### S3 bucket setup
You can configure the s3 bucket with these configurations
This works with minio too
```hcl
s3 {
  endpoint             = "127.0.0.1:9000"
  access_key           = "secret-access-key"
  secret_key           = "secret-key"
  use_https            = true
  insecure_skip_verify = true
}
```

### Recaptcha setup

`recaptcha.secret` is a secret key that you get on hcaptcha admin console!
for more information about hcaptcha [click here](https://www.hcaptcha.com/)

`recaptcha.type` is only available for hcaptcha, google will add soon!

```hcl
recaptcha {
  enabled = false
  type    = "hcaptcha"
  secret  = "hcaptcha-secret-token"
}
```

for more information about hcaptcha [click here](https://www.hcaptcha.com/)

You're ready to Go!

## Run project with go compiler
You can simply run the project with following command
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
