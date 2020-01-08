<a href="http://www.gnu.org/licenses/gpl-3.0.html">
    <img src="https://img.shields.io/badge/License-GPL%20v3-blue.svg" alt="License" />
</a>
<a href="https://gitlab.com/flattrack/flattrack.io/releases">
    <img src="https://img.shields.io/badge/version-1.3.0-brightgreen.svg" />
</a>

# url-redirector

> A simple cloud-native golang + yaml URL redirector app

## Features
- multi-host routes
- wildcard hosts
- root redirects
- wildcard paths
- simple YAML configuration
- restart-free configuration updating (only redeploy the ConfigMap)
- TLS/SSL support
- \>3MB container image

## Config definitions

Example of `./config.yaml`
```yaml
'*':
  routes:
    a: https://about.gitlab.com
    b: https://github.com
    c: https://duckduckgo.com
  root: https://reddit.com
  wildcard: https://github.com

myshortner1.com:
  routes:
    a: https://duckduckgo.com
    b: https://reddit.com
    c: https://github.com
  root: https://about.gitlab.com
  wildcard: https://github.com

myshortner2.com:
  routes:
    a: https://about.gitlab.com
    b: https://duckduckgo.com
    c: https://github.com
  root: https://reddit.com
  wildcard: https://github.com
```

Given the config above, if [`https://localhost:8080/duck`](https://localhost:8080/duck) is visited, the request will redirect to [`https://duckduckgo.com`](https://duckduckgo.com). If [`https://localhost:8080`](https://localhost:8080) is visited, the request will be redirected to [`https://gitlab.com`](https://gitlab.com). If the path that doesn't exist is visited, the request will be redirected to [`https://github.com`](https://github.com).  

For more examples, check out [docs/EXAMPLES.md](docs/EXAMPLES.md)

## Local usage
```bash
docker run -it --rm -v "$PWD"/config.yaml:/app/config.yaml:z,ro -p 8080:8080 registry.gitlab.com/bobymcbobs/url-redirector:latest
```

## Building
```bash
docker build -t registry.gitlab.com/bobymcbobs/url-redirector:latest .
```

## Deployment in k8s
Make sure you update the values in the yaml files
```bash
kubectl apply -f k8s-deployment/
```

### Notes
- the ConfigMap can be updated at any time, give a few seconds for the patched version to be active

## Environment variables

| Name | Purpose | Defaults |
| - | - | - |
| `APP_PORT` | the port and interface which the app serves from | `:8080` |
| `APP_PORT_TLS` | the port and interface which the app serves from | `:4433` |
| `APP_USE_TLS` | run the app with TLS enabled | `false` |
| `APP_TLS_PUBLIC_CERT` | the public certificate for the app to use | `server.crt` |
| `APP_TLS_PRIVATE_CERT` | the private cert for the app to use | `server.tls` |
| `APP_CONFIG_YAML` | the location of where the config.yaml is (for the [routes](#definitions)) | `./config.yaml` |
| `APP_USE_LOGGING` | toggle the saving of logs | `false` |
| `APP_LOG_FILE` | the location of where a log file will be created and written to | `./redirector.log` |

## License
Copyright 2019 Caleb Woodbine.
This project is licensed under the [GPL-3.0](http://www.gnu.org/licenses/gpl-3.0.html) and is [Free Software](https://www.gnu.org/philosophy/free-sw.en.html).
This program comes with absolutely no warranty.
