package parser

import (
	//image ...
	"c/GoExam/imagesUrlColor/model"
	"image"

	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/cenkalti/dominantcolor"
	"golang.org/x/image/draw"
	"golang.org/x/net/html"
)

//Result ...
type Result struct {
	resultURLColor []*model.URLImage
}

//GetImagesLinks ...
func GetImagesLinks(count int) ([]Result, error) {
	links := generate(2)

	var resultURLColor []st.URLImage
	var result st.URLImage
	for _, url := range links {
		result.URLImg = url
		result.Color = findFromUrl(url)
		resultURLColor = append(resultURLColor, result)
	}
	log.Print(resultURLColor)
	return resultURLColor, nil
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

func findFromURL(pageURL string) (img Image) {
	resp, err := http.Get(pageURL)
	if err != nil {
		log.Print(err)
	}

	defer resp.Body.Close()

	img, _, errr := image.Decode(resp.Body)
	if err != nil {
		log.Print(errr)
	}

	return img
}

func ImgColorProcessor(img Image) string {
	// создаём пустое изображение для записи необходимого размера
	dst := image.NewRGBA(image.Rect(0, 0, 200, 200))
	// изменение размера
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
