# TL;DR

Little nginx auth_proxy sidecar tests whether a request's client is authorized to perform some
action on some resource by asking Fence or Arborist for gating access to URLs from the reverse proxy - ex:
     http://localhost:7777/some-action/some-resource

