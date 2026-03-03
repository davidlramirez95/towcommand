export const EVENT_CATALOG = {
  booking: {
    source: 'tc.booking' as const,
    events: {
      BookingCreated: 'BookingCreated',
      BookingAccepted: 'BookingAccepted',
      BookingCancelled: 'BookingCancelled',
      BookingCompleted: 'BookingCompleted',
      BookingStatusChanged: 'BookingStatusChanged',
    },
  },
  matching: {
    source: 'tc.matching' as const,
    events: {
      ProviderMatched: 'ProviderMatched',
      MatchTimeout: 'MatchTimeout',
      MatchReassigned: 'MatchReassigned',
    },
  },
  tracking: {
    source: 'tc.tracking' as const,
    events: {
      LocationUpdated: 'LocationUpdated',
      DriverArrived: 'DriverArrived',
      RouteDeviation: 'RouteDeviation',
    },
  },
  payment: {
    source: 'tc.payment' as const,
    events: {
      PaymentInitiated: 'PaymentInitiated',
      PaymentCompleted: 'PaymentCompleted',
      PaymentFailed: 'PaymentFailed',
      RefundProcessed: 'RefundProcessed',
    },
  },
  sos: {
    source: 'tc.sos' as const,
    events: {
      SOSActivated: 'SOSActivated',
      SOSResolved: 'SOSResolved',
    },
  },
  auth: {
    source: 'tc.auth' as const,
    events: {
      UserRegistered: 'UserRegistered',
      UserSuspended: 'UserSuspended',
    },
  },
  provider: {
    source: 'tc.provider' as const,
    events: {
      ProviderOnline: 'ProviderOnline',
      ProviderOffline: 'ProviderOffline',
      ProviderVerified: 'ProviderVerified',
    },
  },
  evidence: {
    source: 'tc.evidence' as const,
    events: {
      EvidenceUploaded: 'EvidenceUploaded',
      ReportCompleted: 'ReportCompleted',
    },
  },
} as const;
