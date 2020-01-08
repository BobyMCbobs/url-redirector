# Examples

> Configuration ideas for url-redirector

## Any host redirection
Set a root property of `'*'` (wildcard) to redirect on any host.

``` yaml
'*':
  routes:
    a: https://about.gitlab.com
    b: https://github.com
    c: https://duckduckgo.com
  root: https://reddit.com
  wildcard: https://github.com
```

### Notes
- A config similar to this may be useful if you don't mind which host uses these values
- An issue with this config is if a third-party assigns a DNS entry to the this host, they can use the service

## Any host any path redirection
Completely redirect any request to a given URL.

```yaml
'*':
  root: https://github.com
  wildcard: https://github.com
```

### Notes
- A config similar to this may be useful if you have a service which is being renamed or deprecated in order to tell users about the new site, etc...
- An alternative is to set up a DNS entry to a host to redirect
