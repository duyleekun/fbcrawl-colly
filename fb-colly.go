package main

/*
typedef struct FacebookGroup {
    char *name;
} FacebookGroup;
*/
import "C"
import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/google/logger"
	"os"
	"strings"
)

type CFacebookGroup struct {
	FacebookGroup C.FacebookGroup
}

type FacebookGroup struct {

}

type FacebookPost struct {

}

type FacebookComment struct {

}

type FacebookAuthor struct {

}

func sharedOnRequest(request *colly.Request) {
	logger.Info("OnRequest")
	//request.Headers.Set("Host", "facebook.com")
	request.Headers.Set("Accept-Language", "en-US,en;q=0.9")
	request.Headers.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	request.Headers.Set("origin", "https://mbasic.facebook.com")

	//logger.Info("Saved referrer is", request.Ctx.Get("_referer"))
	request.Headers.Set("referer", "https://mbasic.facebook.com/checkpoint/?_rdr")
	request.Headers.Set("cache-control", "max-age=0")
	request.Headers.Set("upgrade-insecure-requests", "1")
	//accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
	//origin: https://mbasic.facebook.com
	//referer: https://mbasic.facebook.com/checkpoint/?_rdr
	request.Headers.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0.1; Moto G (4)) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Mobile Safari/537.36")
	request.ResponseCharacterEncoding = "utf-8"
}
func setupGroupCollector() error {
	err, collector := setupSharedCollector(groupCollector)
	currentPage := 1

	collector.OnHTML("#m_group_stories_container > :last-child a", func(element *colly.HTMLElement) {
		currentPage++
		if currentPage < 3 {
			logger.Info("Will fetch page", currentPage)
			collector.Visit("http://mbasic.facebook.com" + element.Attr("href"))
		}
	})

	collector.OnXML("//a[text()=\"Full Story\"]", func(element *colly.XMLElement) {
		url := "http://mbasic.facebook.com" + element.Attr("href")
		logger.Info("Post url found ", url)
		detailCollector.Visit(url)
	})

	return err
}

func setupGroupPostCollector() error {
	err, collector := setupSharedCollector(detailCollector)

	collector.OnHTML("#m_story_permalink_view", func(element *colly.HTMLElement) {
		dataElement := element.DOM.Find("div[data-ft]")
		if dataElement.Length() > 0 {
			var result map[string]interface{}
			jsonData, isExist := dataElement.Attr("data-ft")
			if isExist {
				json.Unmarshal([]byte(jsonData), &result)
				logger.Info("Post ", result)
				//Content
				logger.Info(strings.Join(dataElement.Find("p").Map(func(i int, selection *goquery.Selection) string {
					return selection.Text()
				}),"\n"))

			}

			//Comment
			element.DOM.Find("h3 + div + div + div").Parent().Each(func(i int, selection *goquery.Selection) {
				//author
				logger.Info("comment")
				logger.Info(selection.Find("h3 > a").Text())
				logger.Info(selection.Find("h3 + div").Text())
			})
		}
	})

	return err
}

func setupAuthCollector() error {
	err, collector := setupSharedCollector(authCollector)

	collector.OnHTML("#login_form", func(element *colly.HTMLElement) {
		logger.Info("OnHTML login_form")
		loginURL, err, reqMap := getForm(element, err)
		if err != nil {
			logger.Error(err)
			return
		}
		reqMap["email"] = *email
		reqMap["pass"] = *password
		logger.Info("req map:", reqMap)
		err = collector.Post(loginURL, reqMap)
		if err != nil {
			logger.Error("post err:", err)
		}
	})

	collector.OnHTML("form[action=\"/login/checkpoint/\"]", func(element *colly.HTMLElement) {

		checkpointUrl, err, reqMap := getForm(element, err)
		if err != nil {
			logger.Error(err)
			return
		}

		if element.DOM.Find("input[name=\"name_action_selected\"]").Length() > 0 {
			//Save Device
			logger.Info("OnHTML Save Device checkpoint")
			reqMap["name_action_selected"] = "dont_save"
		} else if element.DOM.Find("input[name=\"approvals_code\"]").Length() > 0 {
			logger.Info("OnHTML OTP checkpoint")
			logger.Info("Please input OTP")
			reader := bufio.NewReader(os.Stdin)
			code, _ := reader.ReadString('\n')
			code = code[0:6]
			reqMap["approvals_code"] = code
		} else {
			logger.Info("OnHTML Only Continue checkpoint")
		}
		logger.Info("req map:", reqMap)
		err = collector.Post(checkpointUrl, reqMap)
		if err != nil {
			logger.Error("post err:", err)
		}
	})

	collector.OnHTML("form[action=\"/search/\"]", func(element *colly.HTMLElement) {
		//We're in home
		logger.Info("I'm IN HOME, navigate to page now")
	})

	return err
}

func setupSharedCollector(collector *colly.Collector) (error, *colly.Collector) {
	var err error
	extensions.Referer(collector)

	collector.OnRequest(sharedOnRequest)
	collector.OnResponse(sharedOnResponse)
	collector.OnError(func(resp *colly.Response, errHttp error) {
		err = errHttp
		logger.Error("OnError", err)
	})
	return err, collector
}

func sharedOnResponse(response *colly.Response) {
	logger.Info("OnResponse ./last.html")
	_ = response.Save("./last.html")
	//logger.Info(string(resp.Body))
}

func getForm(element *colly.HTMLElement, err error) (string, error, map[string]string) {
	submitUrl, exists := element.DOM.Attr("action")
	if !exists {
		err = errors.New("doesn't have action label")
		return "", nil, nil
	}
	submitUrl = fmt.Sprintf("https://mbasic.facebook.com%s", submitUrl)
	logger.Info("form url is:", submitUrl)
	reqMap := make(map[string]string)
	element.DOM.Find("input").Each(func(i int, s *goquery.Selection) {
		name, _ := s.Attr("name")
		value, _ := s.Attr("value")
		if name != "" && name != "sign_up" && name != "submit[logout-button-with-confirm]" {
			reqMap[name] = value
		}
	})
	return submitUrl, err, reqMap
}

const logPath = "parse.log"

var verbose = flag.Bool("verbose", true, "print info level logs to stdout")
var email = flag.String("email", "change_me@gmail.com", "facebook email")
var password = flag.String("password", "change_me", "facebook password")
var groupId = flag.String("groupId", "334294967318328", "facebook group id, default is 334294967318328")
var authCollector *colly.Collector
var groupCollector *colly.Collector
var detailCollector *colly.Collector

func main() {
	group := CFacebookGroup{FacebookGroup: C.FacebookGroup{name: C.CString("LOL")}}
	fmt.Printf("%+v", group)
	return

	flag.Parse()

	lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
	}
	defer lf.Close()
	defer logger.Init("fb-colly", *verbose, false, lf).Close()

	collector := colly.NewCollector()
	collector.SetProxy("socks5://localhost:8889")


	authCollector = collector.Clone()
	groupCollector = collector.Clone()
	detailCollector = collector.Clone()
	setupAuthCollector()
	setupGroupCollector()
	setupGroupPostCollector()


	err = authCollector.Visit("https://mbasic.facebook.com/")
	if err != nil {
		logger.Error("crawl by colly err:", err)
	}

	err = groupCollector.Visit("https://mbasic.facebook.com/groups/" + *groupId)
	if err != nil {
		logger.Error("crawl by colly err:", err)
	}

}
