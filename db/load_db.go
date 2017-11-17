package db

func LoadDb(){
	GetRedisClient()
	LoadSession()
}

func init() {
	LoadDb()
}
