const go = new Go();
let wasmReady = false;

// Actors State
const actors = {
    issuer: { did: null, keys: null },
    verifier: { did: null, keys: null } // Rent-a-Car
};

// User State
const state = {
    keys: null, // {priv, pub}
    did: null,  // {did, document}
    vc: null,   // jwt string
    challenge: {
        aud: "",
        nonce: ""
    }
};

// Logger
function log(msg, type = 'info') {
    const logEl = document.getElementById('log-output');
    const entry = document.createElement('div');
    entry.className = `log-entry ${type}`;
    entry.innerText = `[${new Date().toLocaleTimeString()}] ${msg}`;
    logEl.prepend(entry);
}

// UI helpers
function el(id) { return document.getElementById(id); }

// WASM Loader
WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
    go.run(result.instance);
    wasmReady = true;
    document.getElementById('wasm-status').innerText = "WASM Ready";
    document.getElementById('wasm-status').style.color = "#58a6ff";

    // Initialize Actors (Hidden)
    setTimeout(initActors, 500);
});

function initActors() {
    log("Initializing Network Actors...", "system");

    // 1. Issuer
    const iKeys = generateKey();
    const iDid = createDID(iKeys.pub);
    actors.issuer = { keys: iKeys, did: iDid.did };
    log(`Issuer (License Authority) Online: ${iDid.did}`, "system");

    // 2. Verifier
    const vKeys = generateKey();
    const vDid = createDID(vKeys.pub);
    actors.verifier = { keys: vKeys, did: vDid.did };
    log(`Verifier (Rent-a-Car) Online: ${vDid.did}`, "system");

    log("System Ready. Please Generate User Key.", "system");
}


// ==========================================
// 1. Generate Key
// ==========================================
el('btn-gen-key').addEventListener('click', () => {
    if (!wasmReady) return;
    const res = generateKey();
    if (res.success) {
        state.keys = { priv: res.priv, pub: res.pub };
        el('out-priv-key').value = res.priv;
        el('out-pub-key').value = res.pub;
        el('btn-create-did').disabled = false;
        log("Key Pair Generated.", "success");

        // Scenario Helper: If we already had a DID, warn about rotation
        if (state.did && state.did.did) {
            log("WARNING: Key rotated but DID remains same. Auth will fail!", "error");
        }
    } else {
        log("Key Generation Failed", "error");
    }
});

// ==========================================
// 2. Create DID
// ==========================================
el('btn-create-did').addEventListener('click', () => {
    if (!state.keys) return;
    const res = createDID(state.keys.pub);
    if (res.did) {
        state.did = res;
        el('out-did').value = res.did;
        el('btn-issue-vc').disabled = false;
        log(`DID Created: ${res.did}`, "success");
    } else {
        log("DID Creation Failed", "error");
    }
});

// ==========================================
// 3. Request DL VC (with DID Auth)
// ==========================================
el('btn-issue-vc').addEventListener('click', async () => {
    if (!state.did) return;

    el('btn-issue-vc').disabled = true;
    el('auth-status-box').style.display = 'flex';

    // Step A: Challenge
    log("Requesting DL VC from Issuer...", "info");
    const challenge = Math.random().toString(36).substring(7);
    log(`[DID Auth] Issuer sent challenge: ${challenge}`, "system");

    // Simulate Network Delay
    await new Promise(r => setTimeout(r, 800));

    // Step B: User Sign
    const sigRes = signData(state.keys.priv, challenge);
    if (sigRes.error) {
        log(`[DID Auth] Signing failed: ${sigRes.error}`, "error");
        updateAuthBadge("fail");
        return;
    }
    log(`[DID Auth] User signed challenge. Sending response...`, "info");

    // Step C: Issuer Verify
    await new Promise(r => setTimeout(r, 600));

    // Note: In real world, Issuer resolves User DID. Here we simulate.
    // If User Rotated Key (Scenario 2), the Key in Registry (from DID creation) won't match the signing key (newly generated).
    // verifyData checks against Registry.

    const vRes = verifyData(state.did.did, challenge, sigRes.signature);

    if (vRes.valid) {
        log("[DID Auth] Authentication SUCCESS!", "success");
        updateAuthBadge("success");
    } else {
        log(`[DID Auth] Authentication FAILED: ${vRes.error}`, "error");
        updateAuthBadge("fail");
        el('btn-issue-vc').disabled = false;
        return;
    }

    // Step D: Issue VC
    const issuerDid = actors.issuer.did;
    const issuerPriv = actors.issuer.keys.priv;
    const subjectDid = state.did.did;

    const claims = JSON.stringify({
        licenseClass: "Type 2 Regular",
        licenseNumber: "12-34-567890",
        validUntil: "2030-12-31"
    });

    const res = issueVC(issuerDid, issuerPriv, subjectDid, claims);
    if (res.vc) {
        state.vc = res.vc;
        el('vc-display').innerHTML = `<div class="vc-content">${res.vc}</div>`;
        el('btn-copy-vc').style.display = 'inline-block';

        // Disable "Request DL" button permanently? or allow re-issue?
        el('btn-issue-vc').disabled = false;

        // Show Issuer Key
        el('issuer-key-box').style.display = 'block';
        el('out-issuer-pub-key').value = actors.issuer.keys.pub;

        log("Driver License VC Issued!", "success");

        // Highlight Next Step in Right Pane
        el('btn-req-contract').classList.add('pulse-btn');
    } else {
        log("VC Issuance Failed: " + res.error, "error");
    }
});

function updateAuthBadge(status) {
    const b = el('auth-status-badge');
    b.className = `auth-badge ${status}`;
    b.innerText = status === 'success' ? 'Verified' : 'Failed';
}

// Copy VC
el('btn-copy-vc').addEventListener('click', () => {
    if (state.vc) {
        navigator.clipboard.writeText(state.vc).then(() => {
            const originalText = el('btn-copy-vc').innerText;
            el('btn-copy-vc').innerText = "Copied!";
            setTimeout(() => el('btn-copy-vc').innerText = originalText, 2000);
            log("VC JWT Copied to clipboard", "info");
        });
    }
});

// Copy Issuer Key
el('btn-copy-issuer-key').addEventListener('click', () => {
    const key = el('out-issuer-pub-key').value;
    if (key) {
        navigator.clipboard.writeText(key).then(() => {
            const originalText = el('btn-copy-issuer-key').innerText;
            el('btn-copy-issuer-key').innerText = "Copied!";
            setTimeout(() => el('btn-copy-issuer-key').innerText = originalText, 2000);
            log("Issuer Public Key Copied", "info");
        });
    }
});

// ==========================================
// 4. Rent-a-Car Scenario
// ==========================================
el('btn-req-contract').addEventListener('click', () => {
    if (!state.vc) {
        log("Cannot request contract without Driver License VC.", "error");
        return;
    }

    // Reset Checks
    resetChecklist();
    el('contract-card').style.display = 'none';

    // Verifier Config
    const verifierDid = actors.verifier.did;
    const nonce = Math.random().toString(36).substring(7);

    state.challenge = { aud: verifierDid, nonce: nonce };

    // Remove Pulse
    el('btn-req-contract').classList.remove('pulse-btn');

    // Update UI
    el('input-aud').value = verifierDid;
    el('input-nonce').value = nonce;
    el('vp-submission-card').style.opacity = '1';
    el('vp-submission-card').style.pointerEvents = 'all';
    el('vp-submission-card').classList.add('highlight-card');
    el('btn-submit-vp').disabled = false;

    log(`Rent-a-Car requested VP. Aud: ${verifierDid}, Nonce: ${nonce}`, "system");
});

el('btn-submit-vp').addEventListener('click', async () => {
    const inputAud = el('input-aud').value;
    const inputNonce = el('input-nonce').value;
    const isTampered = el('chk-tamper').checked;

    log(`Generating VP for Aud: ${inputAud}...`);

    let holderDid = state.did.did;
    let holderPriv = state.keys.priv;

    // Create VP
    const res = createVP(holderDid, holderPriv, state.vc, inputAud, inputNonce);

    if (res.vp) {
        let vpJwt = res.vp;
        log("VP Generated. Submitting to Rent-a-Car...", "info");

        // Show Checklist Panel
        el('checklist-panel').style.display = 'block';

        if (isTampered) {
            vpJwt = vpJwt + "tampered";
            log("Simulating Network Tampering...", "error");
        }

        // Call Verifier
        // Verifier uses ITS OWN expected values (state.challenge)
        // User could have tampered inputs, but we verify against what we requested.
        const verifyNonce = state.challenge.nonce;
        const verifyAud = state.challenge.aud;

        const vRes = verifyVP(vpJwt, verifyAud, verifyNonce);

        // Visualize verification steps
        await visualizeVerification(vRes);

    } else {
        log("VP Generation Failed: " + res.error, "error");
    }
});

async function visualizeVerification(res) {
    const steps = [
        { id: 'chk-vp-sig', key: 'vpSig' },
        { id: 'chk-vc-sig', key: 'vcSig' },
        { id: 'chk-aud', key: 'audMatch' },
        { id: 'chk-nonce', key: 'nonceMatch' },
        { id: 'chk-sub', key: 'subMatch' }
    ];

    let allPass = true;
    const details = res.details || {}; // Fallback if error

    for (const step of steps) {
        await new Promise(r => setTimeout(r, 300)); // Animation delay
        const elItem = el(step.id);
        const passed = details[step.key];

        elItem.classList.remove('pending');
        if (passed) {
            elItem.classList.add('pass');
            elItem.querySelector('.check-icon').innerText = '✓';
        } else {
            elItem.classList.add('fail');
            elItem.querySelector('.check-icon').innerText = '✗';
            allPass = false;
        }
    }

    if (res.valid && allPass) {
        log("Verification Complete. issuing Contract...", "success");
        await new Promise(r => setTimeout(r, 500));

        el('contract-card').style.display = 'block';
        el('contract-details').innerHTML = `
            Holder: ${res.did}<br>
            License Verified<br>
            Contract ID: ${Math.floor(Math.random() * 100000)}
        `;
    } else {
        log(`Verification FAILED: ${res.error}`, "error");
    }
}

function resetChecklist() {
    const items = document.querySelectorAll('.check-item');
    items.forEach(item => {
        item.classList.remove('pass', 'fail');
        item.classList.add('pending');
        item.querySelector('.check-icon').innerText = '!';
    });
    el('checklist-panel').style.display = 'none';
}

// Scenarios helpers
el('btn-scenario-key-rotation').addEventListener('click', () => {
    el('btn-gen-key').click();
    log("SCENARIO: Key Rotated. Now Request Contract again (Expect Verify Fail).", "warning");
});

el('btn-scenario-sub-mismatch').addEventListener('click', () => {
    if (!state.vc) return;

    // Generate New Identity
    const kRes = generateKey();
    state.keys = { priv: kRes.priv, pub: kRes.pub };
    el('out-priv-key').value = kRes.priv;
    el('out-pub-key').value = kRes.pub;

    const dRes = createDID(kRes.pub);
    state.did = dRes;
    el('out-did').value = dRes.did;

    log("SCENARIO: Identity Reset. Old VC retained. Request Contract (Expect Subject Mismatch).", "danger");
});
