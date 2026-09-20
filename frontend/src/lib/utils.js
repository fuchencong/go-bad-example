import { clsx } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs) {
  return twMerge(clsx(inputs))
}

export function initials(value) {
  return value.split(" ").map((x) => x[0]).join("").slice(0, 2).toUpperCase()
}

export function shortDate(value) {
  if (!value) return "No date"
  return new Date(`${value}T12:00:00`).toLocaleDateString("en", { month: "short", day: "numeric" })
}

export function getWeeksInYear(year) {
  const first = new Date(year, 0, 1)
  const last = new Date(year, 11, 31)
  return Math.ceil(((last - first) / 86400000 + first.getDay() + 1) / 7)
}

export function scoreColor(n) {
  if (n > 90) return "red"
  if (n > 70) return "orange"
  if (n > 50) return "yellow"
  return "green"
}
