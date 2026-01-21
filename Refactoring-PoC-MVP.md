# Refactoring PoC → MVP 제안

## 1) 아키텍처/도메인 정리
- 패키지 경계 재정비: `did` 라이브러리(암호/문서/드라이버)와 서비스 레이어(Registry/Registrar/Issuer/RP/REST)를 명확히 분리하고, 공용 유틸(`dkms`, `byd50-jwt`)을 모듈화.
- Config 일원화: `configs.yml`을 서비스별 구조체로 매핑하고, 환경변수 오버라이드/검증 로직을 추가해 배포 환경별 관리.
- 인터페이스 추상화: Registry 스토리지 인터페이스(put/get/update)와 DID Driver 인터페이스에 컨텍스트·옵션을 추가해 테스트/교체 용이성 확보.

## 2) 보안/키 관리
- 하드코딩 키 제거: `eth.go`의 ECDSA 키 상수, PoC용 랜덤 키 생성 흐름을 Secret Manager/KMS나 keystore 파일 기반으로 교체.
- 챌린지 관리 개선: `greeter_server`의 전역 `sourceData` → 요청별 nonce 저장소(메모리 캐시+TTL)로 변경, 재사용/리플레이 방지.
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

## MVP 고도화 체크리스트 (진행 기준)
- [ ] 키/보안: 하드코딩 키 제거 및 `.env`/키스토어 기반 설정 완료.
- [ ] 키/보안: 챌린지 nonce 저장소(TTL)로 리플레이 방지 적용.
- [ ] 운영성: `make dev-up`로 5개 서비스+클라이언트 실행 후 정상 로그 확인.
- [ ] 데이터: Registry 스토리지 인터페이스 분리(레벨DB 의존 제거).
- [ ] 검증: VC/VP 표준 클레임 검증(iss/aud/exp/nbf) 기본 정책 적용.
- [ ] 테스트: 단위/통합 스모크 테스트 스크립트 추가 및 CI 연동.
