package api

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Image struct {
	URL         string
	Name        string
	CharaID     int
	CardID      int
	IsAwakening bool
	IsFramed    bool
}

func GetAllImages(charaID int) ([][]Image, error) {
	ids, err := FindCardIDs(charaID)
	if err != nil {
		return nil, err
	}

	// todo
	fmt.Printf("chara: %d\n", charaID)

	var lists [][]Image
	for _, cardID := range ids {
		url := fmt.Sprintf("https://imas.gamedbs.jp/mlth/chara/show/%d/%d", charaID, cardID)

		fmt.Printf("Getting card: %d...", cardID)
		images, err := getImages(url)
		if err != nil {
			return nil, err
		}
		fmt.Printf("done.\n\n")

		lists = append(lists, images)
	}

	return lists, nil
}

func GetImagesWithID(charaID, cardID int) ([]Image, error) {
	url := fmt.Sprintf("https://imas.gamedbs.jp/mlth/chara/show/%d/%d", charaID, cardID)

	return getImages(url)
}

func getImages(url string) ([]Image, error) {
	time.Sleep(2 * time.Second)

	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}

	urls, err := findImageURLs(doc)
	if err != nil {
		return nil, err
	}

	tokens := strings.Split(url, "/")
	charaID, _ := strconv.Atoi(tokens[len(tokens)-2])
	cardID, _ := strconv.Atoi(tokens[len(tokens)-1])

	var images []Image
	for i, url := range urls {
		images = append(images, Image{
			URL:         url,
			Name:        findCardName(doc),
			CharaID:     charaID,
			CardID:      cardID,
			IsAwakening: i >= 2,
			IsFramed:    i%2 == 0,
		})
	}

	return images, nil
}

func getImageURLsFromPath(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return nil, err
	}

	return findImageURLs(doc)
}

func findImageURLs(doc *goquery.Document) ([]string, error) {
	var urls []string

	selection := doc.Find(`
		body#top > 
		div#container > 
		div#sb-site > 
		div#contents > 
		div#contents-main > 
		section > 
		section.imgbox.flexbox.flexwrap > 
		article.d2_3 > 
		section.imgbox.flexbox.flexwrap >
		article.tc
	`)

	selection = selection.Next()
	selection = selection.Next()

	for i := 0; i < 4; i++ {
		inner := selection.Find("a")
		url, exists := inner.Attr("href")
		if exists {
			urls = append(urls, url)
		} else {
			return nil, errors.New("url not found")
		}

		selection = selection.Next()
	}

	return urls, nil
}

func findCardName(doc *goquery.Document) string {
	selection := doc.Find(`
		body#top > 
		div#container > 
		div#sb-site > 
		div#contents > 
		div#contents-main > 
		section > 
		section.imgbox.flexbox.flexwrap > 
		article.d2_3 > 
		h2
	`)

	name := selection.Text()
	if name == "" {
		name = "default"
	}

	prefix := "カード情報（"
	suffix := "）カード一覧関連画像"

	if !strings.HasPrefix(name, prefix) {
		name = "default"
	} else {
		name = name[len(prefix) : len(name)-len(suffix)]
	}

	return name
}

func FindCardIDs(charaID int) ([]int, error) {
	url := fmt.Sprintf("https://imas.gamedbs.jp/mlth/chara/show/%d", charaID)

	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}

	selection := doc.Find(`
		body#top > 
		div#container > 
		div#sb-site > 
		div#contents > 
		div#contents-main > 
		section > 
		section.imgbox.flexbox.flexwrap > 
		article.d2_3 > 
		ul.dblst.flexbox.flexwrap > 
		li.hvr-grow
	`)

	var ids []int
	for {
		innerSelection := selection.Find("a")

		url, exists := innerSelection.Attr("href")
		if !exists {
			break
		}

		tokens := strings.Split(url, "/")

		id, err := strconv.Atoi(tokens[len(tokens)-1])
		if err != nil {
			return nil, err
		}

		ids = append(ids, id)

		selection = selection.Next()
	}

	return ids, nil
}
