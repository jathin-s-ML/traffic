package generator

import (
	"encoding/json"
	"fmt"
	"math/rand"
)

func RandomData() []byte {
	data := map[string]interface{}{
		"id":    rand.Intn(1000),
		"value": rand.Float64() * 100,
		"info":  fmt.Sprintf("RandomInfo%d", rand.Intn(100)),
	}
	jsonData, _ := json.Marshal(data)
	return jsonData
}
