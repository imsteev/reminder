/**
 * Converts a Date object to a datetime-local string format
 * Handles timezone offset for local datetime input
 */
export function toLocalDateTimeString(date: Date): string {
  const localDateTime = new Date(
    date.getTime() - date.getTimezoneOffset() * 60000
  )
    .toISOString()
    .slice(0, 16);
  return localDateTime;
}

/**
 * Creates a new Date that is X minutes from now
 */
export function addMinutesToNow(minutes: number): Date {
  return new Date(Date.now() + minutes * 60000);
}

/**
 * Creates a Date for tomorrow at a specific hour and minute
 */
export function getTomorrowAtTime(hour: number, minute = 0): Date {
  const tomorrow = new Date();
  tomorrow.setDate(tomorrow.getDate() + 1);
  tomorrow.setHours(hour, minute, 0, 0);
  return tomorrow;
}

/**
 * Gets the current time as a datetime-local string
 */
export function getCurrentDateTimeString(): string {
  return toLocalDateTimeString(new Date());
}

/**
 * Gets a datetime-local string for X minutes from now
 */
export function getDateTimeStringInMinutes(minutes: number): string {
  return toLocalDateTimeString(addMinutesToNow(minutes));
}

/**
 * Gets a datetime-local string for tomorrow at a specific time
 */
export function getTomorrowDateTimeString(hour: number, minute = 0): string {
  return toLocalDateTimeString(getTomorrowAtTime(hour, minute));
}
