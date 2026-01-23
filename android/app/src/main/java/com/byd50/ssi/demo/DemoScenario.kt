package com.byd50.ssi.demo

import org.json.JSONObject

class DemoScenario(private val api: ApiClient, private val native: NativeBridge) {
    var did: String = ""
        private set
    var idVc: String = ""
        private set
    var dlVc: String = ""
        private set
    var rentalVc: String = ""
        private set

    fun createDid(): String {
        native.createKeyPair()
        val pb = native.getPublicKeyBase58()
        did = api.createDid("byd50", pb)
        return did
    }

    fun issueIdCard(): String {
        if (did.isEmpty()) createDid()
        val subject = JSONObject()
        subject.put("name", "Hong Gil-Dong")
        subject.put("birth", "2000-11-08")
        idVc = api.createVc(did, native.getPrivateKeyBase58(), "eIdCardCredential", subject)
        return idVc
    }

    fun issueDriverLicense(): String {
        if (idVc.isEmpty()) issueIdCard()
        val subject = JSONObject()
        subject.put("licenseType", "Type-1")
        subject.put("country", "KR")
        dlVc = api.createVc(did, native.getPrivateKeyBase58(), "eDlCardCredential", subject)
        return dlVc
    }

    fun issueRentalAgreement(): String {
        if (dlVc.isEmpty()) issueDriverLicense()
        val subject = JSONObject()
        subject.put("agreementId", "rent-${System.currentTimeMillis()}")
        subject.put("validDays", 1)
        rentalVc = api.createVc(did, native.getPrivateKeyBase58(), "RentalCarAgreementCredential", subject)
        return rentalVc
    }

    fun carAccessCheck(): String {
        if (rentalVc.isEmpty()) return "계약 없음"
        val vp = api.createVp(did, native.getPrivateKeyBase58(), rentalVc)
        val valid = api.verifyVp(vp)
        return if (valid) "접근 허용(데모)" else "접근 거부(데모)"
    }
}
