package db

type Query interface {
	Query()
}

type RedisQuery struct {

}

type MongoDbQuery struct {
	TableName string
	FindMap map[string]interface{}
}
