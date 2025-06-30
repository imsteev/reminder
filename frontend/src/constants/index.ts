// API Configuration
export const API_BASE_URL = "http://localhost:8080";

// User Configuration
export const DEFAULT_USER_ID = 1;

// Time Configuration
export const DEFAULT_START_DELAY_MINUTES = 5;

// Time Presets for Quick Selection
export const TIME_PRESETS = {
  NOW: 0,
  FIFTEEN_MINUTES: 15,
  ONE_HOUR: 60,
  TOMORROW_9AM: { hour: 9, minute: 0 },
  TOMORROW_1PM: { hour: 13, minute: 0 },
} as const;

// Form Configuration
export const FORM_DEFAULTS = {
  FREQUENCY_PLACEHOLDER: "3",
  INTERVAL_PLACEHOLDER: "8",
  MIN_FREQUENCY: 1,
  MIN_INTERVAL: 1,
} as const;

// UI Text
export const UI_TEXT = {
  START_TIME_LABEL: "When should we start sending reminders?",
  QUICK_SELECT_LABEL: "Quick select:",
  CUSTOM_TIME_LABEL: "Date (local time)",
  PHONE_PLACEHOLDER: "+1234567890",
  MESSAGE_PLACEHOLDER: "Take your medication",
} as const;
