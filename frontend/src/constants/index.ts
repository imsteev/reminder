// API Configuration
export const API_BASE_URL =
  import.meta.env.VITE_API_URL || "http://localhost:8080";

// User Configuration
export const DEFAULT_USER_ID = 1;

// Time Configuration
export const DEFAULT_START_DELAY_MINUTES = 5;

// Time Presets for Quick Selection
export const TIME_PRESETS = {
  NOW: 0,
  FIVE_MINUTES: 5,
  FIFTEEN_MINUTES: 15,
  ONE_HOUR: 60,
  TOMORROW_9AM: { hour: 9, minute: 0 },
  TOMORROW_1PM: { hour: 13, minute: 0 },
} as const;
