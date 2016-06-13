# aws_proxy
Accessing some URLs (e.g. kibana, es) require signed URLs. This is an attempt to proxy and sign the URLs

AWS requires that some requests be signed. It is convoluted. So this is a simple attempt to just proxy and sign requests. 

It works for some URLs (annoying right?) but not all. Pretty much I can curl endpoints through the proxy, but can't make full web requests. I have tried stripping out the different headers but that causes things to break (e.g. kibana). 

[signing using v4](https://docs.aws.amazon.com/general/latest/gr/sigv4_signing.html) 

to start: 

```
  $ glide install
```

```
  $ go run main.go -c config.json  
```

If you don't specify cred it will try to load them from the env. 
