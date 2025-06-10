package models

type UniversityType string
const (
	StateUniversity  UniversityType = "state" // Devlet Üniversitesi
	PrivateUniversity UniversityType = "private" // Özel Üniversite
	KktcUniversity UniversityType = "kktc" // KKTC Üniversiteleri
	AbroadUniversity UniversityType = "abroad" // Yurtdışı Üniversiteleri
)

type DegreeLevel string
const (
	Licence    DegreeLevel = "licence" // Lisans (4 Yıl)
	Associate  DegreeLevel = "associate" // Ön Lisans (2 Yıl)
)

type ScoreType string
const (
	Numerical     ScoreType = "numerical" // Sayısal Bölümü
	Verbal        ScoreType = "verbal" // Sözel Bölümü
	EqualWeight   ScoreType = "equal_weight" // Eşit Ağırlık Bölümü
	TYT           ScoreType = "tyt" // TYT Bölümü
	Language      ScoreType = "language" // Dil Bölümü
)

type EducationType string
const (
	Formal    EducationType = "formal" // Örgün Eğitim
	SecondEducation    EducationType = "second_education" // İkinci Eğitim
	Distance  EducationType = "distance" // Uzaktan Eğitim
)

// Rektörün Bilgileri
type Rector struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

// Okulun İletişim Bilgileri
type Contacts struct {
	Phone   string  `json:"phone"`
	Faks    string  `json:"faks"`
	Website string  `json:"website"`
	Mail    string  `json:"mail"`
	Address string  `json:"address"`
	Rector  Rector  `json:"rector"`
}

// Yıllara göre sıralama puan gibi veriler
type YearlyData struct {
	Year       int     `json:"year"`
	Quota      int     `json:"quota"`
	BaseScore  float64 `json:"base_score"`
	TopScore   float64 `json:"top_score"`
	BaseRank   int     `json:"base_rank"`
	TopRank    int     `json:"top_rank"`
	Placement  int     `json:"placement"`
}

// Üniversiteye ait departmanların bilgileri
type Department struct {
	Name          string        `json:"name"`
	Slug          string        `json:"slug,omitempty"`
	Faculty       string        `json:"faculty,omitempty"`
	Language      string        `json:"language"`
	Duration      string        `json:"duration"`    
	DegreeLevel   DegreeLevel   `json:"degree_level"` 
	ScoreType     ScoreType     `json:"score_type"`    
	EducationType EducationType `json:"education_type"` 
	YearlyData    []YearlyData  `json:"yearly_data"`
}

type University struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Slug           string         `json:"slug,omitempty"`
	City           string         `json:"city"`
	District       string         `json:"district"`
	UniversityType UniversityType `json:"university_type"` 
	Contacts       *Contacts      `json:"contacts,omitempty"`
	Departments    []Department   `json:"departments"`
}
