DATETIME=$(shell date +"%m%d-%H%M" | tr ' :' '__')

ANDROID_OUT=./android/app/src/main/jniLibs
ANDROID_SDK=/Users/ryan9kim/Library/Android/sdk
NDK_BIN=/Users/ryan9kim/Library/Android/sdk/ndk/27.0.12077973/toolchains/llvm/prebuilt/darwin-x86_64/bin

dep:
	go get -u github.com/go-bindata/go-bindata/...
	go get -u github.com/golang/mock/mockgen/...
	go get -u github.com/jstemmer/go-junit-report
	go mod tidy
	go mod vendor

build:
	go build -o ./apps/did-registry ./apps/did-registry/main.go
	go build -o ./apps/did-registrar ./apps/did-registrar/main.go
	go build -o ./apps/demo-rp ./apps/demo-rp/main.go
	go build -o ./apps/demo-issuer ./apps/demo-issuer/main.go
	go build -o ./apps/demo-client ./apps/demo-client/main.go

docker:
	docker build -t did-registry_$(DATETIME) -f ./apps/did-registry/Dockerfile .
	docker build -t did-registrar_$(DATETIME) -f ./apps/did-registrar/Dockerfile .

proto:
	@PATH="$(shell go env GOPATH)/bin:$$PATH" protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto-files/*.proto

coverage:
	go test ./pkg/did/core/ -coverprofile="coverage.out"
	go tool cover -html="coverage.out" -o coverage.html
	go tool cover -html="coverage.out"

unit-test:
	@./scripts/coverage-did.sh

.PHONY: integration-test dev-down dev-status
integration-test:
	@./scripts/integration-test.sh

dev-down:
	@./scripts/dev-down.sh

dev-status:
	@./scripts/dev-status.sh

.PHONY: coverage-did swagger swagger-lint swagger-docs swagger-pdf
coverage-did:
	@./scripts/coverage-did.sh

swagger:
	@./scripts/swagger.sh

swagger-lint:
	@./scripts/swagger-lint.sh

swagger-docs:
	@./scripts/swagger-docs.sh

swagger-pdf:
	@./scripts/swagger-pdf.sh

.PHONY: swagger-all
swagger-all: swagger swagger-lint swagger-docs swagger-pdf

test2jenkins:
	go test -v ./did/core/ -tags="unit integration" -covermode=atomic -coverprofile=coverage.out ./cmd/... ./common/... 2>&1 | go-junit-report -set-exit-code > report.xml
	go tool cover -func coverage.out

android-armv7a:
	CGO_ENABLED=1 \
	GOOS=android \
	GOARCH=arm \
	GOARM=7 \
	CC=$(NDK_BIN)/armv7a-linux-androideabi21-clang \
	CGO_LDFLAGS="-Wl,-soname,libfoo.so" \
	go build -buildmode=c-shared -o $(ANDROID_OUT)/armeabi-v7a/libfoo.so ./pkg/did/c-shared/libfoo

android-arm64:
	CGO_ENABLED=1 \
	GOOS=android \
	GOARCH=arm64 \
	CC=$(NDK_BIN)/aarch64-linux-android21-clang \
	CGO_LDFLAGS="-Wl,-soname,libfoo.so" \
	go build -buildmode=c-shared -o $(ANDROID_OUT)/arm64-v8a/libfoo.so ./pkg/did/c-shared/libfoo

android-x86:
	CGO_ENABLED=1 \
	GOOS=android \
	GOARCH=386 \
	CC=$(NDK_BIN)/i686-linux-android21-clang \
	CGO_LDFLAGS="-Wl,-soname,libfoo.so" \
	go build -buildmode=c-shared -o $(ANDROID_OUT)/x86/libfoo.so ./pkg/did/c-shared/libfoo

android-x86_64:
	CGO_ENABLED=1 \
	GOOS=android \
	GOARCH=amd64 \
	CC=$(NDK_BIN)/x86_64-linux-android21-clang \
	CGO_LDFLAGS="-Wl,-soname,libfoo.so" \
	go build -buildmode=c-shared -o $(ANDROID_OUT)/x86_64/libfoo.so ./pkg/did/c-shared/libfoo


ps_android-armv7a:
	@echo "Building for android-arm..."
	@powershell -Command " \
	$$env:NDK_BIN='$(NDK_BIN)'; \
	$$env:CGO_ENABLED='1'; \
	$$env:GOOS='android'; \
	$$env:GOARCH='arm'; \
	$$env:CC=$$env:NDK_BIN + '/armv7a-linux-androideabi21-clang'; \
	$$env:CGO_LDFLAGS='-Wl,-soname,libfoo.so'; \
	go build -buildmode=c-shared -o $(ANDROID_OUT)/armeabi-v7a/libfoo.so ./pkg/did/c-shared/libfoo"

ps_android-arm64:
	@echo "Building for android-arm64..."
	@powershell -Command " \
	$$env:NDK_BIN='$(NDK_BIN)'; \
	$$env:CGO_ENABLED='1'; \
	$$env:GOOS='android'; \
	$$env:GOARCH='arm64'; \
	$$env:CC=$$env:NDK_BIN + '/aarch64-linux-android21-clang'; \
	$$env:CGO_LDFLAGS='-Wl,-soname,libfoo.so'; \
	go build -buildmode=c-shared -o $(ANDROID_OUT)/arm64-v8a/libfoo.so ./pkg/did/c-shared/libfoo"

ps_android-x86:
	@echo "Building for android-386..."
	@powershell -Command " \
	$$env:NDK_BIN='$(NDK_BIN)'; \
	$$env:CGO_ENABLED='1'; \
	$$env:GOOS='android'; \
	$$env:GOARCH='386'; \
	$$env:CC=$$env:NDK_BIN + '/i686-linux-android21-clang'; \
	$$env:CGO_LDFLAGS='-Wl,-soname,libfoo.so'; \
	go build -buildmode=c-shared -o $(ANDROID_OUT)/x86/libfoo.so ./pkg/did/c-shared/libfoo"

ps_android-x86_64:
	@echo "Building for android-amd64..."
	@powershell -Command " \
	$$env:NDK_BIN='$(NDK_BIN)'; \
	$$env:CGO_ENABLED='1'; \
	$$env:GOOS='android'; \
	$$env:GOARCH='amd64'; \
	$$env:CC=$$env:NDK_BIN + '/x86_64-linux-android21-clang'; \
	$$env:CGO_LDFLAGS='-Wl,-soname,libfoo.so'; \
	go build -buildmode=c-shared -o $(ANDROID_OUT)/x86_64/libfoo.so ./pkg/did/c-shared/libfoo"

android: android-armv7a android-arm64 android-x86 android-x86_64

ps_android: ps_android-armv7a ps_android-arm64 ps_android-x86 ps_android-x86_64

# For CI

# Login AWS ECR
aws_ecr:
	#aws ecr get-login-password --region ap-northeast-2 | docker login --username AWS --password-stdin 086849521175.dkr.ecr.ap-northeast-2.amazonaws.com
	docker run --rm --name aws \
		-e AWS_ACCESS_KEY_ID="${AWS_ACCESS_KEY_ID}" \
		-e AWS_SECRET_ACCESS_KEY="${AWS_SECRET_ACCESS_KEY}" \
		-e AWS_DEFAULT_REGION="${AWS_DEFAULT_REGION}" \
		-v "${PWD}:/aws" \
		amazon/aws-cli ecr get-login-password --region ap-northeast-2 | docker login --username AWS --password-stdin 086849521175.dkr.ecr.ap-northeast-2.amazonaws.com

# Build for AWS ECR Destribution
docker_build_registry:
	docker build -t did-registry -f ./apps/did-registry/Dockerfile .
	docker tag did-registry:latest 086849521175.dkr.ecr.ap-northeast-2.amazonaws.com/did-registry:latest
docker_build_registrar:
	docker build -t did-registrar -f ./apps/did-registrar/Dockerfile .
	docker tag did-registrar:latest 086849521175.dkr.ecr.ap-northeast-2.amazonaws.com/did-registrar:latest

# Push to AWS ECR
docker_push_registry:
	docker push 086849521175.dkr.ecr.ap-northeast-2.amazonaws.com/did-registry:latest
docker_push_registrar:
	docker push 086849521175.dkr.ecr.ap-northeast-2.amazonaws.com/did-registrar:latest

.PHONY: dev-up
dev-up:
	@./scripts/dev-up.sh

.PHONY: test-summary
test-summary:
	@./scripts/test-summary.sh
