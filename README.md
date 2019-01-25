# TL;DR

Little nginx auth_proxy sidecar tests whether a request's client is 
authorized to perform some
action on some resource by asking Fence or Arborist for gating access to URLs from the reverse proxy - ex:
     http://localhost:7777/some-action/some-resource


## Use cases

* Gate access to thirs party tools like prometheus and jupyterhub
* Arborist-enable legacy services like indexd
* Example (still needs end to end testing) `nginx.conf` with http://nginx.org/en/docs/http/ngx_http_auth_request_module.html
```
  location /prometheus/ {
      auth_request http://localhost:7780/prometheus-admin/prometheus;
      ...
  }
```

## Details

Currently just calls out to http://fence/user , and returns true (access allowed)
if the user is an admin.  We can drop in an Arborist implementation of the `AuthzService` interface next ...

## Build, Test, Run

```
mkdir -p ~/go/src/github.com/uc-cdis
cd ~/go/src/github.com/uc-cdis
git clone git@github.com:uc-cdis/authproxy.git
cd authproxy

(cd authProxyServer && go build)
go test -v
bash testWithToken.sh you@dev.csoc you.planx-pla.net you@uchicago.edu

authProxyServer/authProxyServer
```
