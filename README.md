# BYD50 SSI

SSI PoC/MVP codebase with gRPC services, REST demo server, and Android app demo.

## Layout
- `apps/`: runnable services/clients
- `pkg/`: shared libraries (DID/VC/VP/KMS, etc.)
- `configs/`: service configuration
- `proto-files/`: gRPC schema
- `android/`: Android demo app (JNI + REST)
- `docs/`: documentation

## Quick Start (gRPC demo)
```bash
make dev-up
```
Starts:
- DID-Registry
- DID-Registrar
- demo-issuer
- demo-rp
- demo-client (scenario runner)

Logs: `.devlogs/*.log`  
Ports: `50051~50055`

## REST Demo Server (did_service_endpoint)
```bash
go run ./apps/did_service_endpoint/main.go
```
Swagger: `http://localhost:8080/swagger/index.html`

Key endpoints:
- DID: `/v2/testapi/create-did`, `/v2/testapi/get-did/:id`
- VC/VP: `/v2/testapi/vc/*`, `/v2/testapi/vp/*`
- Demo flow: `/v2/testapi/license/*`, `/v2/testapi/rental/*`

## Android Demo App
1) Build JNI libraries:
```bash
make android
```
2) Open `android/` in Android Studio and run.

Default base URL: `http://10.0.2.2:8080` (emulator → local REST server)

## Docs
Key docs under `docs/`:
- `0-1.강의-요약본.md`
- `0-2.교육-용어-시나리오-정리.md`
- `1-1.FBS.md`
- `1-2.PRD.md`
- `1-3.아키텍처.md`
- `2-1.개발문서.md`
- `2-2.실습-체크리스트.md`
- `2-3.코드-정리-지침.md`

## Testing
- `make test-summary` (verbose + summary)
- `make coverage-did` (coverage target)

## API Docs
- `make swagger` / `make swagger-lint` / `make swagger-docs`
- `make swagger-all` (spec → lint → html → pdf)

## Docker (optional)
If you use ECR, update the registry/account values:
```bash
aws ecr get-login-password --region ap-northeast-2 | docker login --username AWS --password-stdin <account>.dkr.ecr.ap-northeast-2.amazonaws.com
docker build -t did-registry -f ./apps/did-registry/Dockerfile .
docker build -t did-registrar -f ./apps/did-registrar/Dockerfile .
```
