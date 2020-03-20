package Database

func AddSessionAmount() {
	var status Status
	DB := GetDatabase()
	DB.First(&status, "id = ?", 1)
	DB.Model(&status).Update("sessions", status.Sessions+1)
}

func AddItemAmount() {
	var status Status
	DB := GetDatabase()
	DB.First(&status, "id = ?", 1)
	DB.Model(&status).Update("items", status.Items+1)
}

func GetStats() *Status {
	var status Status
	DB := GetDatabase()
	DB.First(&status, "id = ?", 1)
	return &status
}
