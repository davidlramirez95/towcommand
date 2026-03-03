import { SurgePricingCache } from '@towcommand/cache';

export interface SurgeAwareResult {
  providerId: string;
  surgeMultiplier: number;
  adjustedScore: number;
}

export class SurgeAwareAlgorithm {
  private surgeCache: SurgePricingCache;

  constructor() {
    this.surgeCache = new SurgePricingCache();
  }

  async calculate(
    providers: Array<{ providerId: string; score: number }>,
    region: string,
  ): Promise<SurgeAwareResult[]> {
    const surgeMultiplier = await this.surgeCache.getSurgeMultiplier(region);

    return providers
      .map(({ providerId, score }) => {
        // During surge, slightly boost providers with higher base scores
        // to incentivize reliable providers to accept jobs
        const surgeBonus = surgeMultiplier > 1.0 ? score * (surgeMultiplier - 1) * 0.1 : 0;
        const adjustedScore = Math.round((score + surgeBonus) * 100) / 100;

        return {
          providerId,
          surgeMultiplier,
          adjustedScore,
        };
      })
      .sort((a, b) => b.adjustedScore - a.adjustedScore);
  }
}
