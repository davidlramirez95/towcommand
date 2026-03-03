import type { EventBridgeEvent } from 'aws-lambda';
import type { BookingCreatedEvent } from '@towcommand/core';
import { estimateEtaMinutes } from '@towcommand/core';
import { BookingRepository, ProviderRepository } from '@towcommand/db';
import { publishEvent, EVENT_CATALOG } from '@towcommand/events';
import { RateLimiter } from '@towcommand/cache';
import { GeoSearchService } from './lib/geo-search';
import { WeightedScoreAlgorithm } from './algorithms/weighted-score';
import { SurgeAwareAlgorithm } from './algorithms/surge-aware';

const geoSearch = new GeoSearchService();
const weightedScore = new WeightedScoreAlgorithm();
const surgeAware = new SurgeAwareAlgorithm();
const bookingRepo = new BookingRepository();
const providerRepo = new ProviderRepository();
const rateLimiter = new RateLimiter();

const SEARCH_RADIUS_KM = [5, 10, 20, 30]; // Cascade: expand radius on each attempt
const MATCH_TIMEOUT_MS = 60_000; // 60 seconds per cascade

interface MatchingEventDetail {
  bookingId: string;
  customerId: string;
  serviceType: string;
  pickupLocation: { lat: number; lng: number; address?: string };
  dropoffLocation: { lat: number; lng: number; address?: string };
  price: Record<string, unknown>;
}

export async function handler(
  event: EventBridgeEvent<string, MatchingEventDetail>,
): Promise<void> {
  const detailType = event['detail-type'];
  const detail = event.detail;

  try {
    if (detailType === EVENT_CATALOG.booking.events.BookingCreated) {
      await handleBookingCreated(detail);
    } else if (detailType === EVENT_CATALOG.matching.events.MatchTimeout) {
      await handleMatchTimeout(detail as any);
    }
  } catch (error) {
    console.error(`Matching handler error for ${detailType}:`, error);
    throw error;
  }
}

async function handleBookingCreated(detail: MatchingEventDetail): Promise<void> {
  const { bookingId, pickupLocation } = detail;

  // Acquire lock to prevent duplicate matching
  const lockAcquired = await rateLimiter.acquireJobLock(bookingId, MATCH_TIMEOUT_MS / 1000);
  if (!lockAcquired) {
    console.log(`Matching already in progress for booking ${bookingId}`);
    return;
  }

  await findAndMatchProvider(bookingId, pickupLocation.lat, pickupLocation.lng, 0);
}

async function findAndMatchProvider(
  bookingId: string,
  lat: number,
  lng: number,
  cascade: number,
): Promise<void> {
  const radiusKm = SEARCH_RADIUS_KM[cascade] ?? SEARCH_RADIUS_KM[SEARCH_RADIUS_KM.length - 1];

  // Find nearby online providers
  const nearbyProviders = await geoSearch.findNearbyProviders(lat, lng, radiusKm);

  if (nearbyProviders.length === 0) {
    if (cascade < SEARCH_RADIUS_KM.length - 1) {
      // Escalate: publish timeout event to trigger wider search
      await publishEvent(
        EVENT_CATALOG.matching.source,
        EVENT_CATALOG.matching.events.MatchTimeout,
        { bookingId, cascade: cascade + 1, attemptedProviders: [], lat, lng },
      );
      return;
    }

    // No providers found after all cascades
    console.log(`No providers found for booking ${bookingId} after ${cascade + 1} cascades`);
    return;
  }

  // Score providers with weighted algorithm
  const scored = weightedScore.calculate(
    nearbyProviders.map((np) => ({ provider: np.provider, distance: np.distance })),
    lat, lng,
  );

  // Apply surge-aware adjustments
  const surgeAdjusted = await surgeAware.calculate(
    scored.map((s) => ({ providerId: s.providerId, score: s.score })),
    'NCR',
  );

  // Select the best provider
  const bestMatch = surgeAdjusted[0];
  if (!bestMatch) return;

  const matchedProvider = await providerRepo.getById(bestMatch.providerId);
  if (!matchedProvider) return;

  const distanceToPickup = nearbyProviders.find(
    (p) => p.providerId === bestMatch.providerId,
  )?.distance ?? 0;

  // Update booking with matched provider
  await bookingRepo.updateStatus(bookingId, 'MATCHED' as any, {
    providerId: bestMatch.providerId,
    matchedAt: new Date().toISOString(),
  });

  // Publish match event
  await publishEvent(
    EVENT_CATALOG.matching.source,
    EVENT_CATALOG.matching.events.ProviderMatched,
    {
      bookingId,
      providerId: matchedProvider.providerId,
      providerName: matchedProvider.name,
      providerPhone: matchedProvider.phone,
      truckPlate: matchedProvider.plateNumber,
      eta: estimateEtaMinutes(distanceToPickup),
      score: bestMatch.adjustedScore,
    },
  );
}

async function handleMatchTimeout(detail: {
  bookingId: string;
  cascade: number;
  attemptedProviders: string[];
  lat: number;
  lng: number;
}): Promise<void> {
  // Verify booking is still pending
  const booking = await bookingRepo.getById(detail.bookingId);
  if (!booking || booking.status !== ('PENDING' as any)) return;

  await findAndMatchProvider(detail.bookingId, detail.lat, detail.lng, detail.cascade);
}
