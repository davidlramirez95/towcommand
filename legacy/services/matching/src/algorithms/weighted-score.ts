import type { Provider } from '@towcommand/core';
import { haversineDistance } from '@towcommand/core';

export interface WeightedScoreResult {
  providerId: string;
  score: number;
  factors: {
    distance: number;
    rating: number;
    acceptance: number;
    experience: number;
  };
}

// Matching Algorithm Weights (from TRS):
// Distance: 40%, Trust/Rating: 25%, Acceptance Rate: 15%, Vehicle Compatibility: 10%, Current Load: 10%
const WEIGHTS = {
  distance: 0.40,
  rating: 0.25,
  acceptance: 0.15,
  experience: 0.20,
};

const MAX_DISTANCE_KM = 30;
const MAX_RATING = 5;
const MAX_JOBS = 500;

export class WeightedScoreAlgorithm {
  calculate(
    providers: Array<{ provider: Provider; distance: number }>,
    pickupLat: number,
    pickupLng: number,
  ): WeightedScoreResult[] {
    return providers
      .map(({ provider, distance }) => {
        // Distance score: closer is better (0-100)
        const distanceScore = Math.max(0, (1 - distance / MAX_DISTANCE_KM)) * 100;

        // Rating score: higher is better (0-100)
        const ratingScore = (provider.rating / MAX_RATING) * 100;

        // Acceptance rate score: already a percentage (0-100)
        const acceptanceScore = provider.acceptanceRate;

        // Experience score: more jobs = more experienced (0-100, capped)
        const experienceScore = Math.min(provider.totalJobsCompleted / MAX_JOBS, 1) * 100;

        const score = Math.round(
          (WEIGHTS.distance * distanceScore +
           WEIGHTS.rating * ratingScore +
           WEIGHTS.acceptance * acceptanceScore +
           WEIGHTS.experience * experienceScore) * 100,
        ) / 100;

        return {
          providerId: provider.providerId,
          score,
          factors: {
            distance: Math.round(distanceScore * 100) / 100,
            rating: Math.round(ratingScore * 100) / 100,
            acceptance: Math.round(acceptanceScore * 100) / 100,
            experience: Math.round(experienceScore * 100) / 100,
          },
        };
      })
      .sort((a, b) => b.score - a.score);
  }
}
