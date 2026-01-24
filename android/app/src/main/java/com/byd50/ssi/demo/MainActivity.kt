package com.byd50.ssi.demo

import android.os.Bundle
import android.os.CountDownTimer
import android.os.Looper
import android.util.Log
import android.widget.Button
import android.widget.EditText
import android.widget.TextView
import androidx.appcompat.app.AppCompatActivity
import com.byd50.ssi.demo.databinding.ActivityMainBinding
import java.text.SimpleDateFormat
import java.util.Date
import java.util.Locale

class MainActivity : AppCompatActivity() {

    private val tag = TAG
    private lateinit var binding: ActivityMainBinding
    private val logBuffer = StringBuilder()
    private val timeFormat = SimpleDateFormat(TIME_FORMAT_PATTERN, Locale.getDefault())

    private var dlTimer: CountDownTimer? = null
    private var rentalTimer: CountDownTimer? = null

    private var scenario: DemoScenario? = null
    private var scenarioBaseUrl: String = ""

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityMainBinding.inflate(layoutInflater)
        setContentView(binding.root)

        binding.inputBaseUrl.setText(DEFAULT_BASE_URL)
        logLine("${LOG_APP_READY} ${NativeBridge.isLoaded}")
        setDidLabel("")

        setupExpiryInput(binding.inputDlExpirySec)
        setupExpiryInput(binding.inputRentalExpirySec)
        setStatus(binding.txtDlStatus, STATUS_READY)
        setStatus(binding.txtRentalStatus, STATUS_READY)
        resetAllChecks()

        binding.btnCreateDid.setOnClickListener {
            runAsync(
                ACTION_CREATE_DID,
                block = { getScenario(it).createDid() },
                onSuccess = { result ->
                    setDidLabel(result)
                },
                onError = {
                    setDidLabel("")
                }
            )
        }

        binding.btnIssueDl.setOnClickListener {
            val expirySec = readExpirySeconds(binding.inputDlExpirySec, DEFAULT_LICENSE_EXPIRY_SEC)
            resetLicenseChecks()
            runAsync(
                ACTION_ISSUE_LICENSE,
                block = { getScenario(it).issueDriverLicense(expirySec) },
                onSuccess = { result ->
                    if (isJwt(result)) {
                        startCountdown(
                            expirySec,
                            binding.txtDlStatus,
                            binding.btnIssueDl,
                            isLong = expirySec > 0
                        )
                    } else {
                        setStatus(binding.txtDlStatus, STATUS_ISSUE_FAILED)
                    }
                    updateLicenseChecks()
                },
                onError = { msg ->
                    setStatus(binding.txtDlStatus, "$STATUS_ERROR_PREFIX $msg")
                    updateLicenseChecks()
                }
            )
        }

        binding.btnIssueRental.setOnClickListener {
            val expirySec = readExpirySeconds(binding.inputRentalExpirySec, DEFAULT_RENTAL_EXPIRY_SEC)
            resetRentalChecks()
            runAsync(
                ACTION_ISSUE_RENTAL,
                block = { getScenario(it).issueRentalAgreement(expirySec) },
                onSuccess = { result ->
                    if (result.startsWith(LICENSE_PREFIX)) {
                        setStatus(binding.txtRentalStatus, "$STATUS_PREFIX $result")
                        updateRentalChecks()
                        return@runAsync
                    }
                    if (isJwt(result)) {
                        startCountdown(
                            expirySec,
                            binding.txtRentalStatus,
                            binding.btnIssueRental,
                            isLong = expirySec > 0
                        )
                    } else {
                        setStatus(binding.txtRentalStatus, STATUS_ISSUE_FAILED)
                    }
                    updateRentalChecks()
                },
                onError = { msg ->
                    setStatus(binding.txtRentalStatus, "$STATUS_ERROR_PREFIX $msg")
                    updateRentalChecks()
                }
            )
        }

        binding.btnCarAccess.setOnClickListener {
            runAsync(
                ACTION_CAR_ACCESS,
                block = { getScenario(it).carAccessCheck() },
                onSuccess = {
                    updateCarAccessCheck()
                },
                onError = {
                    updateCarAccessCheck()
                }
            )
        }

        binding.btnShowKeys.setOnClickListener {
            if (!NativeBridge.isLoaded) {
                logLine(LOG_JNI_NOT_LOADED)
                return@setOnClickListener
            }
            val pb = NativeBridge.getPublicKeyBase58()
            logLine("$LOG_JNI_PUBKEY ${pb.take(48)}...")
        }

        binding.btnClear.setOnClickListener {
            logBuffer.clear()
            binding.txtLog.text = ""
        }
    }

    private fun getScenario(baseUrl: String): DemoScenario {
        if (scenario == null || baseUrl != scenarioBaseUrl) {
            scenario = DemoScenario(ApiClient(baseUrl), NativeBridge)
            scenarioBaseUrl = baseUrl
            logLine("$LOG_SCENARIO_RESET $baseUrl")
            resetAllChecks()
            setDidLabel("")
        }
        return scenario!!
    }

    private fun runAsync(
        title: String,
        block: (baseUrl: String) -> String,
        onSuccess: (String) -> Unit = {},
        onError: (String) -> Unit = {}
    ) {
        logLine("[START] $title")
        val baseUrl = binding.inputBaseUrl.text.toString().trim()
        Thread {
            try {
                val result = block(baseUrl)
                runOnUiThread {
                    logLine("[DONE] $title: ${result.take(LOG_RESULT_MAX)}")
                    onSuccess(result)
                }
            } catch (e: Exception) {
                runOnUiThread {
                    logLine("[ERROR] $title: ${e.message}")
                    onError(e.message ?: "unknown error")
                }
                Log.e(tag, "Scenario error: $title", e)
            }
        }.start()
    }

    private fun logLine(message: String) {
        if (Looper.myLooper() == Looper.getMainLooper()) {
            logLineInternal(message)
        } else {
            runOnUiThread { logLineInternal(message) }
        }
    }

    private fun logLineInternal(message: String) {
        logBuffer.append(message).append('\n')
        binding.txtLog.text = logBuffer.toString()
        Log.d(tag, message)
    }

    private fun setupExpiryInput(input: EditText) {
        input.setTextColor(TEXT_COLOR)
        input.addTextChangedListener(SimpleTextWatcher {
            val hasText = input.text?.isNotBlank() == true
            input.setTextColor(if (hasText) TEXT_COLOR else HINT_COLOR)
        })
        input.setTextColor(HINT_COLOR)
    }

    private fun readExpirySeconds(input: EditText, defaultSeconds: Int): Int {
        val raw = input.text?.toString()?.trim().orEmpty()
        val parsed = raw.toIntOrNull()
        return if (parsed != null && parsed > 0) parsed else defaultSeconds
    }

    private fun startCountdown(
        totalSeconds: Int,
        statusView: TextView,
        button: Button,
        isLong: Boolean
    ) {
        if (totalSeconds <= 0) {
            setStatus(statusView, STATUS_EXPIRED_NOW)
            return
        }
        val endAt = System.currentTimeMillis() + totalSeconds * 1000L
        val endTimeText = timeFormat.format(Date(endAt))
        cancelTimerFor(statusView)
        val timer = object : CountDownTimer(totalSeconds * 1000L, 1000L) {
            override fun onTick(millisUntilFinished: Long) {
                val remainSec = (millisUntilFinished / 1000L).toInt()
                val mm = remainSec / 60
                val ss = remainSec % 60
                val status = String.format(
                    Locale.getDefault(),
                    STATUS_VALID_TEMPLATE,
                    mm,
                    ss,
                    endTimeText
                )
                setStatus(statusView, status, remainSec <= EXPIRY_WARN_THRESHOLD_SEC)
            }

            override fun onFinish() {
                setStatus(statusView, STATUS_EXPIRED, true)
            }
        }
        assignTimer(statusView, timer)
        timer.start()
        if (!isLong) {
            setStatus(statusView, String.format(STATUS_VALID_UNTIL_TEMPLATE, endTimeText))
        }
    }

    private fun setStatus(view: TextView, text: String, warn: Boolean = false) {
        view.text = text
        view.setTextColor(if (warn) COLOR_WARN else COLOR_TEXT_DARK)
    }

    private fun assignTimer(view: TextView, timer: CountDownTimer) {
        when (view.id) {
            binding.txtDlStatus.id -> {
                dlTimer?.cancel()
                dlTimer = timer
            }
            binding.txtRentalStatus.id -> {
                rentalTimer?.cancel()
                rentalTimer = timer
            }
        }
    }

    private fun cancelTimerFor(view: TextView) {
        when (view.id) {
            binding.txtDlStatus.id -> dlTimer?.cancel()
            binding.txtRentalStatus.id -> rentalTimer?.cancel()
        }
    }

    private fun resetAllChecks() {
        resetLicenseChecks()
        resetRentalChecks()
        setCheck(binding.txtCarCheckVp, CHECK_CAR_VP, null)
    }

    private fun resetLicenseChecks() {
        setCheck(binding.txtDlCheckDidAuth, CHECK_DID_AUTH, null)
        setCheck(binding.txtDlCheckVcIssued, CHECK_VC_ISSUED, null)
    }

    private fun resetRentalChecks() {
        setCheck(binding.txtRentalCheckVpSig, CHECK_VP_SIG, null)
        setCheck(binding.txtRentalCheckAudNonce, CHECK_VP_FIELDS, null)
        setCheck(binding.txtRentalCheckVcIntegrity, CHECK_VC_INTEGRITY, null)
    }

    private fun updateLicenseChecks() {
        val s = scenario ?: return
        setCheck(binding.txtDlCheckDidAuth, CHECK_DID_AUTH, s.licenseDidAuthValid)
        setCheck(binding.txtDlCheckVcIssued, CHECK_VC_ISSUED, s.licenseVcIssued)
    }

    private fun updateRentalChecks() {
        val s = scenario ?: return
        setCheck(binding.txtRentalCheckVpSig, CHECK_VP_SIG, s.rentalVpSigValid)
        setCheck(binding.txtRentalCheckAudNonce, CHECK_VP_FIELDS, s.rentalAudNonceValid)
        val integrityOk = when {
            s.rentalVcIntegrityValid == null && s.rentalHolderMatchValid == null -> null
            s.rentalVcIntegrityValid == true && s.rentalHolderMatchValid == true -> true
            s.rentalVcIntegrityValid == false || s.rentalHolderMatchValid == false -> false
            else -> null
        }
        setCheck(binding.txtRentalCheckVcIntegrity, CHECK_VC_INTEGRITY, integrityOk)
    }

    private fun updateCarAccessCheck() {
        val s = scenario ?: return
        setCheck(binding.txtCarCheckVp, CHECK_CAR_VP, s.carVpValid)
    }

    private fun setCheck(view: TextView, label: String, state: Boolean?) {
        val prefix = when (state) {
            true -> "[v]"
            false -> "[x]"
            null -> "[ ]"
        }
        view.text = "$prefix $label"
        val color = when (state) {
            true -> COLOR_OK
            false -> COLOR_WARN
            null -> COLOR_TEXT_MUTED
        }
        view.setTextColor(color)
    }

    private fun setDidLabel(did: String) {
        val value = if (did.startsWith(DID_PREFIX)) did else DID_EMPTY
        binding.txtDid.text = "$DID_LABEL $value"
    }

    private fun isJwt(value: String) = value.startsWith(JWT_PREFIX)

    companion object {
        private const val TAG = "MainActivity"
        private const val DEFAULT_BASE_URL = "http://10.0.2.2:8080"
        private const val TIME_FORMAT_PATTERN = "HH:mm:ss"
        private const val DID_PREFIX = "did:"
        private const val DID_LABEL = "현재 DID:"
        private const val DID_EMPTY = "없음"
        private const val JWT_PREFIX = "ey"
        private const val LICENSE_PREFIX = "면허 VC"
        private const val LOG_APP_READY = "App ready. JNI loaded:"
        private const val LOG_SCENARIO_RESET = "Scenario reset for baseUrl="
        private const val LOG_JNI_NOT_LOADED = "[JNI] library not loaded"
        private const val LOG_JNI_PUBKEY = "[JNI] PublicKeyBase58:"
        private const val LOG_RESULT_MAX = 72
        private const val DEFAULT_LICENSE_EXPIRY_SEC = 300
        private const val DEFAULT_RENTAL_EXPIRY_SEC = 60
        private const val EXPIRY_WARN_THRESHOLD_SEC = 30
        private const val STATUS_PREFIX = "상태:"
        private const val STATUS_READY = "상태: 대기"
        private const val STATUS_EXPIRED_NOW = "상태: 즉시 만료"
        private const val STATUS_EXPIRED = "상태: 만료"
        private const val STATUS_ISSUE_FAILED = "상태: 발급 실패"
        private const val STATUS_ERROR_PREFIX = "상태: 오류 -"
        private const val STATUS_VALID_TEMPLATE = "상태: 유효 (남은 %02d:%02d / 만료 %s)"
        private const val STATUS_VALID_UNTIL_TEMPLATE = "상태: 유효 (만료 %s)"
        private const val ACTION_CREATE_DID = "DID 발급"
        private const val ACTION_ISSUE_LICENSE = "전자면허증 VC 발급"
        private const val ACTION_ISSUE_RENTAL = "렌터카 계약 VC 발급"
        private const val ACTION_CAR_ACCESS = "차량 접근 테스트"
        private const val CHECK_DID_AUTH = "DID 인증 체크"
        private const val CHECK_VC_ISSUED = "VC 생성 발급"
        private const val CHECK_VP_SIG = "VP 서명 검증"
        private const val CHECK_VP_FIELDS = "VP에 포함 요구한 2개 필드 체크"
        private const val CHECK_VC_INTEGRITY = "VP에 포함된 VC 무결성 체크"
        private const val CHECK_CAR_VP = "VP 무결성 체크"
        private const val COLOR_WARN = 0xFFD32F2F.toInt()
        private const val COLOR_OK = 0xFF2E7D32.toInt()
        private const val COLOR_TEXT_DARK = 0xFF444444.toInt()
        private const val COLOR_TEXT_MUTED = 0xFF777777.toInt()
        private const val HINT_COLOR = 0xFF888888.toInt()
        private const val TEXT_COLOR = 0xFF222222.toInt()
    }
}
