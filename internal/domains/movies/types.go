package movies

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Movie struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Year        int                `bson:"year" json:"year"`
	Genre       []string           `bson:"genre" json:"genre"`
	Director    []string           `bson:"director" json:"director"`
	Description string             `bson:"description" json:"description"`
	Duration    int                `bson:"duration_minutes" json:"duration_minutes"`

	// Rich metadata that benefits from document storage
	Cast []CastMember `bson:"cast" json:"cast"`
	Crew []CrewMember `bson:"crew" json:"crew"`

	// Production details
	Production ProductionInfo `bson:"production" json:"production"`

	// Technical specifications
	Technical TechnicalInfo `bson:"technical" json:"technical"`

	// External ratings and reviews
	ExternalRatings map[string]float64 `bson:"external_ratings" json:"external_ratings"` // IMDb, RT, etc.

	// Metadata
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type CastMember struct {
	Name      string `bson:"name" json:"name"`
	Character string `bson:"character" json:"character"`
	Order     int    `bson:"order" json:"order"`
}

type CrewMember struct {
	Name string `bson:"name" json:"name"`
	Role string `bson:"role" json:"role"` // Director, Producer, Writer, etc.
}

type ProductionInfo struct {
	Budget      int64    `bson:"budget" json:"budget"`
	BoxOffice   int64    `bson:"box_office" json:"box_office"`
	Studios     []string `bson:"studios" json:"studios"`
	Countries   []string `bson:"countries" json:"countries"`
	Languages   []string `bson:"languages" json:"languages"`
	ReleaseDate string   `bson:"release_date" json:"release_date"`
}

type TechnicalInfo struct {
	AspectRatio string   `bson:"aspect_ratio" json:"aspect_ratio"`
	SoundMix    []string `bson:"sound_mix" json:"sound_mix"`
	Color       string   `bson:"color" json:"color"`
	Camera      []string `bson:"camera" json:"camera"`
}

// IDString converts MongoDB ObjectID to string for cross-database references
func (m *Movie) IDString() string {
	return m.ID.Hex()
}
