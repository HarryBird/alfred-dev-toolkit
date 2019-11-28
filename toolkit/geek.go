package toolkit

import (
	"log"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	alfred "github.com/HarryBird/alfred-toolkit-go"
	"github.com/parnurzeal/gorequest"
	"github.com/urfave/cli"
)

const (
	CateAll     string = ""
	CateArticle string = "article"
	CateColumn  string = "column"
	CateDaily   string = "daily"
)

const (
	ItemArticle string = "article"
	ItemColumn  string = "product"
	ItemColl    string = "collection"
)

const (
	ArticleURL       = "https://time.geekbang.org/column/article/"
	ArticleColumnURL = "https://time.geekbang.org/column/intro/"
	VideoURL         = "https://time.geekbang.org/course/detail/"
	VideoColumnURL   = "https://time.geekbang.org/course/intro/"
	DailyVideoURL    = "https://time.geekbang.org/dailylesson/detail/"
	DailyCollURL     = "https://time.geekbang.org/dailylesson/collection/"
)

const (
	SearchURL = "https://time.geekbang.org/serv/v3/search"
)

type SearchResult struct {
	Code  int              `json:"code"`
	Data  SearchResultData `json:"data"`
	Error interface{}      `json:"error"`
}

type SearchResultData struct {
	List []SearchResultItem `json:"list"`
}

type SearchResultItem struct {
	ItemType   string                 `json:"item_type"`
	Article    SearchResultArticle    `json:"article"`
	Column     SearchResultColumn     `json:"product"`
	Collection SearchResultCollection `json:"collection"`
}

type SearchResultArticle struct {
	Id           int64  `json:"id"`
	Title        string `json:"title"`
	Type         string `json:"product_type"`
	Subtitle     string `json:"subtitle"`
	ColumnTitle  string `json:"product_title"`
	ColumnAuthor string `json:"product_author"`
	Digest       string `json:"content"`
	Sku          int64  `json:"product_id"`
}

type SearchResultColumn struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	Subtitle    string `json:"subtitle"`
	AuthorName  string `json:"author_name"`
	AuthorIntro string `json:"author_intro"`
}

type SearchResultCollection struct {
	Id         int64  `json:"id"`
	Title      string `json:"title"`
	Type       string `json:"product_type"`
	Subtitle   string `json:"subtitle"`
	AuthorName string `json:"author_name"`
	Digest     string `json:"content"`
}

type Result struct {
	Title    string
	Subtitle string
	URL      string
}

func GeekSearchAction(ctx *cli.Context, al *alfred.Alfred) {
	query, cate, err := handleArgs(ctx.Args())

	if query == "" {
		al.ResultAppend(alfred.NewErrorTitleItem("Invalid Query", err.Error())).Output()
		return
	}

	res, err := search(query, cate)

	if len(res) == 0 {
		al.ResultAppend(alfred.NewNoResultItem()).Output()
		return
	}

	log.Println("Items:", res)

	for _, item := range res {
		al.ResultAppend(buildGeekItem(item.Title, item.Subtitle, item.URL))
	}

	al.Output()
}

func buildGeekItem(title, subTitle, arg string) alfred.Item {
	return alfred.NewItem(title, subTitle, arg, arg, "", "", true, alfred.NewIcon("", "./icons/geektime/read.png"))
}

func handleArgs(input cli.Args) (string, string, error) {
	args := []string(input)
	log.Println("Args:", args)

	if len(args) == 0 {
		return "", CateAll, errors.New(sign("Empty Args"))
	}

	query := args[0]
	cate := CateAll

	if len(args) > 1 {
		cate = strings.ToLower(args[1])

		if cate != CateArticle && cate != CateColumn && cate != CateDaily {
			cate = CateAll
		}
	}

	return query, cate, nil
}

func genBaseURL(cate, ty string) string {
	url := ArticleURL

	switch ty {
	case "c1", "c2":
		switch cate {
		case ItemColumn:
			url = ArticleColumnURL // 专栏/微课课程
		case ItemArticle:
			url = ArticleURL // 专栏/微课文章
		}
	case "c3":
		switch cate {
		case ItemColumn:
			url = VideoColumnURL // 视频课程
		case ItemArticle:
			url = VideoURL // 视频
		}
	case "d":
		url = DailyVideoURL // 每日一课视频
	case "c8", "c9":
		url = ArticleURL // 极客视点/二叉树
	}

	if cate == ItemColl {
		url = DailyCollURL
	}

	log.Println("GenBaseURL", cate, ty, url)
	return url
}

func trimHighLight(st string) string {
	st = strings.Replace(st, "<em>", "", -1)
	st = strings.Replace(st, "</em>", "", -1)
	return st
}

func handleTitle(cate, ty, title string) string {
	ret := ""
	switch cate {
	case ItemColumn:
		switch ty {
		case "c1":
			ret = "[专栏] " + title
		case "c2":
			ret = "[微课] " + title
		case "c3":
			ret = "[视频课] " + title
		}
	case ItemArticle:
		switch ty {
		case "c1", "c2", "c8":
			ret = "[文章] " + title
		case "c3":
			ret = "[视频] " + title
		case "d":
			ret = "[每日一课] " + title
		case "c9":
			ret = "[二叉树] " + title
		}
	case ItemColl:
		ret = "[合辑] " + title
	}

	log.Println("HandleTitle", cate, ty, title, ret)

	return ret
}

func handleItem(item SearchResultItem) (string, string, string) {
	t, st, url := "", "", ""

	switch item.ItemType {
	case ItemColumn:
		t = handleTitle(item.ItemType, item.Column.Type, trimHighLight(item.Column.Title))
		st = trimHighLight(item.Column.Subtitle) + " | " + trimHighLight(item.Column.AuthorName)
		url = genBaseURL(item.ItemType, item.Column.Type) + strconv.FormatInt(item.Column.Id, 10)
	case ItemArticle:
		t = handleTitle(item.ItemType, item.Article.Type, trimHighLight(item.Article.Title))
		st = trimHighLight(item.Article.Digest)

		switch item.Article.Type {
		case "d":
			url = genBaseURL(item.ItemType, item.Article.Type) + strconv.FormatInt(item.Article.Sku, 10)
		case "c3":
			url = genBaseURL(item.ItemType, item.Article.Type) + strconv.FormatInt(item.Article.Sku, 10) + "-" + strconv.FormatInt(item.Article.Id, 10)
		default:
			url = genBaseURL(item.ItemType, item.Article.Type) + strconv.FormatInt(item.Article.Id, 10)
		}

	case ItemColl:
		t = handleTitle(item.ItemType, item.Collection.Type, trimHighLight(item.Collection.Title))
		st = trimHighLight(item.Collection.Subtitle) + " | " + trimHighLight(item.Collection.AuthorName)
		url = genBaseURL(item.ItemType, item.Collection.Type) + strconv.FormatInt(item.Collection.Id, 10)
	}

	log.Println("HandleItem", t, st, url)

	return t, st, url
}

func search(query string, cate string) ([]Result, error) {
	var res []Result
	var resp SearchResult

	// change fix
	if cate == CateColumn {
		cate = ItemColumn
	}

	postBody := map[string]interface{}{
		"category": cate,
		"prev":     1,
		"keyword":  query,
		"size":     20,
	}

	log.Println("Query Action:", SearchURL, postBody)

	response, _, errs := gorequest.New().Post(SearchURL).
		Set("User-Agent", SelfName).
		Set("Referer", "https://alfred.geekbang.org").
		Send(postBody).
		EndStruct(&resp)

	log.Println(SearchURL, response, resp, errs)

	if len(errs) > 0 {
		return res, errors.Wrap(errs[0], sign("Search Fail"))
	}

	if response.StatusCode != 200 {
		return res, errors.New(sign("Invalid Response Status: " + strconv.Itoa(response.StatusCode)))
	}

	if resp.Code != 0 {
		return res, errors.New(sign("Invalid Response Code: " + strconv.Itoa(resp.Code)))
	}

	for _, item := range resp.Data.List {
		t, st, url := handleItem(item)
		res = append(res, Result{
			Title:    t,
			Subtitle: st,
			URL:      url,
		})
	}

	return res, nil
}
