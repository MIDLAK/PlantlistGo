package models

/* Местность */
type Terrain struct {
	City 		string
	Country 	string
	Locality 	string
	GPS			GPSPoint
}

/* GPS-координата */
type GPSPoint struct {
	ID			int
	Latitude	float64
	Longitude	float64
}

/* Страна */
type City struct {
	ID		int
	Name	string
}

/* Город */
type Country struct {
	ID		int
	Name	string
}

/* Участок */
type Locality struct {
	ID		int
	Name	string
}