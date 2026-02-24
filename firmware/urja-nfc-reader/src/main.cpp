// ────────────────────────────────────────────────────────────────
// Urja NFC Door Reader — ESP32-WROOM + RC522
//
// Flow:
//   1. First boot: generates a random device secret, prints it on
//      serial for admin to register via the dashboard.
//   2. Admin configures WiFi + API URL via serial commands.
//   3. Normal operation: reads cards → POST /api/v1/devices/check-in
//      → unlocks relay on success.
//   4. Heartbeat every 60s → POST /api/v1/devices/heartbeat
// ────────────────────────────────────────────────────────────────

#include <Arduino.h>
#include <WiFi.h>
#include <HTTPClient.h>
#include <SPI.h>
#include <MFRC522.h>
#include <Preferences.h>
#include <ArduinoJson.h>
#include <esp_random.h>
#include "config.h"

// ── Globals ───────────────────────────────────────────────────
MFRC522 rfid(RC522_SS_PIN, RC522_RST_PIN);
Preferences prefs;

String wifiSSID;
String wifiPass;
String apiURL;      // e.g. "http://192.168.1.100:8080"
String deviceKey;   // 64-char hex secret

unsigned long lastHeartbeat = 0;
unsigned long bootTime = 0;
String lastCardUID = "";
unsigned long lastCardTime = 0;
bool doorOpen = false;
unsigned long doorOpenTime = 0;

// ── Forward declarations ─────────────────────────────────────
void setupPins();
void loadConfig();
void serialConfig();
void connectWiFi();
void generateDeviceKey();
String readCardUID();
void checkIn(const String& cardUID);
void sendHeartbeat();
void openDoor();
void closeDoor();
void beepSuccess();
void beepError();
void blinkGreen(int ms);
void blinkRed(int ms);

// ── Setup ────────────────────────────────────────────────────
void setup() {
    Serial.begin(115200);
    delay(500);

    Serial.println();
    Serial.println("╔══════════════════════════════════════╗");
    Serial.println("║     URJA NFC Door Reader v1.0       ║");
    Serial.println("╚══════════════════════════════════════╝");
    Serial.println();

    setupPins();
    SPI.begin();
    rfid.PCD_Init();

    // Check RC522
    byte v = rfid.PCD_ReadRegister(MFRC522::VersionReg);
    if (v == 0x00 || v == 0xFF) {
        Serial.println("[ERROR] RC522 not detected! Check wiring.");
        while (true) {
            blinkRed(200);
            delay(200);
        }
    }
    Serial.printf("[OK] RC522 firmware version: 0x%02X\n", v);

    loadConfig();

    // If not configured, enter serial config mode
    if (wifiSSID.isEmpty() || apiURL.isEmpty()) {
        Serial.println();
        Serial.println("[SETUP] Device not configured. Enter config via serial.");
        Serial.println("Commands:");
        Serial.println("  SSID <your_wifi_name>");
        Serial.println("  PASS <your_wifi_password>");
        Serial.println("  API  <http://your-server:8080>");
        Serial.println("  SAVE");
        Serial.println("  SHOW (print current config)");
        Serial.println();
        serialConfig();
    }

    // Print device key for registration
    Serial.println();
    Serial.println("┌──────────────────────────────────────────────────────────────────────┐");
    Serial.println("│ DEVICE SECRET (enter this in the dashboard to register this device): │");
    Serial.printf( "│ %s │\n", deviceKey.c_str());
    Serial.println("└──────────────────────────────────────────────────────────────────────┘");
    Serial.println();

    connectWiFi();
    bootTime = millis();

    Serial.println("[READY] Waiting for card taps...");
    Serial.println();
}

// ── Main Loop ────────────────────────────────────────────────
void loop() {
    // Handle serial commands even during normal operation
    if (Serial.available()) {
        String line = Serial.readStringUntil('\n');
        line.trim();
        if (line.equalsIgnoreCase("SHOW")) {
            Serial.printf("SSID: %s\n", wifiSSID.c_str());
            Serial.printf("API:  %s\n", apiURL.c_str());
            Serial.printf("KEY:  %s\n", deviceKey.c_str());
            Serial.printf("WiFi: %s\n", WiFi.isConnected() ? "connected" : "disconnected");
        } else if (line.equalsIgnoreCase("REBOOT")) {
            ESP.restart();
        }
    }

    // Reconnect WiFi if dropped
    if (!WiFi.isConnected()) {
        Serial.println("[WIFI] Connection lost, reconnecting...");
        connectWiFi();
    }

    // Close door after timeout
    if (doorOpen && (millis() - doorOpenTime >= DOOR_OPEN_MS)) {
        closeDoor();
    }

    // Heartbeat
    if (millis() - lastHeartbeat >= HEARTBEAT_INTERVAL) {
        sendHeartbeat();
        lastHeartbeat = millis();
    }

    // Check for card
    if (!rfid.PICC_IsNewCardPresent() || !rfid.PICC_ReadCardSerial()) {
        delay(50);
        return;
    }

    String uid = readCardUID();
    rfid.PICC_HaltA();
    rfid.PCD_StopCrypto1();

    if (uid.isEmpty()) return;

    // Debounce: ignore same card within cooldown
    if (uid == lastCardUID && (millis() - lastCardTime < CARD_COOLDOWN_MS)) {
        return;
    }
    lastCardUID = uid;
    lastCardTime = millis();

    Serial.printf("[CARD] UID: %s\n", uid.c_str());
    checkIn(uid);
}

// ── Pin Setup ────────────────────────────────────────────────
void setupPins() {
    pinMode(RELAY_PIN, OUTPUT);
    pinMode(LED_GREEN_PIN, OUTPUT);
    pinMode(LED_RED_PIN, OUTPUT);
    pinMode(BUZZER_PIN, OUTPUT);

    digitalWrite(RELAY_PIN, LOW);
    digitalWrite(LED_GREEN_PIN, LOW);
    digitalWrite(LED_RED_PIN, LOW);
    digitalWrite(BUZZER_PIN, LOW);
}

// ── NVS Config ───────────────────────────────────────────────
void loadConfig() {
    prefs.begin(NVS_NAMESPACE, true); // read-only
    wifiSSID  = prefs.getString(NVS_KEY_WIFI_SSID, "");
    wifiPass  = prefs.getString(NVS_KEY_WIFI_PASS, "");
    apiURL    = prefs.getString(NVS_KEY_API_URL, "");
    deviceKey = prefs.getString(NVS_KEY_DEVICE_KEY, "");
    prefs.end();

    // Generate device key on first boot
    if (deviceKey.isEmpty()) {
        generateDeviceKey();
    }
}

void generateDeviceKey() {
    // Generate 32 random bytes → 64 hex chars
    char hex[65];
    for (int i = 0; i < 32; i++) {
        uint8_t b = (uint8_t)(esp_random() & 0xFF);
        sprintf(hex + (i * 2), "%02x", b);
    }
    hex[64] = '\0';
    deviceKey = String(hex);

    // Save to NVS
    prefs.begin(NVS_NAMESPACE, false);
    prefs.putString(NVS_KEY_DEVICE_KEY, deviceKey);
    prefs.end();

    Serial.println("[INIT] Generated new device secret.");
}

// ── Serial Configuration ─────────────────────────────────────
void serialConfig() {
    while (true) {
        // Blink red to indicate config mode
        blinkRed(500);

        if (!Serial.available()) {
            delay(100);
            continue;
        }

        String line = Serial.readStringUntil('\n');
        line.trim();
        if (line.isEmpty()) continue;

        if (line.startsWith("SSID ")) {
            wifiSSID = line.substring(5);
            wifiSSID.trim();
            Serial.printf("[OK] SSID set to: %s\n", wifiSSID.c_str());
        } else if (line.startsWith("PASS ")) {
            wifiPass = line.substring(5);
            wifiPass.trim();
            Serial.println("[OK] Password set.");
        } else if (line.startsWith("API ")) {
            apiURL = line.substring(4);
            apiURL.trim();
            // Remove trailing slash
            if (apiURL.endsWith("/")) apiURL.remove(apiURL.length() - 1);
            Serial.printf("[OK] API URL set to: %s\n", apiURL.c_str());
        } else if (line.equalsIgnoreCase("SHOW")) {
            Serial.printf("SSID: %s\n", wifiSSID.c_str());
            Serial.printf("API:  %s\n", apiURL.c_str());
            Serial.printf("KEY:  %s\n", deviceKey.c_str());
        } else if (line.equalsIgnoreCase("SAVE")) {
            if (wifiSSID.isEmpty() || apiURL.isEmpty()) {
                Serial.println("[ERROR] SSID and API URL are required before saving.");
                continue;
            }
            prefs.begin(NVS_NAMESPACE, false);
            prefs.putString(NVS_KEY_WIFI_SSID, wifiSSID);
            prefs.putString(NVS_KEY_WIFI_PASS, wifiPass);
            prefs.putString(NVS_KEY_API_URL, apiURL);
            prefs.putBool(NVS_KEY_CONFIGURED, true);
            prefs.end();
            Serial.println("[OK] Config saved. Continuing boot...");
            Serial.println();
            break;
        } else {
            Serial.println("[?] Unknown command. Use: SSID, PASS, API, SAVE, SHOW");
        }
    }
}

// ── WiFi ─────────────────────────────────────────────────────
void connectWiFi() {
    Serial.printf("[WIFI] Connecting to %s", wifiSSID.c_str());
    WiFi.mode(WIFI_STA);
    WiFi.begin(wifiSSID.c_str(), wifiPass.c_str());

    int attempts = 0;
    while (!WiFi.isConnected() && attempts < 20) {
        delay(500);
        Serial.print(".");
        attempts++;
    }

    if (WiFi.isConnected()) {
        Serial.printf(" connected! IP: %s\n", WiFi.localIP().toString().c_str());
        blinkGreen(200);
    } else {
        Serial.println(" FAILED!");
        Serial.println("[WIFI] Will retry in background...");
        blinkRed(1000);
    }
}

// ── Card Reading ─────────────────────────────────────────────
String readCardUID() {
    if (rfid.uid.size == 0) return "";

    String uid = "";
    for (byte i = 0; i < rfid.uid.size; i++) {
        if (rfid.uid.uidByte[i] < 0x10) uid += "0";
        uid += String(rfid.uid.uidByte[i], HEX);
    }
    uid.toUpperCase();
    return uid;
}

// ── API: Check-in ────────────────────────────────────────────
void checkIn(const String& cardUID) {
    if (!WiFi.isConnected()) {
        Serial.println("[ERROR] No WiFi — cannot check in.");
        beepError();
        blinkRed(500);
        return;
    }

    HTTPClient http;
    String url = apiURL + "/api/v1/devices/check-in";
    http.begin(url);
    http.addHeader("Content-Type", "application/json");
    http.addHeader("X-Device-Key", deviceKey);
    http.setTimeout(5000);

    // Build JSON body
    JsonDocument doc;
    doc["card_uid"] = cardUID;
    String body;
    serializeJson(doc, body);

    int httpCode = http.POST(body);
    String response = http.getString();
    http.end();

    if (httpCode == 201) {
        // Parse response for member name and streak
        JsonDocument resDoc;
        DeserializationError err = deserializeJson(resDoc, response);

        String memberName = "Member";
        int streak = 0;
        if (!err) {
            memberName = resDoc["member_name"].as<String>();
            streak = resDoc["streak"].as<int>();
        }

        Serial.printf("[CHECK-IN] %s (streak: %d days)\n", memberName.c_str(), streak);
        beepSuccess();
        openDoor();
    } else {
        // Parse error message
        JsonDocument errDoc;
        DeserializationError err = deserializeJson(errDoc, response);
        String errMsg = "unknown error";
        if (!err && errDoc["error"].is<const char*>()) {
            errMsg = errDoc["error"].as<String>();
        }

        Serial.printf("[DENIED] HTTP %d: %s\n", httpCode, errMsg.c_str());
        beepError();
        blinkRed(1000);
    }
}

// ── API: Heartbeat ───────────────────────────────────────────
void sendHeartbeat() {
    if (!WiFi.isConnected()) return;

    HTTPClient http;
    String url = apiURL + "/api/v1/devices/heartbeat";
    http.begin(url);
    http.addHeader("Content-Type", "application/json");
    http.addHeader("X-Device-Key", deviceKey);
    http.setTimeout(3000);

    unsigned long uptime = (millis() - bootTime) / 1000;

    JsonDocument doc;
    doc["uptime_seconds"] = (long)uptime;
    String body;
    serializeJson(doc, body);

    int httpCode = http.POST(body);
    http.end();

    if (httpCode == 200) {
        Serial.printf("[HEARTBEAT] OK (uptime: %lus)\n", uptime);
    } else {
        Serial.printf("[HEARTBEAT] Failed: HTTP %d\n", httpCode);
    }
}

// ── Door Control ─────────────────────────────────────────────
void openDoor() {
    digitalWrite(RELAY_PIN, HIGH);
    digitalWrite(LED_GREEN_PIN, HIGH);
    doorOpen = true;
    doorOpenTime = millis();
    Serial.println("[DOOR] Unlocked");
}

void closeDoor() {
    digitalWrite(RELAY_PIN, LOW);
    digitalWrite(LED_GREEN_PIN, LOW);
    doorOpen = false;
    Serial.println("[DOOR] Locked");
}

// ── Feedback ─────────────────────────────────────────────────
void beepSuccess() {
    // Two short beeps
    tone(BUZZER_PIN, 2000, 100);
    delay(150);
    tone(BUZZER_PIN, 2500, 100);
    delay(100);
    noTone(BUZZER_PIN);
}

void beepError() {
    // One long low beep
    tone(BUZZER_PIN, 400, 500);
    delay(500);
    noTone(BUZZER_PIN);
}

void blinkGreen(int ms) {
    digitalWrite(LED_GREEN_PIN, HIGH);
    delay(ms);
    digitalWrite(LED_GREEN_PIN, LOW);
}

void blinkRed(int ms) {
    digitalWrite(LED_RED_PIN, HIGH);
    delay(ms);
    digitalWrite(LED_RED_PIN, LOW);
}
