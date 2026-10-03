package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateNomorRekening() string {
	waktuSekarang := time.Now().Format("020106")

	angkaAcak := rand.Intn(900000) + 100000

	return fmt.Sprintf("110%s%06d", waktuSekarang, angkaAcak)
}
