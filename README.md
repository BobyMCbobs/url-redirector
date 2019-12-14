# url-redirector

> A simple golang + yaml URL redirector


## Definitions
Example of `./config.yaml`
```yaml
duck: https://duckduckgo.com
```

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

## Environment variables

| Name | Purpose |
| - | - |
| `APP_PORT` | the port and interface which the app serves from |
| `APP_CONFIG_YAML` | the location of where the config.yaml is (for the [routes](#definitions)) |
| `APP_LOG_FILE` | the location of where a log file will be created and written to |

## License
Copyright 2019 Caleb Woodbine.
This project is licensed under the [GPL-3.0](http://www.gnu.org/licenses/gpl-3.0.html) and is [Free Software](https://www.gnu.org/philosophy/free-sw.en.html).
This program comes with absolutely no warranty.
