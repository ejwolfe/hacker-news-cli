package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Item struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

func main() {
	storyIds := getData[[]int]("https://hacker-news.firebaseio.com/v0/topstories.json?print=pretty")
	var items []Item
	var wg sync.WaitGroup
	for _, id := range storyIds {
		wg.Add(1)
		go func(id int) {
			item := getData[Item](
				"https://hacker-news.firebaseio.com/v0/item/" + strconv.Itoa(
					id,
				) + ".json?print=pretty",
			)
			items = append(items, item)
			wg.Done()
		}(id)
	}
	wg.Wait()
	for _, item := range items {
		fmt.Println("")
		fmt.Println(item.Title)
		fmt.Println(item.Url)
		fmt.Println("")
	}
}

func handleError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func getData[T any](url string) T {
	resp, err := http.Get(url)
	handleError(err)
	body, err := io.ReadAll(resp.Body)
	handleError(err)
	return parseBody[T](body)
}

func parseBody[T any](body []byte) T {
	var result T
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatalln(err)
	}
	return result
}
