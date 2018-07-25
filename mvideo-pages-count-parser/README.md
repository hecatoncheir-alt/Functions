# Mvideo pages count parser

Do not forget change image in _**mvideo-pages-count-parser.yaml**_:

from
```
image: mvideo-pages-count-parser
```
to 
```
image: some-repository/mvideo-pages-count-paresr
```

## For build user [faas-cli](https://github.com/openfaas/faas-cli):


In **_Dockerfile_** version of Go can be changed to: 
```
FROM golang:1.10.3-alpine3.8 as builder
```

```
faas-cli build -f .\mvideo-pages-count-parser.yml

cd .\build\mvideo-pages-count-parser\

docker build . -t some-repository/mvideo-pages-count-parser
```

Then push image to docker registry:
```
docker push some-repository/mvideo-pages-count-parser
```

## For deploy call faas-cli deploy.
Use **--gateway** if you have gateway on another server:

```
faas-cli deploy -f .\mvideo-pages-count-parser.yml --gateway http://192.168.99.100:31112
```
