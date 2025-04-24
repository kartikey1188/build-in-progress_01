package types

type CollectorUpdate struct {
	FullName      *string `json:"full_name,omitempty"`
	PhoneNumber   *string `json:"phone_number,omitempty"`
	Address       *string `json:"address,omitempty"`
	ProfileImage  *string `json:"profile_image,omitempty"`
	CompanyName   *string `json:"company_name,omitempty"`
	LicenseNumber *string `json:"license_number,omitempty"`
	Capacity      *int64  `json:"capacity,omitempty"`
	LicenseExpiry *Date   `json:"license_expiry,omitempty"`
}
