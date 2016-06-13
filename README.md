# aws_proxy
Accessing some URLs (e.g. kibana, es) require signed URLs. This is an attempt to proxy and sign the URLs

AWS requires that some requests be signed. It is convoluted. So this is a simple attempt to just proxy and sign requests. 

This works, but signing doesn't seem to support keep alive.

[signing using v4](https://docs.aws.amazon.com/general/latest/gr/sigv4_signing.html) 

to start: 

```
  $ glide install
```

```
  $ go run main.go -c config.json  
```

If you don't specify cred it will try to load them from the env. 
