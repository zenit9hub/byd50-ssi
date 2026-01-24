package com.byd50.ssi.demo

import android.util.Log
import org.json.JSONArray
import org.json.JSONObject
import java.io.BufferedReader
import java.io.InputStreamReader
import java.net.HttpURLConnection
import java.net.URL

class ApiClient(private val baseUrl: String) {
    private val tag = TAG

    data class Challenge(val aud: String, val nonce: String)
    data class LicenseIssueResult(
        val simplePresentationValid: Boolean,
        val vcJwt: String,
        val error: String
    )
    data class RentalIssueResult(
        val vpSignatureValid: Boolean,
        val audNonceValid: Boolean,
        val vcValid: Boolean,
        val vcNotExpired: Boolean,
        val holderDidMatch: Boolean,
        val vcJwt: String,
        val error: String
    )

    fun createDid(method: String, publicKeyBase58: String): String {
        val body = JSONObject()
        body.put("method", method)
        body.put("public_key_base58", publicKeyBase58)
        val resp = post(PATH_CREATE_DID, body)
        return resp.optString("did", "")
    }

    fun createVc(
        kid: String,
        pvKeyBase58: String,
        credType: String,
        credentialSubject: JSONObject,
        expiresInMinutes: Int
    ): String {
        val body = JSONObject()
        body.put("kid", kid)
        body.put("pv_key_base58", pvKeyBase58)
        body.put("type", credType)
        body.put("issuer", DEFAULT_ISSUER)
        body.put("subject", kid)
        if (expiresInMinutes > 0) {
            body.put("expires_in_minutes", expiresInMinutes)
        }
        body.put("credential_subject", credentialSubject)
        val resp = post(PATH_VC_CREATE, body)
        return resp.optString("vc_jwt", "")
    }

    fun verifyVc(vcJwt: String): Boolean {
        val body = JSONObject()
        body.put("vc_jwt", vcJwt)
        val resp = post(PATH_VC_VERIFY, body)
        return resp.optBoolean("valid", false)
    }

    fun createVp(
        holderDid: String,
        pvKeyBase58: String,
        vcJwts: List<String>,
        aud: String,
        nonce: String,
        simplePresentation: Boolean
    ): String {
        val body = JSONObject()
        body.put("holder_did", holderDid)
        body.put("pv_key_base58", pvKeyBase58)
        body.put("type", VP_TYPE_DEFAULT)
        body.put("issuer", holderDid)
        body.put("subject", holderDid)
        body.put("expires_in_minutes", DEFAULT_VP_EXPIRES_MIN)
        val vcArray = JSONArray()
        vcJwts.forEach { vcArray.put(it) }
        body.put("vc_jwts", vcArray)
        body.put("aud", aud)
        body.put("nonce", nonce)
        body.put("simple_presentation", simplePresentation)
        val resp = post(PATH_VP_CREATE, body)
        return resp.optString("vp_jwt", "")
    }

    fun verifyVp(vpJwt: String, expectedAud: String, expectedNonce: String): Boolean {
        val body = JSONObject()
        body.put("vp_jwt", vpJwt)
        if (expectedAud.isNotBlank()) {
            body.put("expected_aud", expectedAud)
        }
        if (expectedNonce.isNotBlank()) {
            body.put("expected_nonce", expectedNonce)
        }
        val resp = post(PATH_VP_VERIFY, body)
        return resp.optBoolean("valid", false)
    }

    fun requestLicenseChallenge(): Challenge {
        val resp = post(PATH_LICENSE_CHALLENGE, JSONObject())
        return Challenge(resp.optString("aud", ""), resp.optString("nonce", ""))
    }

    fun requestRentalChallenge(): Challenge {
        val resp = post(PATH_RENTAL_CHALLENGE, JSONObject())
        return Challenge(resp.optString("aud", ""), resp.optString("nonce", ""))
    }

    fun issueLicense(
        holderDid: String,
        simpleVpJwt: String,
        expectedAud: String,
        expectedNonce: String,
        expiresInSeconds: Int
    ): LicenseIssueResult {
        val body = JSONObject()
        body.put("holder_did", holderDid)
        body.put("simple_vp_jwt", simpleVpJwt)
        body.put("expected_aud", expectedAud)
        body.put("expected_nonce", expectedNonce)
        if (expiresInSeconds > 0) {
            body.put("expires_in_seconds", expiresInSeconds)
        }
        val resp = post(PATH_LICENSE_ISSUE, body)
        return LicenseIssueResult(
            resp.optBoolean("simple_presentation_valid", false),
            resp.optString("vc_jwt", ""),
            resp.optString("error", "")
        )
    }

    fun issueRental(
        vpJwt: String,
        expectedAud: String,
        expectedNonce: String,
        expiresInSeconds: Int
    ): RentalIssueResult {
        val body = JSONObject()
        body.put("vp_jwt", vpJwt)
        body.put("expected_aud", expectedAud)
        body.put("expected_nonce", expectedNonce)
        if (expiresInSeconds > 0) {
            body.put("expires_in_seconds", expiresInSeconds)
        }
        val resp = post(PATH_RENTAL_ISSUE, body)
        return RentalIssueResult(
            resp.optBoolean("vp_signature_valid", false),
            resp.optBoolean("aud_nonce_valid", false),
            resp.optBoolean("vc_valid", false),
            resp.optBoolean("vc_not_expired", false),
            resp.optBoolean("holder_did_match", false),
            resp.optString("vc_jwt", ""),
            resp.optString("error", "")
        )
    }

    private fun post(path: String, body: JSONObject): JSONObject {
        val url = URL(baseUrl.trimEnd('/') + path)
        val conn = url.openConnection() as HttpURLConnection
        conn.requestMethod = "POST"
        conn.setRequestProperty("Content-Type", "application/json")
        conn.doOutput = true
        Log.d(tag, "POST ${url} body=${body}")

        conn.outputStream.use { os ->
            os.write(body.toString().toByteArray())
        }

        val code = conn.responseCode
        Log.d(tag, "HTTP ${code}")
        val reader = if (code in 200..299) {
            BufferedReader(InputStreamReader(conn.inputStream))
        } else {
            BufferedReader(InputStreamReader(conn.errorStream))
        }

        val content = reader.readText()
        reader.close()
        if (content.isBlank()) {
            Log.w(tag, "Empty response body")
            return JSONObject()
        }
        Log.d(tag, "Response body: ${content.take(RESPONSE_LOG_MAX)}")
        return JSONObject(content)
    }

    companion object {
        private const val TAG = "ApiClient"
        private const val DEFAULT_ISSUER = "http://demo-issuer.example"
        private const val VP_TYPE_DEFAULT = "CredentialManagerPresentation"
        private const val DEFAULT_VP_EXPIRES_MIN = 5
        private const val RESPONSE_LOG_MAX = 512
        private const val PATH_CREATE_DID = "/v2/testapi/create-did"
        private const val PATH_VC_CREATE = "/v2/testapi/vc/create"
        private const val PATH_VC_VERIFY = "/v2/testapi/vc/verify"
        private const val PATH_VP_CREATE = "/v2/testapi/vp/create"
        private const val PATH_VP_VERIFY = "/v2/testapi/vp/verify"
        private const val PATH_LICENSE_CHALLENGE = "/v2/testapi/license/challenge"
        private const val PATH_LICENSE_ISSUE = "/v2/testapi/license/issue"
        private const val PATH_RENTAL_CHALLENGE = "/v2/testapi/rental/challenge"
        private const val PATH_RENTAL_ISSUE = "/v2/testapi/rental/issue"
    }
}
