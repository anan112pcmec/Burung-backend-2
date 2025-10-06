package helper

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/meilisearch/meilisearch-go"

	"github.com/anan112pcmec/Burung-backend-2/watcher_app/database/models"
)

func ConvertJenisBarang(jenis string) string {
	// Map internal -> DB
	mapJenis := map[string]string{
		"Pakaian&Fashion":     "Pakaian & Fashion",
		"Kosmetik&Kecantikan": "Kosmetik & Kecantikan",
		"Elektronik&Gadget":   "Elektronik & Gadget",
		"Buku&Media":          "Buku & Media",
		"Makanan&Minuman":     "Makanan & Minuman",
		"Ibu&Bayi":            "Ibu & Bayi",
		"Mainan":              "Mainan",
		"Olahraga&Outdoor":    "Olahraga & Outdoor",
		"Otomotif&Sparepart":  "Otomotif & Sparepart",
		"RumahTangga":         "Rumah Tangga",
		"AlatTulis":           "Alat Tulis",
		"Perhiasan&Aksesoris": "Perhiasan & Aksesoris",
		"ProdukDigital":       "ProdukDigital",
		"Bangunan&Perkakas":   "Bangunan & Perkakas",
		"Musik&Instrumen":     "Musik & Instrumen",
		"Film&Broadcasting":   "Film & Broadcasting",
		"SemuaBarang":         "Semua Barang",
	}

	if val, ok := mapJenis[jenis]; ok {
		return val
	}
	// fallback kalau tidak ada mapping
	return jenis
}

func ConvertJenisBarangReverse(jenis string) string {
	// Map DB -> internal
	mapReverse := map[string]string{
		"Pakaian & Fashion":     "Pakaian&Fashion",
		"Kosmetik & Kecantikan": "Kosmetik&Kecantikan",
		"Elektronik & Gadget":   "Elektronik&Gadget",
		"Buku & Media":          "Buku&Media",
		"Makanan & Minuman":     "Makanan&Minuman",
		"Ibu & Bayi":            "Ibu&Bayi",
		"Mainan":                "Mainan",
		"Olahraga & Outdoor":    "Olahraga&Outdoor",
		"Otomotif & Sparepart":  "Otomotif&Sparepart",
		"Rumah Tangga":          "RumahTangga",
		"Alat Tulis":            "AlatTulis",
		"Perhiasan & Aksesoris": "Perhiasan&Aksesoris",
		"ProdukDigital":         "ProdukDigital",
		"Bangunan & Perkakas":   "Bangunan&Perkakas",
		"Musik & Instrumen":     "Musik&Instrumen",
		"Film & Broadcasting":   "Film&Broadcasting",
		"Semua Barang":          "SemuaBarang",
	}

	if val, ok := mapReverse[jenis]; ok {
		return val
	}
	// fallback kalau tidak ada mapping
	return jenis
}

func BarangDataFromSearchEngine(data meilisearch.Hits) []models.BarangInduk {
	var hasilnyo []models.BarangInduk

	hitsjson, err := json.Marshal(data)
	if err != nil {
		log.Fatal("❌ Gagal marshal hits:", err)
	}

	if err := json.Unmarshal(hitsjson, &hasilnyo); err != nil {
		log.Fatal("❌ Gagal unmarshal hits ke struct:", err)
	}

	return hasilnyo
}

func ShouldDelete(inputTimeStr string) bool {
	const layout = "2006-01-02 15:04:05"

	inputTime, err := time.Parse(layout, inputTimeStr)
	if err != nil {
		fmt.Printf("❌ Gagal parse waktu: %v\n", err)
		return false
	}

	currentTime := time.Now()
	duration := currentTime.Sub(inputTime)

	const twentyFourHours = 24 * time.Hour

	if duration >= twentyFourHours {
		return true // Sudah lebih dari atau sama dengan 24 jam
	}

	return false // Belum mencapai 24 jam
}
