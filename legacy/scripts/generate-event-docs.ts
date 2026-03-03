// Auto-generate event catalog documentation from code
const EVENT_CATALOG = {
  'tc.booking': ['BookingCreated', 'BookingAccepted', 'BookingCancelled', 'BookingCompleted', 'BookingStatusChanged'],
  'tc.matching': ['ProviderMatched', 'MatchTimeout', 'MatchReassigned'],
  'tc.tracking': ['LocationUpdated', 'DriverArrived', 'RouteDeviation'],
  'tc.payment': ['PaymentInitiated', 'PaymentCompleted', 'PaymentFailed', 'RefundProcessed'],
  'tc.sos': ['SOSActivated', 'SOSResolved'],
  'tc.auth': ['UserRegistered', 'UserSuspended'],
  'tc.provider': ['ProviderOnline', 'ProviderOffline', 'ProviderVerified'],
  'tc.evidence': ['EvidenceUploaded', 'ReportCompleted'],
};

function generateDocs() {
  let md = '# TowCommand PH — Event Catalog\n\n';
  md += `Generated: ${new Date().toISOString()}\n\n`;
  md += '| Source | Event | Description |\n|--------|-------|-------------|\n';

  for (const [source, events] of Object.entries(EVENT_CATALOG)) {
    for (const event of events) {
      md += `| \`${source}\` | \`${event}\` | TODO: Add description |\n`;
    }
  }

  md += `\n\nTotal: ${Object.values(EVENT_CATALOG).flat().length} events across ${Object.keys(EVENT_CATALOG).length} sources\n`;

  const fs = require('fs');
  fs.writeFileSync('docs/event-catalog.md', md);
  console.log('✅ Event catalog docs generated at docs/event-catalog.md');
}

generateDocs();
