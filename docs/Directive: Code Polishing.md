# [Directive: Code Polishing]

작업을 수행할 때, **'수정하는 코드와 그 주변(Context)'**을 내가 발견했을 때보다 더 빛나고 매끄럽게 다듬으세요. (Boy Scout Rule)

## ✅ Polishing 대상 (Checklist)
1.  **Naming**: 모호한 이름(`data`, `flag`)을 의도가 드러나는 이름(`userProfile`, `isLoaded`)으로 변경.
2.  **Logging**: 단순 `console.log` 삭제보다는, **디버깅 레벨(`console.debug`)을 적용하거나 `console.group`으로 구조화**하여 유의미한 로그 자산(Log Enrichment)으로 남김. (단, 잡음/노이즈성 로그는 삭제)
3.  **Clean-up**: 죽은 코드(Dead Code), 불필요한 주석 제거.
4.  **Typing**: `any`는 구체적 타입으로 명시하되, **타입 정의에 10분 이상 소요될 경우 생략**한다.
5.  **Structure**: 중첩된 `if/else`는 **Guard Clause (Early Return)** 패턴으로 평탄화.
6.  **Constants**: 문자열 리터럴, 매직 넘버 등 의미를 알 수 없는 값은 **상수(Enum/Constants)**로 추출하여 관리한다.

## ⛔️ 주의사항 (Constraints)
1.  **동작 변경 절대 금지 (Behavior Preservation)**: 
    *   Polishing은 로직(Business Logic)을 바꾸는 것이 아니라 코드의 품질과 마감을 개선하는 것입니다.
    *   **"이 기능 있으면 좋겠는데?"라는 생각이 들어도, 본래 과제가 아니라면 절대 임의로 추가하지 마세요.** (Feature Creep 엄금)
2.  **과몰입 금지**: 배보다 배꼽이 커지지 않도록, 전체 작업 시간의 **20%**를 넘기지 마세요.
3.  **보고**: PR이나 작업 요약 시, **[Polish]** 항목으로 지침 적용 내역을 분리하여 보고해주세요.
