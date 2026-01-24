package com.byd50.ssi.demo

import java.util.UUID

class DemoScenario(private val api: ApiClient, private val native: NativeBridge) {
    var did: String = ""
        private set
    var dlVc: String = ""
        private set
    var rentalVc: String = ""
        private set

    var licenseDidAuthValid: Boolean? = null
        private set
    var licenseVcIssued: Boolean? = null
        private set
    var rentalVpSigValid: Boolean? = null
        private set
    var rentalAudNonceValid: Boolean? = null
        private set
    var rentalVcIntegrityValid: Boolean? = null
        private set
    var rentalHolderMatchValid: Boolean? = null
        private set
    var carVpValid: Boolean? = null
        private set

    fun createDid(): String {
        native.createKeyPair()
        val pb = native.getPublicKeyBase58()
        did = api.createDid("byd50", pb)
        return did
    }

    fun issueDriverLicense(expiresInSeconds: Int): String {
        resetLicenseChecks()
        if (did.isEmpty()) return ERR_NEED_DID
        val challenge = api.requestLicenseChallenge()
        val vp = api.createVp(
            did,
            native.getPrivateKeyBase58(),
            emptyList(),
            challenge.aud,
            challenge.nonce,
            true
        )
        if (vp.isBlank()) {
            licenseDidAuthValid = false
            return ERR_DID_AUTH_FAILED
        }
        val issue = api.issueLicense(
            did,
            vp,
            challenge.aud,
            challenge.nonce,
            expiresInSeconds
        )
        licenseDidAuthValid = issue.simplePresentationValid
        dlVc = issue.vcJwt
        licenseVcIssued = dlVc.isNotBlank()
        if (!licenseDidAuthValid!!) {
            return ERR_DID_AUTH_FAILED
        }
        if (!licenseVcIssued!!) {
            return if (issue.error.isNotBlank()) issue.error else ERR_LICENSE_ISSUE_FAILED
        }
        return dlVc
    }

    fun issueRentalAgreement(expiresInSeconds: Int): String {
        resetRentalChecks()
        if (did.isEmpty()) return ERR_NEED_DID
        if (dlVc.isEmpty()) return ERR_NEED_LICENSE
        val challenge = api.requestRentalChallenge()
        val vp = api.createVp(
            did,
            native.getPrivateKeyBase58(),
            listOf(dlVc),
            challenge.aud,
            challenge.nonce,
            false
        )
        if (vp.isBlank()) {
            rentalVpSigValid = false
            return ERR_VP_CREATE_FAILED
        }
        val issue = api.issueRental(
            vp,
            challenge.aud,
            challenge.nonce,
            expiresInSeconds
        )
        rentalVpSigValid = issue.vpSignatureValid
        rentalAudNonceValid = issue.audNonceValid
        rentalVcIntegrityValid = issue.vcValid && issue.vcNotExpired
        rentalHolderMatchValid = issue.holderDidMatch
        rentalVc = issue.vcJwt
        if (!issue.vpSignatureValid || !issue.audNonceValid || !issue.vcValid || !issue.vcNotExpired || !issue.holderDidMatch) {
            return if (issue.error.isNotBlank()) issue.error else ERR_RENTAL_ISSUE_FAILED
        }
        return rentalVc
    }

    fun carAccessCheck(): String {
        if (rentalVc.isEmpty()) {
            carVpValid = null
            return ERR_NO_RENTAL
        }
        val nonce = UUID.randomUUID().toString().take(8)
        val vp = api.createVp(
            did,
            native.getPrivateKeyBase58(),
            listOf(rentalVc),
            "",
            nonce,
            false
        )
        val valid = api.verifyVp(vp, "", nonce)
        carVpValid = valid
        return if (valid) ACCESS_ALLOWED else ACCESS_DENIED
    }

    private fun resetLicenseChecks() {
        licenseDidAuthValid = null
        licenseVcIssued = null
    }

    private fun resetRentalChecks() {
        rentalVpSigValid = null
        rentalAudNonceValid = null
        rentalVcIntegrityValid = null
        rentalHolderMatchValid = null
    }

    companion object {
        private const val ERR_NEED_DID = "DID 없음 - 1번 먼저 실행"
        private const val ERR_NEED_LICENSE = "면허 VC 없음 - 2번 먼저 실행"
        private const val ERR_DID_AUTH_FAILED = "DID 인증 실패"
        private const val ERR_LICENSE_ISSUE_FAILED = "면허 VC 발급 실패"
        private const val ERR_VP_CREATE_FAILED = "VP 생성 실패"
        private const val ERR_RENTAL_ISSUE_FAILED = "렌터카 계약 실패"
        private const val ERR_NO_RENTAL = "계약 없음"
        private const val ACCESS_ALLOWED = "접근 허용(데모)"
        private const val ACCESS_DENIED = "접근 거부(데모)"
    }
}
