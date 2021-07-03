# Debug mode
debug = false

# Application environment
env = "dev"

# Configure grpc client
#
# You can find more information about how to run grpc server 
# is available here [https://github.com/castyapp/grpc.server#readme]
grpc {
  host = "localhost"
  port = 55283
}

http "rules" {
  access_control_allow_origin = "*"
}

# S3 bucket config
s3 {
  endpoint             = "127.0.0.1:9000"
  access_key           = "secret-access-key"
  secret_key           = "secret-key"
  use_https            = false
  insecure_skip_verify = true
}

# Sentry config
sentry {
  enabled = false
  dsn     = "sentry.dsn.here"
}

# Recaptcha config, it can be google or hcaptcha
recaptcha {
  enabled = false
  type    = "hcaptcha"
  secret  = "hcaptcha-secret-token"
}
