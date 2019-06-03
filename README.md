# Proxlet

Proxlet is a POC, naïve implemetation of serverless reverse proxy for any website.

## Supported Providers

* [Now.sh](https://zeit.co/now)
* [AWS Lambda](https://aws.amazon.com/lambda)
* [Google Cloud Functions](https://cloud.google.com/functions)

### Now.sh

dev dependency: [now-cli](https://github.com/zeit/now-cli#usage)

```bash
now

curl https://{app_name}.now.sh/https://httpbin.org/get
```

### AWS Lambda

dev dependency: [Serverless Framework](https://serverless.com/framework/docs/getting-started/)

```bash
# First customize serverless.yml's fields: provider.role, provider.region, etc.
make aws

# Then go to AWS API Gateway console, Setting, add Binary Media Type "*/*".
# (All contents are base64-encoded and will be decoded by the API Gateway.)

curl https://{id}.execute-api.{region}.amazonaws.com/dev/https://httpbin.org/get
```

CAVEAT: All subsequent requests fail because of they are lack of the stage path `/dev`

### Google Cloud Functions

dev dependency: [gcloud](https://cloud.google.com/sdk/docs/quickstarts)

```bash
make gcf

curl https://{region}-{project-id}.cloudfunctions.net/proxlet/https://httpbin.org/get
```

CAVEAT: All subsequent requests fail because of they are lack of the function path `/proxlet`

## About the Cookie

To get the real host of subsequent request of relative path, Proxlet sets a session cookie `proxlet-host` with the value of real host. The server does not store the cookie.

## Known Issues

* Cross-site requests are not proxied.
* Now.sh and AWS Lambda have an HTTP response size limit of 6MB. For GCF, it's 10MB.
* AWS Lambda's stage path / GCF's function name path makes all subsequent requests impossible *(unless using a custom domain)*.
