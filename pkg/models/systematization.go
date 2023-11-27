package models

/* Систематика */
type Systematization struct {
	Domain		Domain	
	Kingdom		Kingdom
	Order		Order
	Family		Family
	Genus		Genus
	Department	Department
	Class		Class
}

/* Домен */
type Domain struct {
	ID		int
	Name	string
}

/* Царство */
type Kingdom struct {
	ID		int
	Name	string
}

/* Порядок */
type Order struct {
	ID		int
	Name	string
}

/* Семейство */
type Family struct {
	ID		int
	Name	string
}

/* Род */
type Genus struct {
	ID		int
	Name	string
}

/* Отдел */
type Department struct {
	ID		int
	Name	string
}

/* Класс */
type Class struct {
	ID		int
	Name	string
}