package main

import (
	"log"
	"os"
	"strconv"

	"github.com/earlgray283/imasdb-cli/api"
)

func main() {
	charaID := 44
	imageimages, err := api.GetAllImages(charaID)
	if err != nil {
		log.Fatal(err)
	}

	os.MkdirAll(strconv.Itoa(charaID), 0755)
	for _, images := range imageimages {
		for _, image := range images {
			name := image.Name + func(isAwakening, isFramed bool) string {
				suffix := ""
				if isAwakening {
					suffix += "-覚醒後"
				} else {
					suffix += "-覚醒前"
				}
				if isFramed {
					suffix += "(枠有)"
				}

				return suffix + ".png"
			}(image.IsAwakening, image.IsFramed)

			if err := DownloadFile(name, strconv.Itoa(charaID), image.URL); err != nil {
				log.Fatal(err)
			}
		}
	}
}
