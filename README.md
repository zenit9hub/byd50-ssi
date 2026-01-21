# Byd50-SSI

## Docker

### docker push

```bash
$ aws ecr get-login-password --region ap-northeast-2 | docker login --username AWS --password-stdin 086849521175.dkr.ecr.ap-northeast-2.amazonaws.com
$ docker build -t rs-did-registry -f ./did_registry/Dockerfile .
$ docker tag rs-did-registry:latest 086849521175.dkr.ecr.ap-northeast-2.amazonaws.com/rs-did-registry:latest
$ docker push 086849521175.dkr.ecr.ap-northeast-2.amazonaws.com/rs-did-registry:latest

$ docker build -t rs-did-registrar -f ./did_registrar/Dockerfile .
$ docker tag rs-did-registrar:latest 086849521175.dkr.ecr.ap-northeast-2.amazonaws.com/rs-did-registrar:latest
$ docker push 086849521175.dkr.ecr.ap-northeast-2.amazonaws.com/rs-did-registrar:latest
```

### docker run

```bash
$ aws ecr get-login-password --region ap-northeast-2 | docker login --username AWS --password-stdin 086849521175.dkr.ecr.ap-northeast-2.amazonaws.com

$ docker run --rm -d \
--name rs-did-registry \
-p 50051:50051 \
086849521175.dkr.ecr.ap-northeast-2.amazonaws.com/rs-did-registry:latest

$ docker run --rm -d \
--name rs-did-registrar \
-p 50052:50052 \
086849521175.dkr.ecr.ap-northeast-2.amazonaws.com/rs-did-registrar:latest
```
