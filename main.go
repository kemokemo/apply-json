package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

func main() {
	os.Exit(run())
}

func run() int {
	b, err := ioutil.ReadFile("source.json")
	if err != nil {
		fmt.Println("failed to read source.html,", err)
		return 1
	}

	var da InputData
	err = json.Unmarshal(b, &da)
	if err != nil {
		fmt.Println("failed to unmarshal data,", err)
		return 1
	}

	f, err := os.Open("source.html")
	if err != nil {
		fmt.Println("failed to open file,", err)
		return 1
	}
	defer func() {
		e := f.Close()
		if e != nil {
			fmt.Println("failed to close file, ", e)
		}
	}()

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		fmt.Println("failed to goquery, ", err)
		return 1
	}

	for _, d := range da.DataArray {
		doc.Find(fmt.Sprintf("#%s", d.ID)).Each(func(index int, s *goquery.Selection) {
			if len(s.Nodes) == 0 {
				return
			}
			nodeKind := s.Nodes[0].Data

			switch nodeKind {
			case "input":
				t, exists := s.Attr("type")
				if !exists {
					return
				}

				switch t {
				case "checkbox":
					s.SetAttr("checked", d.Value)
				case "text":
					s.SetAttr("value", d.Value)
				}
			case "select":
				s.Find("option").Each(func(i int, cs *goquery.Selection) {
					v, _ := cs.Attr("value")
					if v == d.Value {
						cs.SetAttr("selected", "")
					}
				})
			default:
				s.SetText(d.Value)
			}
		})
	}

	err = html.Render(os.Stdout, doc.Nodes[0])
	if err != nil {
		fmt.Println("failed to render, ", err)
		return 1
	}

	return 0
}
