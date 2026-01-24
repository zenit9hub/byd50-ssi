package com.byd50.ssi.demo

import android.util.Log
import org.json.JSONArray
import org.json.JSONObject
import java.io.BufferedReader
import java.io.InputStreamReader
import java.net.HttpURLConnection
import java.net.URL

class ApiClient(private val baseUrl: String) {
    private val tag = "ApiClient"

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
        val resp = post("/v2/testapi/create-did", body)
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
        body.put("issuer", "http://demo-issuer.example")
        body.put("subject", kid)
        if (expiresInMinutes > 0) {
            body.put("expires_in_minutes", expiresInMinutes)
        }
        body.put("credential_subject", credentialSubject)
        val resp = post("/v2/testapi/vc/create", body)
        return resp.optString("vc_jwt", "")
    }

    fun verifyVc(vcJwt: String): Boolean {
        val body = JSONObject()
        body.put("vc_jwt", vcJwt)
        val resp = post("/v2/testapi/vc/verify", body)
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
        body.put("type", "CredentialManagerPresentation")
        body.put("issuer", holderDid)
        body.put("subject", holderDid)
        body.put("expires_in_minutes", 5)
        val vcArray = JSONArray()
        vcJwts.forEach { vcArray.put(it) }
        body.put("vc_jwts", vcArray)
        body.put("aud", aud)
        body.put("nonce", nonce)
        body.put("simple_presentation", simplePresentation)
        val resp = post("/v2/testapi/vp/create", body)
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
        val resp = post("/v2/testapi/vp/verify", body)
        return resp.optBoolean("valid", false)
    }

    fun requestLicenseChallenge(): Challenge {
        val resp = post("/v2/testapi/license/challenge", JSONObject())
        return Challenge(resp.optString("aud", ""), resp.optString("nonce", ""))
    }

    fun requestRentalChallenge(): Challenge {
        val resp = post("/v2/testapi/rental/challenge", JSONObject())
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
        val resp = post("/v2/testapi/license/issue", body)
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
        val resp = post("/v2/testapi/rental/issue", body)
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
        Log.d(tag, "Response body: ${content.take(512)}")
        return JSONObject(content)
    }
}
