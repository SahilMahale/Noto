import { twMerge } from "tailwind-merge";
import clsx, { type ClassValue } from 'clsx';
// Utility function (often called cn)
export function classMerge(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}
