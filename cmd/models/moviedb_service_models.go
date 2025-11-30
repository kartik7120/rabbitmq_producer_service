package models

import (
	"time"
)

type CastAndCrew struct {
	Type      string `json:"type" gorm:"not null"` // cast or crew
	Name      string `json:"name" gorm:"not null"`
	Character string `json:"character"`
	PhotoURL  string `json:"photo_url"`
	MovieID   uint
}

type SeatMatrix struct {
	SeatNumber string `json:"seat_number" gorm:"not null;uniqueIndex:idx_unique_seat"`
	Row        int    `json:"row" gorm:"not null;uniqueIndex:idx_unique_seat"`
	Column     int    `json:"column" gorm:"not null;uniqueIndex:idx_unique_seat"`
	Price      int    `json:"price" gorm:"not null;uniqueIndex:idx_unique_seat"`
	VenueID    uint   `json:"venue_id" gorm:"not null;uniqueIndex:idx_unique_seat"`
	Type       string `json:"type"`
}

type BookedSeats struct {
	SeatNumber      string     `json:"seat_number" gorm:"not null;uniqueIndex:idx_unique_booked_seats"`
	MovieTimeSlotID uint       `json:"movie_time_slot_id" gorm:"not null;uniqueIndex:idx_unique_booked_seats"` // Link booking to a movie show
	SeatMatrixID    uint       `json:"seat_matrix_id" gorm:"not null;uniqueIndex:idx_unique_booked_seats"`     // Reference seat matrix for consistency
	IsBooked        bool       `json:"is_booked" gorm:"default:false"`
	Email           *string    `json:"email" validate:"required,email"`
	PhoneNumber     string     `json:"phone_number" validate:"required,e164"`
	LockedUntil     *time.Time `json:"locked_until"` // Optional field to lock the seat for a certain period
}

type Review struct {
	MovieID         uint       `json:"movie_id" gorm:"not null;index"` // FK to movies table
	UserID          int        `json:"user_id" gorm:"index"`           // FK to users table (-1 for anonymous)
	Title           string     `json:"title" gorm:"type:varchar(255);not null"`
	Comment         string     `json:"comment" gorm:"type:text"`
	Rating          int        `json:"rating" gorm:"not null;check:rating >= 1 AND rating <= 5"` // 1–5 only
	ContainsSpoiler bool       `json:"contains_spoiler" gorm:"default:false"`                    // marks if review has spoilers
	Language        string     `json:"language" gorm:"type:varchar(50);default:'English'"`       // e.g., Hindi, Tamil, etc.
	Format          string     `json:"format" gorm:"type:varchar(20);default:'2D'"`              // 2D, 3D, IMAX, 4DX, etc.
	Likes           int        `json:"likes" gorm:"default:0"`                                   // count of likes (can be incremented)
	Dislikes        int        `json:"dislikes" gorm:"default:0"`                                // count of dislikes
	WatchedAt       *time.Time `json:"watched_at"`                                               // when user watched the movie (optional)
	IsEdited        bool       `json:"is_edited" gorm:"default:false"`                           // if review was updated after posting
}

type MovieTimeSlot struct {
	StartTime   time.Time `json:"start_time" gorm:"not null"`
	EndTime     time.Time `json:"end_time" gorm:"not null"`
	Duration    int       `json:"duration" gorm:"not null"` // in minutes
	MovieID     uint      `json:"movie_id"`
	Date        time.Time `json:"date" gorm:"not null"`
	MovieFormat string    `json:"movie_format" gorm:"not null"` // movie format (e.g., 2D, 3D)
	VenueID     uint      `json:"venue_id"`
}

type Movie struct {
	ID uint `gorm:"primaryKey"`
	// CreatedAt           time.Time
	// UpdatedAt           time.Time
	// DeletedAt           gorm.DeletedAt `gorm:"index"`
	Title               string        `json:"title" gorm:"not null;unique"`
	Description         string        `json:"description" gorm:"not null"`
	Duration            int           `json:"duration" gorm:"not null"`
	Language            []string      `json:"language" gorm:"type:text[];not null"`
	Type                []string      `json:"type" gorm:"type:text[];not null"`
	CastCrew            []CastAndCrew `json:"cast_crew" gorm:"foreignKey:MovieID"`
	PosterURL           string        `json:"poster_url"`
	TrailerURL          string        `json:"trailer_url"`
	ReleaseDate         time.Time     `json:"release_date" gorm:"not null"`
	MovieResolution     []string      `json:"movie_resolution" gorm:"type:text[];not null"`
	Venues              []Venue       `json:"venues" gorm:"many2many:movie_venues;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Ranking             uint          `json:"ranking"`
	Votes               uint          `json:"votes"`
	Reviews             []Review      `json:"reviews" gorm:"foreignKey:MovieID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ScreenWidePosterURL string        `json:"screen_wide_poster_url" gorm:"null,default:null"` // URL for a wide poster suitable for screens
	LogoImageURL        string        `json:"logo_image_url" gorm:"null"`
}

type Venue struct {
	Name                 string   `json:"name" gorm:"not null"`
	Type                 string   `json:"type" gorm:"not null"`
	Address              string   `json:"address" gorm:"not null"`
	Rows                 int      `json:"rows" gorm:"not null"`
	Columns              int      `json:"columns" gorm:"not null"`
	ScreenNumber         int      `json:"screen_number" gorm:"not null;unique"`
	Longitude            float64  `json:"longitude" gorm:"not null"`
	Latitude             float64  `json:"latitude" gorm:"not null"`
	MovieFormatSupported []string `json:"movie_format_supported" gorm:"type:text[];not null"`
	LanguagesSupported   []string `json:"languages_supported" gorm:"type:text[];not null"`

	// Relationships
	Seats          []SeatMatrix    `json:"seats" gorm:"foreignKey:VenueID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	MovieTimeSlots []MovieTimeSlot `json:"movie_time_slots" gorm:"foreignKey:VenueID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Movies         []Movie         `json:"movies" gorm:"many2many:movie_venues;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CinemaName     string          `json:"cinema_name" gorm:"type:text;default:null"`
}
