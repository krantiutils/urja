#pragma once

// ── Pin Assignments (ESP32-WROOM + RC522) ─────────────────────
// RC522 SPI pins
#define RC522_SS_PIN    5
#define RC522_RST_PIN   22
// SPI defaults: MOSI=23, MISO=19, SCK=18

// Output pins
#define RELAY_PIN       4    // Door lock relay (active HIGH)
#define LED_GREEN_PIN   2    // Success indicator
#define LED_RED_PIN     15   // Error indicator
#define BUZZER_PIN      13   // Piezo buzzer

// ── Timing ────────────────────────────────────────────────────
#define DOOR_OPEN_MS        5000   // How long the door stays unlocked
#define HEARTBEAT_INTERVAL  60000  // Heartbeat every 60 seconds
#define CARD_COOLDOWN_MS    3000   // Ignore same card re-tap for 3s
#define WIFI_RETRY_DELAY    5000   // Retry WiFi every 5s

// ── NVS Keys ──────────────────────────────────────────────────
#define NVS_NAMESPACE       "urja"
#define NVS_KEY_WIFI_SSID   "wifi_ssid"
#define NVS_KEY_WIFI_PASS   "wifi_pass"
#define NVS_KEY_API_URL     "api_url"
#define NVS_KEY_DEVICE_KEY  "device_key"
#define NVS_KEY_CONFIGURED  "configured"
