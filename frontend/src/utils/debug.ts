export function debugLog(component: string, message: string, data?: any) {
  if (import.meta.env.DEV) {
    console.log(`[${component}] ${message}`, data || '');
  }
}

export function safeAccess<T>(obj: any, path: string, defaultValue: T): T {
  try {
    return path.split('.').reduce((acc, key) => acc && acc[key], obj) || defaultValue;
  } catch {
    return defaultValue;
  }
}
