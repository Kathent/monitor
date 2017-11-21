package db

func LoadDb(){
	GetRedisClient()
	LoadSession()
}

func Init() {
	LoadDb()
}
