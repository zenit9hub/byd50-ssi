package com.byd50.ssi.demo

import android.os.Bundle
import android.os.Looper
import android.util.Log
import androidx.appcompat.app.AppCompatActivity
import com.byd50.ssi.demo.databinding.ActivityMainBinding

class MainActivity : AppCompatActivity() {

    private val tag = "MainActivity"
    private lateinit var binding: ActivityMainBinding
    private val logBuffer = StringBuilder()

    private var scenario: DemoScenario? = null
    private var scenarioBaseUrl: String = ""

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityMainBinding.inflate(layoutInflater)
        setContentView(binding.root)

        binding.inputBaseUrl.setText("http://10.0.2.2:8080")
        logLine("App ready. JNI loaded: ${NativeBridge.isLoaded}")

        binding.btnCreateDid.setOnClickListener {
            runAsync("DID 발급") {
                getScenario(it).createDid()
            }
        }

        binding.btnIssueId.setOnClickListener {
            runAsync("전자주민증 VC 발급") {
                getScenario(it).issueIdCard()
            }
        }

        binding.btnIssueDl.setOnClickListener {
            runAsync("전자면허증 VC 발급") {
                getScenario(it).issueDriverLicense()
            }
        }

        binding.btnIssueRental.setOnClickListener {
            runAsync("렌터카 계약 VC 발급") {
                getScenario(it).issueRentalAgreement()
            }
        }

        binding.btnCarAccess.setOnClickListener {
            runAsync("차량 접근 테스트") {
                getScenario(it).carAccessCheck()
            }
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
        }
        return scenario!!
    }

    private fun runAsync(title: String, block: (baseUrl: String) -> String) {
        logLine("[START] $title")
        val baseUrl = binding.inputBaseUrl.text.toString().trim()
        Thread {
            try {
                val result = block(baseUrl)
                runOnUiThread {
                    logLine("[DONE] $title: ${result.take(72)}")
                }
            } catch (e: Exception) {
                runOnUiThread {
                    logLine("[ERROR] $title: ${e.message}")
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
}
