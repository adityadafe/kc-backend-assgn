package process

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math/rand/v2"
	"net/http"
	"sync"
	"time"

	"github.com/adityadafe/kc-backend-assgn/internal/models"
	"github.com/adityadafe/kc-backend-assgn/internal/storage"
	"github.com/adityadafe/kc-backend-assgn/internal/utils"
)

// func ProcessJob(jobId string, job models.JobPayload, db storage.Storage) {

// 	totalImages := 0
// 	for _, visit := range job.Visits {
// 		totalImages += len(visit.ImageUrls)
// 	}
// 	results := make(chan models.Result, totalImages)
// 	var wg sync.WaitGroup

// 	for _, visit := range job.Visits {
// 		for _, imageURL := range visit.ImageUrls {
// 			wg.Add(1)
// 			go processImage(visit.StoreId, imageURL, visit.VisitTime, &wg, results)
// 		}
// 	}

// 	wg.Wait()
// 	close(results)

// 	for res := range results {
// 		if res.Error != "" {
// 			fmt.Println("res", res.Error)
// 			db.UpdateJob(jobId, res.StoreID, utils.JobFailed, res.Error)
// 			continue
// 		}
// 		db.UpdateJob(jobId, utils.JobCompleted, "", "")
// 	}

// }

// func processImage(storeID, imageURL, visitTime string, wg *sync.WaitGroup, results chan<- models.Result) {
// 	defer wg.Done()

// 	resp, err := http.Get(imageURL)
// 	if err != nil {
// 		results <- models.Result{
// 			StoreID: storeID,
// 			Error:   utils.FailToDownload,
// 		}
// 		return
// 	}
// 	defer resp.Body.Close()

// 	img, _, err := image.Decode(resp.Body)
// 	if err != nil {
// 		results <- models.Result{
// 			StoreID: storeID,
// 			Error:   utils.FailToDecode,
// 		}
// 		return
// 	}

// 	bounds := img.Bounds()
// 	width := bounds.Dx()
// 	height := bounds.Dy()
// 	perimeter := 2 * (width + height)

// 	//sleepDuration := time.Duration(rand.Float64()*(0.4-0.1)+0.1) * time.Second
// 	time.Sleep(10 * time.Second)

// 	results <- models.Result{
// 		StoreID:   storeID,
// 		ImageURL:  imageURL,
// 		Perimeter: perimeter,
// 	}
// }

func processImage(storeID, imageURL, visitTime string, results chan<- models.Result) {
	resp, err := http.Get(imageURL)
	if err != nil {
		results <- models.Result{
			StoreID: storeID,
			Error:   utils.FailToDownload,
		}
		return
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		results <- models.Result{
			StoreID: storeID,
			Error:   utils.FailToDecode,
		}
		return
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	perimeter := 2 * (width + height)

	sleepDuration := time.Duration(rand.Float64()*(0.4-0.1)+0.1) * time.Second
	time.Sleep(sleepDuration)

	results <- models.Result{
		StoreID:   storeID,
		ImageURL:  imageURL,
		Perimeter: perimeter,
	}
}

func ProcessJob(jobId string, job models.JobPayload, db storage.Storage) {
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
				processImage(storeID, url, visitTime, results)
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
