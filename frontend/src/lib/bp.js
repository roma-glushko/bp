export const MIN_SYSTOLIC = 70
export const MAX_SYSTOLIC = 250
export const MIN_DIASTOLIC = 40
export const MAX_DIASTOLIC = 150
export const MIN_PULSE = 30
export const MAX_PULSE = 220

export function classifyBP(systolic, diastolic) {
  if (systolic >= 180 || diastolic >= 120) {
    return { status: 'very_high', color: 'darkred', label: 'Very high' }
  }
  if (systolic >= 140 || diastolic >= 90) {
    return { status: 'high', color: 'red', label: 'Too high' }
  }
  if (systolic >= 120 || diastolic >= 80) {
    return { status: 'elevated', color: 'yellow', label: 'Elevated' }
  }
  return { status: 'great', color: 'green', label: 'Great' }
}

export function isValidSystolic(v) {
  return v >= MIN_SYSTOLIC && v <= MAX_SYSTOLIC
}

export function isValidDiastolic(v) {
  return v >= MIN_DIASTOLIC && v <= MAX_DIASTOLIC
}

export function isValidPulse(v) {
  return v >= MIN_PULSE && v <= MAX_PULSE
}
