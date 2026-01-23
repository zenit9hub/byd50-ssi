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

    fun createDid(method: String, publicKeyBase58: String): String {
        val body = JSONObject()
        body.put("method", method)
        body.put("public_key_base58", publicKeyBase58)
        val resp = post("/v2/testapi/create-did", body)
        return resp.optString("did", "")
    }

    fun createVc(kid: String, pvKeyBase58: String, credType: String, credentialSubject: JSONObject): String {
        val body = JSONObject()
        body.put("kid", kid)
        body.put("pv_key_base58", pvKeyBase58)
        body.put("type", credType)
        body.put("issuer", "http://demo-issuer.example")
        body.put("subject", kid)
        body.put("expires_in_minutes", 5)
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

    fun createVp(holderDid: String, pvKeyBase58: String, vcJwt: String): String {
        val body = JSONObject()
        body.put("holder_did", holderDid)
        body.put("pv_key_base58", pvKeyBase58)
        body.put("type", "CredentialManagerPresentation")
        body.put("issuer", "client make this vp")
        body.put("subject", holderDid)
        body.put("expires_in_minutes", 5)
        val vcArray = JSONArray()
        vcArray.put(vcJwt)
        body.put("vc_jwts", vcArray)
        val resp = post("/v2/testapi/vp/create", body)
        return resp.optString("vp_jwt", "")
    }

    fun verifyVp(vpJwt: String): Boolean {
        val body = JSONObject()
        body.put("vp_jwt", vpJwt)
        val resp = post("/v2/testapi/vp/verify", body)
        return resp.optBoolean("valid", false)
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
