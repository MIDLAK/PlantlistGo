package models

/* Растение */
type Plant struct {
	ID				int
	RusName			string
	LatinName		string
	Description		string
	Image			string	//путь к файлу
	IDClass			int
	IDDepartment	int
	IDDomain		int
	IDFamily		int
	IDGenus			int
	IDKingdom		int
	IDOrder			int
	IDStatus		int
}

//превью таксона
type PlantPreview struct {
	Name		string
	Description	string
}