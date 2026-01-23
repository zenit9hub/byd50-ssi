# SSI PoC Demo PRD

## 목적
- SSI(Self-Sovereign Identity) 기본 개념을 학습하고, DID/VC/VP 흐름을 실습 가능한 데모 플랫폼 제공
- gRPC 기반 서비스를 통해 DID 발급/등록/해결, VC 발급, VP 제출/검증 시나리오를 체험
- 교육 과정에서 이론-실습-앱 데모로 이어지는 일관된 학습 동선을 확보

## 범위
- DID 생성/등록/해결(Registry/SEP)
- VC 발급(issuer)
- VP 생성 및 검증(relying party)
- 데모 클라이언트 시나리오 실행(3단계 use case)
- REST API 샘플(did_service_endpoint) 및 Swagger 문서

## 비범위(Out of Scope)
- 메인넷/실서비스 배포용 보안 하드닝
- 분산 KMS/탈중앙 키 관리
- 멀티테넌트/권한관리/감사 로깅
- 고가용성(HA), 자동 스케일링

## 대상 사용자
- 교육 참가자(개발자/기획자)
- SSI 개념 학습 및 실습 진행자
- 데모 앱 개발 담당자(안드로이드)

## 핵심 사용자 시나리오
1. DID 생성 및 레지스트리 등록
2. 발급자(issuer)로부터 VC 발급
3. RP(relying party)에 VP 제출 및 검증
4. 인증 챌린지/리스폰스 기반 로그인 인증

## 성공 지표
- 강의/실습 2일 과정 내 전체 시나리오 성공률 95% 이상
- 데모 클라이언트 3개 use case 모두 정상 실행
- 교육용 문서/슬라이드에서 시스템 구성과 흐름 설명 가능

## 주요 구성요소
- DID-Registry: DID Document 저장소(LevelDB 기반)
- DID-Registrar: DID method 기반 라우팅 및 Resolver
- demo-issuer: VC 발급 기관 역할
- demo-rp: VP 검증 및 챌린지 인증
- demo-client: 시나리오 실행 클라이언트
- did_service_endpoint: REST API 샘플 + Swagger
- pkg/did/kms: 키 생성/서명/복호화 제공

## 제약/가정
- 로컬 환경에서 5개 서비스 동시 구동
- gRPC 통신 유지(교육용), REST는 보조 채널
- 키는 .env에 저장(교육 시나리오 우선)

## 릴리스 목표(교육용)
- v1: 로컬 데모 + CLI 기반 실습 완주
- v1.1: 안드로이드 데모 앱 연동
