export class TimeoutManager {
  private timeouts: Map<string, NodeJS.Timeout> = new Map();

  setTimeout(id: string, callback: () => void, delayMs: number): void {
    const existing = this.timeouts.get(id);
    if (existing) {
      clearTimeout(existing);
    }
    const timeout = setTimeout(callback, delayMs);
    this.timeouts.set(id, timeout);
  }

  clearTimeout(id: string): void {
    const timeout = this.timeouts.get(id);
    if (timeout) {
      clearTimeout(timeout);
      this.timeouts.delete(id);
    }
  }

  clearAll(): void {
    for (const timeout of this.timeouts.values()) {
      clearTimeout(timeout);
    }
    this.timeouts.clear();
  }
}
