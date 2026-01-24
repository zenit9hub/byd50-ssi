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

    private val tag = "MainActivity"
    private lateinit var binding: ActivityMainBinding
    private val logBuffer = StringBuilder()
    private val timeFormat = SimpleDateFormat("HH:mm:ss", Locale.getDefault())

    private var dlTimer: CountDownTimer? = null
    private var rentalTimer: CountDownTimer? = null

    private var scenario: DemoScenario? = null
    private var scenarioBaseUrl: String = ""

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityMainBinding.inflate(layoutInflater)
        setContentView(binding.root)

        binding.inputBaseUrl.setText("http://10.0.2.2:8080")
        logLine("App ready. JNI loaded: ${NativeBridge.isLoaded}")
        setDidLabel("")

        setupExpiryInput(binding.inputDlExpirySec)
        setupExpiryInput(binding.inputRentalExpirySec)
        setStatus(binding.txtDlStatus, "상태: 대기")
        setStatus(binding.txtRentalStatus, "상태: 대기")
        resetAllChecks()

        binding.btnCreateDid.setOnClickListener {
            runAsync(
                "DID 발급",
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
            val expirySec = readExpirySeconds(binding.inputDlExpirySec, 300)
            resetLicenseChecks()
            runAsync(
                "전자면허증 VC 발급",
                block = { getScenario(it).issueDriverLicense(expirySec) },
                onSuccess = { result ->
                    if (result.startsWith("ey")) {
                        startCountdown(
                            "전자면허증 VC",
                            expirySec,
                            binding.txtDlStatus,
                            binding.btnIssueDl,
                            isLong = expirySec > 0
                        )
                    } else {
                        setStatus(binding.txtDlStatus, "상태: 발급 실패")
                    }
                    updateLicenseChecks()
                },
                onError = { msg ->
                    setStatus(binding.txtDlStatus, "상태: 오류 - $msg")
                    updateLicenseChecks()
                }
            )
        }

        binding.btnIssueRental.setOnClickListener {
            val expirySec = readExpirySeconds(binding.inputRentalExpirySec, 60)
            resetRentalChecks()
            runAsync(
                "렌터카 계약 VC 발급",
                block = { getScenario(it).issueRentalAgreement(expirySec) },
                onSuccess = { result ->
                    if (result.startsWith("면허 VC")) {
                        setStatus(binding.txtRentalStatus, "상태: ${result}")
                        updateRentalChecks()
                        return@runAsync
                    }
                    if (result.startsWith("ey")) {
                        startCountdown(
                            "렌터카 계약 VC",
                            expirySec,
                            binding.txtRentalStatus,
                            binding.btnIssueRental,
                            isLong = expirySec > 0
                        )
                    } else {
                        setStatus(binding.txtRentalStatus, "상태: 발급 실패")
                    }
                    updateRentalChecks()
                },
                onError = { msg ->
                    setStatus(binding.txtRentalStatus, "상태: 오류 - $msg")
                    updateRentalChecks()
                }
            )
        }

        binding.btnCarAccess.setOnClickListener {
            runAsync(
                "차량 접근 테스트",
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
                logLine("[JNI] library not loaded")
                return@setOnClickListener
            }
            val pb = NativeBridge.getPublicKeyBase58()
            logLine("[JNI] PublicKeyBase58: ${pb.take(48)}...")
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
            logLine("Scenario reset for baseUrl=$baseUrl")
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
                    logLine("[DONE] $title: ${result.take(72)}")
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
        val hintColor = 0xFF888888.toInt()
        val textColor = 0xFF222222.toInt()
        input.setTextColor(textColor)
        input.addTextChangedListener(SimpleTextWatcher {
            val hasText = input.text?.isNotBlank() == true
            input.setTextColor(if (hasText) textColor else hintColor)
        })
        input.setTextColor(hintColor)
    }

    private fun readExpirySeconds(input: EditText, defaultSeconds: Int): Int {
        val raw = input.text?.toString()?.trim().orEmpty()
        val parsed = raw.toIntOrNull()
        return if (parsed != null && parsed > 0) parsed else defaultSeconds
    }

    private fun startCountdown(
        label: String,
        totalSeconds: Int,
        statusView: TextView,
        button: Button,
        isLong: Boolean
    ) {
        if (totalSeconds <= 0) {
            setStatus(statusView, "상태: 즉시 만료")
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
                    "상태: 유효 (남은 %02d:%02d / 만료 %s)",
                    mm,
                    ss,
                    endTimeText
                )
                setStatus(statusView, status, remainSec <= 30)
            }

            override fun onFinish() {
                setStatus(statusView, "상태: 만료")
            }
        }
        assignTimer(statusView, timer)
        timer.start()
        if (!isLong) {
            setStatus(statusView, "상태: 유효 (만료 $endTimeText)")
        }
    }

    private fun setStatus(view: TextView, text: String, warn: Boolean = false) {
        view.text = text
        view.setTextColor(if (warn) 0xFFD32F2F.toInt() else 0xFF444444.toInt())
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
        setCheck(binding.txtCarCheckVp, "VP 무결성 체크", null)
    }

    private fun resetLicenseChecks() {
        setCheck(binding.txtDlCheckDidAuth, "DID 인증 체크", null)
        setCheck(binding.txtDlCheckVcIssued, "VC 생성 발급", null)
    }

    private fun resetRentalChecks() {
        setCheck(binding.txtRentalCheckVpSig, "VP 서명 검증", null)
        setCheck(binding.txtRentalCheckAudNonce, "VP에 포함 요구한 2개 필드 체크", null)
        setCheck(binding.txtRentalCheckVcIntegrity, "VP에 포함된 VC 무결성 체크", null)
    }

    private fun updateLicenseChecks() {
        val s = scenario ?: return
        setCheck(binding.txtDlCheckDidAuth, "DID 인증 체크", s.licenseDidAuthValid)
        setCheck(binding.txtDlCheckVcIssued, "VC 생성 발급", s.licenseVcIssued)
    }

    private fun updateRentalChecks() {
        val s = scenario ?: return
        setCheck(binding.txtRentalCheckVpSig, "VP 서명 검증", s.rentalVpSigValid)
        setCheck(binding.txtRentalCheckAudNonce, "VP에 포함 요구한 2개 필드 체크", s.rentalAudNonceValid)
        val integrityOk = when {
            s.rentalVcIntegrityValid == null && s.rentalHolderMatchValid == null -> null
            s.rentalVcIntegrityValid == true && s.rentalHolderMatchValid == true -> true
            s.rentalVcIntegrityValid == false || s.rentalHolderMatchValid == false -> false
            else -> null
        }
        setCheck(binding.txtRentalCheckVcIntegrity, "VP에 포함된 VC 무결성 체크", integrityOk)
    }

    private fun updateCarAccessCheck() {
        val s = scenario ?: return
        setCheck(binding.txtCarCheckVp, "VP 무결성 체크", s.carVpValid)
    }

    private fun setCheck(view: TextView, label: String, state: Boolean?) {
        val prefix = when (state) {
            true -> "[v]"
            false -> "[x]"
            null -> "[ ]"
        }
        view.text = "$prefix $label"
        val color = when (state) {
            true -> 0xFF2E7D32.toInt()
            false -> 0xFFD32F2F.toInt()
            null -> 0xFF777777.toInt()
        }
        view.setTextColor(color)
    }

    private fun setDidLabel(did: String) {
        val value = if (did.startsWith("did:")) did else "없음"
        binding.txtDid.text = "현재 DID: $value"
    }
}
