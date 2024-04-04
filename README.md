[![Build Status](https://aggmet-ci.lma.wbx2.com/api/badges/aggmet/docker-iam-rev-proxy/status.svg)](https://aggmet-ci.lma.wbx2.com/aggmet/docker-iam-rev-proxy)


Generate the AWS Creds

```
aws sts assume-role --role-arn AWS_IAM_ROLE_TO_ASSUME --role-session-name "YOUR_SESSION_NAME" --profile YOUR_AWS_PROFILE_NAME > assume-role-output.json
```

Make a request and pass in the headers

```
{
  'access-key': 'YOUR_ACCESS_KEY_ID',
  'secret-access-key': 'YOUR_SECRET_ACCESS_KEY_ID',
  'session-token': 'YOUR_SESSION_TOKEN',
}
```

# Example

### Update the env vars for your use case

```
cp .env.example .env

# update the needed ENV vars
IAM_PROXY_BACKEND_URL=http://grafana-iam-proxy:3000
IAM_PROXY_AWS_REGION=YOUR_AWS_REGION
IAM_PROXY_IAM_CONFIG_LOCATION=LOCATION_OF_YOUR_IAM_YAML_CONFIG
IAM_PROXY_CACHE_ENABLED=true
IAM_PROXY_CACHE_TIMEOUT_IN_MINUTES=60
IAM_PROXY_CACHE_CAPACITY=1000
IAM_PROXY_DEBUG=true
```

### Create the iam-proxy.yaml

```
cp example-iam-proxy.yaml iam-proxy.yaml
```

```yaml
accessKeyHeaderKey: "access-key"
secretAccessKeyHeaderKey: "secret-access-key"
sessionTokenHeaderKey: "session-token"
headerAttributeMapping:
  "X-WEBAUTH-USER": "Arn"              # can be Arn, UserId, or Account (default is Arn)
headerAttributeAdditional:
  "FOO": "Bar"                         # be able to set some extra headers
  "FOOENV": "${FOOENV}"                # load in some env var
allowedArns:
  - "FULL_ARN_ALLOWED"
```

### Customize endpoints

You can change the default endpoints for the main proxy, prometheus metrics, and health check

```
IAM_PROXY_ENDPOINT=/
IAM_PROXY_METRICS_ENDPOINT=/metrics
IAM_PROXY_HEALTH_CHECK_ENDPOINT=/_health
```

### Start grafana, postgres, and proxy

```
# start the docker compose file
make start

# run the tests, it will create the assume role file with credentials and make a curl request with the credentials
make test
```

