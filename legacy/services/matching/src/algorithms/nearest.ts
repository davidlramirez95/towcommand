export interface NearestResult {
  providerId: string;
  distance: number;
  score: number;
}

export class NearestAlgorithm {
  calculate(providers: Array<{ providerId: string; distance: number }>): NearestResult[] {
    // Score inversely proportional to distance. Closer = higher score.
    const maxDistance = Math.max(...providers.map((p) => p.distance), 1);

    return providers
      .map((p) => ({
        providerId: p.providerId,
        distance: p.distance,
        score: Math.round(((1 - p.distance / (maxDistance + 1)) * 100) * 100) / 100,
      }))
      .sort((a, b) => b.score - a.score);
  }
}
