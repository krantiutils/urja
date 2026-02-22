-- Add API key to NFC devices for device-level authentication (ESP32 readers).
ALTER TABLE nfc_devices ADD COLUMN api_key VARCHAR(64) UNIQUE;
