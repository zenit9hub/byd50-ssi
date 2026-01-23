package com.byd50.ssi.demo

import android.util.Log

object NativeBridge {
    private const val TAG = "NativeBridge"
    val isLoaded: Boolean

    init {
        var loaded = false
        try {
            System.loadLibrary("foo")
            System.loadLibrary("native-bridge")
            loaded = true
            Log.d(TAG, "JNI libraries loaded")
        } catch (e: UnsatisfiedLinkError) {
            Log.e(TAG, "JNI load failed: ${e.message}")
        }
        isLoaded = loaded
    }

    external fun createKeyPairNative(): String
    external fun getPublicKeyBase58Native(): String
    external fun getPrivateKeyBase58Native(): String
    external fun createVpNative(did: String, issuer: String, pvKeyBase58: String, credType: String, vcJwt: String): String

    fun createKeyPair(): String {
        return if (isLoaded) createKeyPairNative() else "not_loaded"
    }

    fun getPublicKeyBase58(): String {
        return if (isLoaded) getPublicKeyBase58Native() else ""
    }

    fun getPrivateKeyBase58(): String {
        return if (isLoaded) getPrivateKeyBase58Native() else ""
    }
}
