package repository

import "github.com/davidlramirez95/towcommand/internal/usecase/port"

// Compile-time interface satisfaction checks.
var (
	_ port.BookingSaver          = (*DynamoBookingRepository)(nil)
	_ port.BookingFinder         = (*DynamoBookingRepository)(nil)
	_ port.BookingByUserLister   = (*DynamoBookingRepository)(nil)
	_ port.BookingStatusUpdater  = (*DynamoBookingRepository)(nil)
	_ port.BookingByStatusLister = (*DynamoBookingRepository)(nil)

	_ port.UserSaver         = (*DynamoUserRepository)(nil)
	_ port.UserFinder        = (*DynamoUserRepository)(nil)
	_ port.UserByEmailFinder = (*DynamoUserRepository)(nil)
	_ port.UserByPhoneFinder = (*DynamoUserRepository)(nil)
	_ port.VehicleAdder      = (*DynamoUserRepository)(nil)
	_ port.VehicleLister     = (*DynamoUserRepository)(nil)

	_ port.ProviderSaver               = (*DynamoProviderRepository)(nil)
	_ port.ProviderFinder              = (*DynamoProviderRepository)(nil)
	_ port.ProviderByTierLister        = (*DynamoProviderRepository)(nil)
	_ port.ProviderLocationUpdater     = (*DynamoProviderRepository)(nil)
	_ port.ProviderAvailabilityUpdater = (*DynamoProviderRepository)(nil)
	_ port.ProviderDocSaver            = (*DynamoProviderRepository)(nil)
	_ port.ProviderDocLister           = (*DynamoProviderRepository)(nil)

	_ port.PaymentSaver           = (*DynamoPaymentRepository)(nil)
	_ port.PaymentFinder          = (*DynamoPaymentRepository)(nil)
	_ port.PaymentByBookingLister = (*DynamoPaymentRepository)(nil)
	_ port.PaymentStatusUpdater   = (*DynamoPaymentRepository)(nil)

	_ port.RatingSaver            = (*DynamoRatingRepository)(nil)
	_ port.RatingByBookingFinder  = (*DynamoRatingRepository)(nil)
	_ port.RatingByProviderLister = (*DynamoRatingRepository)(nil)

	_ port.EvidenceSaver           = (*DynamoEvidenceRepository)(nil)
	_ port.EvidenceByBookingLister = (*DynamoEvidenceRepository)(nil)
	_ port.MediaItemAdder          = (*DynamoEvidenceRepository)(nil)
)
