const EARTH_RADIUS_KM = 6371;

export function haversineDistance(
  lat1: number, lng1: number,
  lat2: number, lng2: number,
): number {
  const dLat = toRadians(lat2 - lat1);
  const dLng = toRadians(lng2 - lng1);
  const a =
    Math.sin(dLat / 2) * Math.sin(dLat / 2) +
    Math.cos(toRadians(lat1)) * Math.cos(toRadians(lat2)) *
    Math.sin(dLng / 2) * Math.sin(dLng / 2);
  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
  return EARTH_RADIUS_KM * c;
}

function toRadians(degrees: number): number {
  return degrees * (Math.PI / 180);
}

export function isWithinRadius(
  centerLat: number, centerLng: number,
  pointLat: number, pointLng: number,
  radiusKm: number,
): boolean {
  return haversineDistance(centerLat, centerLng, pointLat, pointLng) <= radiusKm;
}

export function isValidPhilippineCoordinate(lat: number, lng: number): boolean {
  return lat >= 4.5 && lat <= 21.5 && lng >= 116.0 && lng <= 127.0;
}

export function estimateEtaMinutes(distanceKm: number, isMetroManila = true): number {
  const avgSpeedKmh = isMetroManila ? 15 : 40;
  return Math.ceil((distanceKm / avgSpeedKmh) * 60);
}

export interface BoundingBox {
  minLat: number; maxLat: number;
  minLng: number; maxLng: number;
}

export function getBoundingBox(lat: number, lng: number, radiusKm: number): BoundingBox {
  const latDelta = radiusKm / 111.32;
  const lngDelta = radiusKm / (111.32 * Math.cos(toRadians(lat)));
  return {
    minLat: lat - latDelta, maxLat: lat + latDelta,
    minLng: lng - lngDelta, maxLng: lng + lngDelta,
  };
}
