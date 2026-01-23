# Refactoring PoC → MVP 제안

## 1) 아키텍처/도메인 정리
- 패키지 경계 재정비: `did` 라이브러리(암호/문서/드라이버)와 서비스 레이어(Registry/Registrar/Issuer/RP/REST)를 명확히 분리하고, 공용 유틸(`dkms`, `byd50-jwt`)을 모듈화.
- Config 일원화: `configs/configs.yml`을 서비스별 구조체로 매핑하고, 환경변수 오버라이드/검증 로직을 추가해 배포 환경별 관리.
- 인터페이스 추상화: Registry 스토리지 인터페이스(put/get/update)와 DID Driver 인터페이스에 컨텍스트·옵션을 추가해 테스트/교체 용이성 확보.
- 서비스 역할 명확화: DID-Registrar는 method 라우팅/정책/검증 경계, DID-Registry는 저장/조회 책임으로 분리한다.

## 2) 보안/키 관리
- 하드코딩 키 제거: `eth.go`의 ECDSA 키 상수, PoC용 랜덤 키 생성 흐름을 내부 KMS/keystore 파일 기반으로 교체.
- 챌린지 관리 개선: `demo-rp`의 전역 `sourceData` → 요청별 nonce 저장소(메모리 캐시+TTL)로 변경, 재사용/리플레이 방지.
- 서명 알고리즘 검증: VC/VP 생성 시 kid·alg 일관성 확인, 지원 키 타입(ECDSA/RSA/Ed25519) 명시 및 거부 정책 추가.

## 3) 데이터/스토리지
- Registry 스토리지 교체: `/tmp/foo.db` LevelDB 대신 내구성 DB(PostgreSQL/Badger/Cloud KV)로 전환하고 마이그레이션 스크립트 추가.
- DID 메타데이터 관리: 문서 메타(`resolutionMetadata`) 저장/조회, 업데이트 이력(버전·타임스탬프·서명) 관리.
- 캐싱 전략: Resolver 측에 DID Document 캐시(TTL, ETag) 추가해 응답 지연/비용 감소.

## 4) 서비스 품질
- 에러 처리: gRPC/REST 전반에 표준 에러 코드와 메시지 정의, 로깅 레벨 구분 및 구조화 로깅 도입(zap/logrus).
- 관측성: healthz/readiness, Prometheus 메트릭, 요청 트레이싱(OTel) 추가. 주요 플로우(VC 발행, VP 검증) 대시보드화.
- 동시성/타임아웃: 모든 외부 호출에 `context` 기반 타임아웃/리트라이, Registrar/Registry 클라이언트 커넥션 풀 관리.

## 5) VC/VP 사양 정합성
- 스키마 정리: VC/VP 클레임을 JSON Schema/LD Context로 명시하고, 검증 로직에 스키마 검증 추가.
- 표준 프로필: W3C DID Core/VC Data Model 준수 확인, DID Document의 service/verificationMethod/authentication 필드 확장.
- 만료/검증 정책: VC 만료·발급자 신뢰 정책, VP 제출 시 Audience/Nonce 체크를 명문화.

## 6) API/UX 개선
- REST·gRPC 정렬: Registrar/Registry/Issuer/RP 모두에 HTTP 게이트웨이 제공(OpenAPI/Swagger)하여 외부 연계 용이화.
- 에러 메시지/응답 모델 정리: 인증·검증 실패 사유를 세분화하여 클라이언트 UX 개선.
- 샘플/SDK: Go/TypeScript 예제 SDK 추가하여 DID 생성·VC 발행·VP 제출을 손쉽게 호출하도록 정리.

## 7) DevOps/테스트
- 테스트 전략: 유닛(드라이버/암호/VC 검증), 통합(gRPC 플로우), 시나리오 테스트(VC 체인 발급→VP 검증) 추가. 로컬/CI에서 docker-compose로 서비스 묶음 구동.
- CI/CD: lint(fmt/vet/golangci-lint), 테스트, 빌드(멀티 플랫폼), 컨테이너 이미지 생성 파이프라인 구축.
- 릴리즈 관리: 버전 태그, 변경 로그, 환경별 배포 설정(dev/stage/prod) 분리.

## 8) 기능 로드맵(우선순위 제안)
1. 키/보안 정리: 하드코딩 키 제거, 챌린지 저장소 개선, 에러 코드 통합.
2. 스토리지/데이터 무결성: Registry DB 교체+마이그레이션, DID 문서 버전 관리.
3. VC/VP 표준화: 스키마·검증 정책 명문화, 테스트 케이스 확충.
4. 관측성/운영: health/metrics/tracing, 로깅 개선, CI 파이프라인 도입.
5. 확장성: HTTP 게이트웨이 정비, SDK/예제 제공, 드라이버 플러그인 구조 고도화.

## 9) 지금 단계 우선순위(3,4 실행 계획)
### 3) 서비스 실행 단일 진입점 정리
- `make dev-up` 표준화: 실행 순서/의존성(Registry→Registrar→Issuer→RP→Client) 고정.
- `.env` 기반 환경변수 표준화: 최소 필수 키만 유지하고 누락 시 명확한 에러 제공.
- 로그/상태 가시성: `.devlogs`에 서비스별 로그 분리, 실패 시 위치 안내.

### 4) 리팩토링 우선순위 선정
1. `did/pkg/controller`의 외부 의존 최소화: gRPC 클라이언트 주입 경계 유지 및 에러 처리 개선.
2. `did/core`의 DID/VC/VP 생성·검증 경로 정리: 입력 검증/에러 일관화.
3. 키 유틸과 KMS 계층 분리: RSA/ECDSA 키 변환·서명 로직 단위화.
4. Registry 스토리지 인터페이스 분리: 현재 LevelDB 의존을 인터페이스로 분리해 교체 가능하게 준비.

### KMS/키 유틸 구조 분리 초안
- 목표 패키지 경계:
  - `pkg/did/core/kms`: 내부 KMS 도메인(키 보관/생성/표현, 서비스 내부 사용).
  - `pkg/keys`: RSA/ECDSA 변환/서명/암복호화 유틸(순수 함수 중심).
- 단계적 이동 계획:
  1) 유틸 함수 그룹핑(Export/Parse/Encrypt/Sign) → `did/pkg/keys`로 이동.
  2) `dkms`는 `keys`만 의존하도록 정리, 외부에서 키 형식 직접 접근 최소화.
  3) controller/greeter는 `dkms`를 통해서만 키 사용(직접 유틸 호출 금지).

## MVP 고도화 체크리스트 (진행 기준)
- [ ] 키/보안: 하드코딩 키 제거 및 `.env`/키스토어 기반 설정 완료.
- [ ] 키/보안: 챌린지 nonce 저장소(TTL)로 리플레이 방지 적용.
- [ ] 운영성: `make dev-up`로 5개 서비스+클라이언트 실행 후 정상 로그 확인.
- [ ] 데이터: Registry 스토리지 인터페이스 분리(레벨DB 의존 제거).
- [ ] 검증: VC/VP 표준 클레임 검증(iss/aud/exp/nbf) 기본 정책 적용.
- [ ] 테스트: 단위/통합 스모크 테스트 스크립트 추가 및 CI 연동.
