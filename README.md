# S3bucket Exporter

s3bucket_exporter collects informations about size and object list about all the buckets accessible by user. Was designed to work with ceph, but should work will all S3 compatible endpoints.

## Getting started

Run from command-line:

```
./s3bucket_exporter [flags]
```

Run from command-line - example with minimal parameter list:
```
./s3bucket_exporter --s3_endpoint=192.168.0.1:7480 --s3_access_key=akces123 --s3_secret_key=secret123 --s3_name=MyS3Endpoint
```

The exporter supports two different configuration ways: command-line arguments take precedence over environment variables.

As for available flags and equivalent environment variables, here is a list:

|     environment variable          |    argument                       |     description                                    |        default        |     example              |
| --------------------------------- | --------------------------------- | -------------------------------------------------- | --------------------- | ------------------------ |
| S3_NAME                           | --s3_name                         | S3 configuration name, visible as a tag            |                       | MyS3Endpoint             |
| S3_ENDPOINT                       | --s3_endpoint                     | S3 endpoint url with port  |                       |                       | 192.168.0.1:7480         |
| S3_ACCESS_KEY                     | --s3_access_key                   | S3 access_key (aws_access_key)                     |                       | myAkcesKey               |
| S3_SECRET_KEY                     | --s3_secret_key                   | S3 secret key (aws_secret_key)                     |                       | mySecretKey              |
| S3_REGION                         | --s3_region                       | S3 region name                                     | default               | "default" or "eu-west-1" |
| LISTEN_PORT                       | --listen_port                     | Exporter listen Port cluster                       | :9655                 | :9123                   |
| LOG_LEVEL                         | --log_level                       | Log level. Info or Debug                           | Info                  | Debug                    |
| S3_DISABLE_SSL                    | --s3_disable_ssl                  | If S3 endpoint is not secured by SSL set to True   | False                 | True                     |
| S3_DISABLE_ENDPOINT_HOST_PREFIX   | --s3_disable_endpoint_host_prefix | Disable endpoint host prefix                       | False                 | True                     |
| S3_FORCE_PATH_STYLE               | --s3_force_path_style             | Force use path style (bucketname not added to url) | True                  | False                    |

> Warning: For security reason is not advised to use credential from command line

## Prometheus configuration example:

```yaml
  - job_name: 's3bucket'
    static_configs:
    - targets: ['192.168.0.5:9655']
```

## Grafana
Grafana dashboad ([resources/grafana-s3bucket-dashboard.json] (resources/grafana-s3bucket-dashboard.json)):

![](images/grafana-s3bucket-dashboard.png)

