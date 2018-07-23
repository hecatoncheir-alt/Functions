# Mvideo page parser

Do not forget change image in _**mvideo-page-parser.yaml**_:

from
```
image: mvideo-page-parser
```
to 
```
image: some-repository/mvideo-page-parser
```

## For build user [faas-cli](https://github.com/openfaas/faas-cli):


In **_Dockerfile_** version of Go can be changed to: 
```
FROM golang:1.10.3-alpine3.8 as builder
```

```
faas-cli build -f .\mvideo-page-parser.yml

cd .\build\mvideo-page-parser\

docker build . -t some-repository/mvideo-page-parser
```

Then push image to docker registry:
```
docker push some-repository/mvideo-page-parser
```

## For deploy call faas-cli deploy.
Use **--gateway** if you have gateway on another server:

```
faas-cli deploy -f .\mvideo-page-parser.yml --gateway http://192.168.99.100:31112
```
