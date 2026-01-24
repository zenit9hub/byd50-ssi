# 리포지토리 분석 – 도메인별 정의 및 기능 상세

## 전체 개요
- 목적: SSI(Self-Sovereign Identity) PoC 아키텍처. gRPC 기반으로 DID 발급/해결, VC 발행, VP 검증 플로우를 시연.
- 핵심 구성: `pkg/did` 라이브러리(암호/문서/드라이버), `did-registry`(저장소), `did-registrar`(리졸버), `apps/did_service_endpoint`(REST 게이트웨이), `demo-*` 데모 클라이언트·릴라잉파티·발급자, `apps/geth_client`(체인 연동 예제), `proto-files`(gRPC 스키마).

## 공용 라이브러리(`pkg/did/`)
- `configs`  
  - `configs/configs.yml` 로드 후 `UseConfig`로 서비스 포트, 채택 드라이버 리스트, 테스트넷 URL/SC 주소 등을 제공.
- `core`  
  - `dids`: DID 생성/문서 초기화(`CreateDID`, `initDocument`), 생성 규칙(`hexdigit`/`uuid`/`base58`) 적용.  
  - `driver`: DID 메서드 드라이버 레지스트리.  
    - `byd50`: gRPC로 `did-registry`에 DID 생성/해결 요청.  
    - `eth`: 테스트넷 RPC, 배포된 컨트랙트 바인딩(`scdid`)으로 DID 생성/해결.  
    - `did_method.go`: 드라이버 등록/조회 인터페이스 정의.  
  - `rc`: `did-registrar` gRPC 클라이언트 싱글턴.  
  - `algorithm`: RSA 기반 암·복호화, 서명/검증, 난수 생성 유틸.  
  - `vc.go` / `vp.go`: VC/VP JWT 생성·검증 래퍼.  
  - `byd50-jwt`: VC/VP용 JWT 클레임 빌더 및 검증 로직.  
  - `service`: 향후 REST 서비스용 스텁.  
  - `byd50-jsonld`: 현재 비어 있는 JSON-LD 확장용 위치.
- `kms`  
  - RSA/ECDSA 키 생성·내보내기(Base58/PEM) 및 DID 연계 관리(내부 KMS).
- `registry`  
  - Store 인터페이스와 LevelDB 구현(`NewLevelDBStore`, `Put/Get/Has`).
- `pkg`  
  - `controller`: DID 생성/해결, 인증 챌린지/리스폰스, SimplePresent/VP 생성·검증을 `did-registrar`와 연계해 제공.  
  - `database`: LevelDB 초기화(`LEVELDB_PATH` 환경변수 기반).  
  - `logger`: 함수 시작/종료 로거.
- `keys`  
  - RSA/ECDSA 키 변환/서명/암복호화 유틸.

## DID Registry 서버(`apps/did-registry/`)
- 역할: PoC용 DID Document 저장소. LevelDB에 DID→문서 바이트 저장.
- gRPC 인터페이스(`proto-files/registry.proto` 기반):  
  - `CreateDid`: 입력 공개키로 DID/문서 생성(`dids.CreateDID`) 후 저장.  
  - `ResolveDid`: DID로 문서 조회, 없으면 `NotFound` 에러 문자열.  
  - `UpdateDid`: 존재 여부 확인 후 문서 업데이트(검증 로직 미구현).
- 구성: `configs.UseConfig.DidRegistryPort`에서 리스닝, 서버 시작/종료 시 DB 열고 닫음.

## DID Registrar 서버(`apps/did-registrar/`)
- 역할: 메서드별 드라이버 라우팅/추상화. DID 생성/해결 요청을 적합한 드라이버로 위임.
- 흐름:  
  - `CreateDID`: 요청 메서드(`byd50` 기본값) 기준 드라이버 선택→`CreateDid` 호출.  
  - `ResolveDID`: 입력 DID 파싱(`did:<method>:...`), 채택 드라이버 리스트 검증 후 드라이버 `ResolveDid` 실행.  
  - `UpdateDID`: 스텁 상태.
- gRPC 인터페이스: `proto-files/registrar.proto`. 포트 `UseConfig.DidRegistrarPort`.

## REST 서비스 엔드포인트(`apps/did_service_endpoint/`)
- 역할: gRPC 사용이 어려운 환경을 위한 간단한 HTTP 게이트웨이(Swagger 문서 포함).
- 엔드포인트:  
  - `POST /v2/testapi/create-did`: 메서드·공개키(Base58) 입력으로 DID 생성 후 반환.  
  - `GET /v2/testapi/get-did/:did`: DID Document 조회.  
  - `GET /v2/testapi/get-did-public-key/:did`: DID Document의 공개키 추출.  
  - `POST /v2/testapi/vc/create`: VC JWT 생성.  
  - `POST /v2/testapi/vc/verify`: VC JWT 검증.  
  - `POST /v2/testapi/vp/create`: VP JWT 생성.  
  - `POST /v2/testapi/vp/verify`: VP JWT 검증.
  - 데모 플로우:  
    - `GET /v2/testapi/demo/actors`: 발급기관/렌터카업체 DID 조회.  
    - `POST /v2/testapi/license/challenge`, `POST /v2/testapi/license/issue`  
    - `POST /v2/testapi/rental/challenge`, `POST /v2/testapi/rental/issue`
- 내부: Gin 서버, `controller`/`core`를 통해 DID/VC/VP 처리.
- 산출물: Swagger/Redoc 문서 `api-docs/`.

## Android 데모 앱(`android/`)
- 역할: REST API 호출 + JNI(c-shared) 기반 키 생성/서명으로 데모 시나리오 실행.
- 흐름: DID 생성 → 신원 VC → 계약 VC → VP 검증(차량 접근).

## Demo 데모 세트
- 공통: `proto-files/relyingparty.proto`·`issuer.proto` 기반 gRPC. PoC 시나리오용 예제 코드.
- `demo-client`:  
  - KMS 초기화(RSA→ECDSA 순서), `controller.CreateDID`로 DID 발급.  
  - Use case 1: Relying party `AuthChallenge` 수신→개인키 복호화 후 `AuthResponse`.  
  - Use case 2: SimplePresent(서명+타임스탬프) 생성/검증.  
  - Use case 3: VC 요청→발급 VC로 VP 구성→Relying party `VerifyVp` 호출.
- `demo-rp`(Relying Party):  
  - `AuthChallenge`: 난수/타임스탬프를 base58+평문으로 구성 후 공개키 암호화 문자열 반환.  
  - `AuthResponse`: 수신 문자열과 기존 챌린지 비교.  
  - `SimplePresent`: 서명 검증 및 만료(10초) 확인.  
  - `VerifyVp`: `core.VerifyVp`로 VP 검증.
- `demo-issuer`(Issuer):  
  - 서버 시작 시 ECDSA 키 생성→DID 발급.  
  - `RequestCredential`: 클라이언트 VP 클레임 검증 후 새 VC 발급.  
  - `ReqCredIdCard`, `ReqCredDlCard`, `ReqCredRentalCarAgreement`: 체인드 검증(이전 VC/VP 검증 후 다음 VC 발급) 및 최종 `RentalCarControl` 액세스 제어.  
  - VC 만료시간이 짧게 설정(1~3분/15초)된 PoC 예시.

## 체인 연동 예제(`apps/geth_client/`)
- 목적: EVM 테스트넷 RPC 연결, 컨트랙트 바인딩(`scdid`) 사용 예시.
- 기능: 컨트랙트 `ResolveDid` 호출, 새 DID/문서를 체인에 `CreateDid` 트랜잭션으로 전송.
- 설정: `ETH_RPC_URL`, `ETH_CHAIN_ID`, `ETH_REGISTRY_ADDRESS`, `ETH_PRIVATE_KEY_HEX` 환경변수 기반.

## 프로토콜 정의(`proto-files/`)
- `registry.proto`, `registrar.proto`, `issuer.proto`, `relyingparty.proto` 및 생성된 gRPC 바인딩.
- 각 서비스가 사용하는 메시지 스키마(Challenge/Response, VC/VP 전달, DID CRUD 등)를 중앙 관리.
