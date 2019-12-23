# url-redirector

> A simple cloud-native golang + yaml URL redirector app

## Config definitions

Example of `./config.yaml`
```yaml
routes:
  duck: https://duckduckgo.com
  gitlab: https://gitlab.com
  github: https://github.com
```

Given the config above, if [`https://localhost:8080/duck`](https://localhost:8080/duck) is visited, it will redirect to [`https://duckduckgo.com`](https://duckduckgo.com).

## Local usage
```bash
docker run -it --rm -v "$PWD"/config.yaml:/opt/redirector/config.yaml:z,ro -p 8080:8080 registry.gitlab.com/bobymcbobs/url-redirector
```

## Building
```bash
docker build -t registry.gitlab.com/bobymcbobs/url-redirector .
```

## Deployment in k8s
Make sure you update the values in the yaml files
```bash
kubectl apply -f k8s-deployment/
```
### Notes
- the ConfigMap can be updated at any time, give at least 1min for the patched version to be active

## Environment variables

| Name | Purpose | Defaults |
| - | - | - |
| `APP_PORT` | the port and interface which the app serves from | `:8080` |
| `APP_CONFIG_YAML` | the location of where the config.yaml is (for the [routes](#definitions)) | `./config.yaml` |
| `APP_USE_LOGGING` | toggle the saving of logs | `false` |
| `APP_LOG_FILE` | the location of where a log file will be created and written to | `./redirector.log` |

## License
Copyright 2019 Caleb Woodbine.
This project is licensed under the [GPL-3.0](http://www.gnu.org/licenses/gpl-3.0.html) and is [Free Software](https://www.gnu.org/philosophy/free-sw.en.html).
This program comes with absolutely no warranty.
