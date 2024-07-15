# HELM GO Project

Upload and parse CSV files(LF format).

## How to use

1. `docker build . -t eugst/tst-app:latest`
2. `docker push eugst/tst-app:latest`
3. `helm upgrade --install test ./helm/`

## Directory Description

`kops` consists configuration for creating cluster with multi-az masters and different instance group by instances types and costs.

`helm` consists helm chart.

`static` - directory for static files served by nginx.

`templates` - for template files used in go.
