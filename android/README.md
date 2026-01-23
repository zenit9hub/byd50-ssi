# Android Demo (Kotlin)

## Open in Android Studio
- Open the `android/` folder as a project.
- Android Studio will sync Gradle and download dependencies.

## Demo Flow (Buttons)
1. DID 발급
2. 전자주민증 VC 발급
3. 전자면허증 VC 발급
4. 렌터카 계약 VC 발급
5. 차량 접근 테스트

## REST Gateway
- Base URL: default `http://10.0.2.2:8080` (Android emulator -> host)
- The app calls `did_service_endpoint` REST APIs for DID/VC/VP flows.

## JNI / c-shared
- Native library name: `libfoo.so`
- Build and copy to `android/app/src/main/jniLibs/<abi>/libfoo.so`
- JNI bridge library: `native-bridge` (built by CMake)

### Build JNI Library
From repo root:
```
make android
```

This produces ABI-specific `libfoo.so` into:
`android/app/src/main/jniLibs/`

## Notes
- Scenario uses REST + JNI keypair generation.
- If JNI lib is missing, buttons still work but keypair calls return fallback values.
