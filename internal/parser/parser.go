package parser

import (
	"c/GoExam/imagesUrlColor/model"
	"image"

	"golang.org/x/image/draw"

	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/cenkalti/dominantcolor"
	"golang.org/x/net/html"

	_ "image/jpeg" //image ...
)

//GetImagesLinks ...
func GetImagesLinks() []model.URLImage {
	links := imgURLParser(2, 100)

	var resURLColor []model.URLImage
	var result model.URLImage
	for _, url := range links {
		result.URLImg = url
		img := findFromURL(url)
		result.Color = imgColorProcessor(img)
		resURLColor = append(resURLColor, result)
	}
	log.Print(resURLColor)

	return resURLColor
}

//imgURLParser ...
func imgURLParser(workers int, count int) []string {
	s := "https://wallpaperstock.net"
	var allURL []string
	ch := make(chan []string)

	for i := 2; len(allURL) < count || i < workers+2; i++ {
		go findLinks(s, ch)
		s = s + "/wallpapers_p" + strconv.Itoa(i) + ".html"
		allURL = append(allURL, <-ch...)
	}

	log.Printf("\nСкачено %v ссылок \n", len(allURL))

	return allURL
}

//findFromURL ...
func findFromURL(pageURL string) image.Image {
	resp, err := http.Get(pageURL)
	if err != nil {
		log.Print(err)
	}

	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		log.Print(err)
	}

	return img
}

func imgColorProcessor(img image.Image) string {
	dst := image.NewRGBA(image.Rect(0, 0, 200, 200))

	draw.CatmullRom.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)

	return dominantcolor.Hex(dominantcolor.Find(dst))
}

func findLinks(url string, c chan []string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Print(err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()

	c <- visit(nil, doc)
}

func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "img" {
		for _, a := range n.Attr {
			if a.Key == "src" && strings.Contains(a.Val, "wallpapers/thumbs") {
				links = append(links, "https:"+a.Val)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}
	return links
}
