package process

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math/rand/v2"
	"net/http"
	"sync"
	"time"

	"github.com/adityadafe/kc-backend-assgn/internal/models"
	"github.com/adityadafe/kc-backend-assgn/internal/storage"
	"github.com/adityadafe/kc-backend-assgn/internal/utils"
)

func processImage(storeID, imageURL, visitTime string, results chan<- models.Result, db storage.Storage, l *log.Logger) {

	err := db.CheckStore(storeID)

	if err != nil {
		l.Println(err)
		results <- models.Result{
			StoreID: storeID,
			Error:   utils.StoreNotFound,
		}
		return
	}

	resp, err := http.Get(imageURL)
	if err != nil {
		l.Println(err)
		results <- models.Result{
			StoreID: storeID,
			Error:   utils.FailToDownload,
		}
		return
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		l.Println(err)
		results <- models.Result{
			StoreID: storeID,
			Error:   utils.FailToDecode,
		}
		return
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	perimeter := getPerimeter(height, width)

	sleepDuration := time.Duration(rand.Float64()*(0.4-0.1)+0.1) * time.Second
	time.Sleep(sleepDuration)

	l.Println("Successfully parsed image with store id ", storeID)
	results <- models.Result{
		StoreID:   storeID,
		ImageURL:  imageURL,
		Perimeter: perimeter,
	}
}

func ProcessJob(jobId string, job models.JobPayload, db storage.Storage, log *log.Logger) {
	const maxConcurrent = 100
	sem := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup

	totalImages := 0
	for _, visit := range job.Visits {
		totalImages += len(visit.ImageUrls)
	}
	results := make(chan models.Result, totalImages)

	for _, visit := range job.Visits {
		for _, imageURL := range visit.ImageUrls {
			sem <- struct{}{}
			wg.Add(1)

			go func(storeID, url, visitTime string) {
				defer wg.Done()
				defer func() { <-sem }()
				processImage(storeID, url, visitTime, results, db, log)
			}(visit.StoreId, imageURL, visit.VisitTime)
		}
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		if res.Error != "" {
			db.UpdateJob(jobId, res.StoreID, utils.JobFailed, res.Error)
			continue
		}
		db.UpdateJob(jobId, res.StoreID, utils.JobCompleted, "")
	}
}

func getPerimeter(height, width int) int {
	return 2 * (width + height)
}
