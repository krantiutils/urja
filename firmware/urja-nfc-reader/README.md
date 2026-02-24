# Urja NFC Door Reader

ESP32-WROOM + RC522 firmware for gym door access control.

## Wiring

```
ESP32-WROOM          RC522
────────────         ─────
GPIO 5  (SS)    →    SDA
GPIO 22 (RST)   →    RST
GPIO 23 (MOSI)  →    MOSI
GPIO 19 (MISO)  →    MISO
GPIO 18 (SCK)   →    SCK
3.3V            →    3.3V
GND             →    GND

ESP32-WROOM          Relay Module (12V)
────────────         ───────────────────
GPIO 4          →    IN
3.3V            →    VCC
GND             →    GND
                     COM  →  12V Power Supply +
                     NO   →  EM Lock +
                     EM Lock - → 12V Power Supply -

ESP32-WROOM          LEDs / Buzzer
────────────         ────────────
GPIO 2          →    Green LED (+) → 220Ω → GND
GPIO 15         →    Red LED (+)   → 220Ω → GND
GPIO 13         →    Buzzer (+)    → GND
```

## Setup

1. Flash the firmware via PlatformIO
2. Open serial monitor (115200 baud)
3. On first boot, configure via serial commands:
   ```
   SSID YourWiFiName
   PASS YourWiFiPassword
   API http://your-server:8080
   SAVE
   ```
4. The device prints its **secret key** — copy it
5. In the Urja dashboard, go to NFC → Register Device:
   - Name: e.g. "Front Door"
   - Device Identifier: ESP32 MAC address
   - Device Secret: paste the key from step 4

## How It Works

1. Member taps NFC card on RC522 reader
2. ESP32 reads card UID (hex)
3. ESP32 sends `POST /api/v1/devices/check-in` with:
   - Header: `X-Device-Key: <device_secret>`
   - Body: `{"card_uid": "<hex>"}`
4. Backend verifies device → looks up card → records attendance
5. On success: relay opens (5s), green LED, two beeps
6. On failure: red LED, long beep

## Serial Commands (runtime)

- `SHOW` — print current config and WiFi status
- `REBOOT` — restart the device

## Security

- Device generates a random 64-char hex secret on first boot
- Secret is stored in ESP32 NVS flash (never transmitted in plaintext except in X-Device-Key header — use HTTPS in production)
- Backend stores only SHA-256(secret), never the raw value
- Even if the database is breached, device keys cannot be recovered
