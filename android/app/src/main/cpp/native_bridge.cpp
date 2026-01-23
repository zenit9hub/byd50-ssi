#include <jni.h>
#include <string>
#include <stdlib.h>
#include "libfoo.h"

extern "C" JNIEXPORT jstring JNICALL
Java_com_byd50_ssi_demo_NativeBridge_createKeyPairNative(JNIEnv *env, jobject /* this */) {
    char* result = createKeyPair();
    jstring out = env->NewStringUTF(result ? result : "");
    if (result) {
        free(result);
    }
    return out;
}

extern "C" JNIEXPORT jstring JNICALL
Java_com_byd50_ssi_demo_NativeBridge_getPublicKeyBase58Native(JNIEnv *env, jobject /* this */) {
    char* result = getPbKey();
    jstring out = env->NewStringUTF(result ? result : "");
    if (result) {
        free(result);
    }
    return out;
}

extern "C" JNIEXPORT jstring JNICALL
Java_com_byd50_ssi_demo_NativeBridge_getPrivateKeyBase58Native(JNIEnv *env, jobject /* this */) {
    char* result = getPvKey();
    jstring out = env->NewStringUTF(result ? result : "");
    if (result) {
        free(result);
    }
    return out;
}

extern "C" JNIEXPORT jstring JNICALL
Java_com_byd50_ssi_demo_NativeBridge_createVpNative(
        JNIEnv *env,
        jobject /* this */,
        jstring did,
        jstring issuer,
        jstring pvKeyBase58,
        jstring credType,
        jstring vcJwt) {
    const char *c_did = env->GetStringUTFChars(did, nullptr);
    const char *c_issuer = env->GetStringUTFChars(issuer, nullptr);
    const char *c_pv = env->GetStringUTFChars(pvKeyBase58, nullptr);
    const char *c_type = env->GetStringUTFChars(credType, nullptr);
    const char *c_vc = env->GetStringUTFChars(vcJwt, nullptr);

    char* result = createVp(const_cast<char*>(c_did), const_cast<char*>(c_issuer), const_cast<char*>(c_pv), const_cast<char*>(c_type), const_cast<char*>(c_vc));
    jstring out = env->NewStringUTF(result ? result : "");

    if (result) {
        free(result);
    }

    env->ReleaseStringUTFChars(did, c_did);
    env->ReleaseStringUTFChars(issuer, c_issuer);
    env->ReleaseStringUTFChars(pvKeyBase58, c_pv);
    env->ReleaseStringUTFChars(credType, c_type);
    env->ReleaseStringUTFChars(vcJwt, c_vc);

    return out;
}
