package fbcolly

import "C"
import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/gocolly/colly/storage"
	"github.com/google/logger"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
	"github.com/thoas/go-funk"
	"net/url"
	"qnetwork.net/fbcrawl/fbcrawl"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Fbcolly struct {
	collector *colly.Collector
	w         *when.Parser
}
type FbDataPostContext struct {
	PublishTime int64 `json:"publish_time"`
}
type FbDataInsight struct {
	FbDataPostContext `json:"post_context"`
}
type FbDataFt struct {
	ContentOwnerIdNew    int64                    `json:"content_owner_id_new"`
	PhotoAttachmentsList []string                 `json:"photo_attachments_list"`
	PhotoId              int64                    `json:"photo_id,string"`
	PageId               int64                    `json:"page_id,string"`
	TopLevelPostId       int64                    `json:"top_level_post_id,string"`
	PageInsights         map[string]FbDataInsight `json:"page_insights"`
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

func setupSharedCollector(collector *colly.Collector) error {
	var err error
	extensions.Referer(collector)
	collector.AllowURLRevisit = true
	collector.OnRequest(sharedOnRequest)
	collector.OnResponse(sharedOnResponse)
	collector.OnError(func(resp *colly.Response, errHttp error) {
		err = errHttp
		logger.Error("OnError", err)
	})
	return err
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

func New() *Fbcolly {
	f := Fbcolly{}
	f.collector = colly.NewCollector()
	f.w = when.New(nil)
	f.w.Add(en.All...)
	f.w.Add(common.All...)
	return &f
}

func (f *Fbcolly) Login(email string, password string, otp string) (string, error) {
	collector := f.collector.Clone()
	err := setupSharedCollector(collector)

	logger.Info("Login using email", email)
	loggedIn := false

	collector.OnHTML("#login_form", func(element *colly.HTMLElement) {
		logger.Info("OnHTML login_form")
		loginURL, err, reqMap := getForm(element, err)
		if err != nil {
			logger.Error(err)
			return
		}
		reqMap["email"] = email
		reqMap["pass"] = password
		logger.Info("req map:", reqMap)
		err = collector.Post(loginURL, reqMap)
		if err != nil {
			logger.Error("post err:", err)
		}
	})

	collector.OnHTML("a[href=\"/login/save-device/cancel/?flow=interstitial_nux&nux_source=regular_login\"]", func(element *colly.HTMLElement) {
		err = collector.Visit("http://mbasic.facebook.com" + element.Attr("href"))
	})

	collector.OnHTML("form[action=\"/login/checkpoint/\"]", func(element *colly.HTMLElement) {

		checkpointUrl, err, reqMap := getForm(element, err)
		shouldSubmit := false
		if err != nil {
			logger.Error(err)
			return
		}

		if element.DOM.Find("input[name=\"name_action_selected\"]").Length() > 0 {
			//Save Device
			logger.Info("OnHTML Save Device checkpoint")
			reqMap["name_action_selected"] = "dont_save"
			shouldSubmit = true
		} else if element.DOM.Find("input[name=\"approvals_code\"]").Length() > 0 {
			logger.Info("OnHTML OTP checkpoint")
			//logger.Info("Please input OTP")
			//reader := bufio.NewReader(os.Stdin)
			//code, _ := reader.ReadString('\n')
			code := otp[0:6]
			reqMap["approvals_code"] = code
			shouldSubmit = true
		} else {
			logger.Info("OnHTML Only Continue checkpoint")

		}
		if shouldSubmit {
			logger.Info("req map:", reqMap)
			err = collector.Post(checkpointUrl, reqMap)
		}
		if err != nil {
			logger.Error("post err:", err)
		}
	})

	collector.OnHTML("form[action=\"/search/\"]", func(element *colly.HTMLElement) {
		//We're in home
		logger.Info("I'm IN HOME, navigate to page now")
		loggedIn = true
	})

	err = collector.Visit("https://mbasic.facebook.com/")
	if err != nil {
		logger.Error("crawl by colly err:", err)
	}

	if loggedIn {
		logger.Info(storage.StringifyCookies(collector.Cookies("https://mbasic.facebook.com/")))
		return storage.StringifyCookies(collector.Cookies("https://mbasic.facebook.com/")), err
	} else {
		return "", err
	}

}

func (f *Fbcolly) FetchGroupFeed(groupId int64) (error, *fbcrawl.FacebookPostList) {
	collector := f.collector.Clone()
	err := setupSharedCollector(collector)
	currentPage := 1
	var result []*fbcrawl.FacebookPost

	collector.OnHTML("#m_group_stories_container > :last-child a", func(element *colly.HTMLElement) {
		currentPage++
		if currentPage < 3 {
			logger.Info("Will fetch page", currentPage)
			err = collector.Visit("http://mbasic.facebook.com" + element.Attr("href"))
		}
	})

	collector.OnXML("//a[text()=\"Full Story\"]", func(element *colly.XMLElement) {
		logger.Info("Post found at", element.Attr("href"))
		u, _ := url.Parse("http://mbasic.facebook.com" + element.Attr("href"))
		postId, _ := strconv.ParseInt(u.Query().Get("id"), 10, 64)
		result = append(result, &fbcrawl.FacebookPost{
			Id:    postId,
			Group: &fbcrawl.FacebookGroup{Id: groupId},
		})
		//f.detailCollector.Visit(url)
	})

	err = collector.Visit(fmt.Sprintf("https://mbasic.facebook.com/groups/%d", groupId))
	if err != nil {
		logger.Error("crawl by colly err:", err)
	}
	return err, &fbcrawl.FacebookPostList{Posts: result}
}

func (f *Fbcolly) FetchContentImages(postId int64) (error, *fbcrawl.FacebookImageList) {
	collector := f.collector.Clone()
	err := setupSharedCollector(collector)
	currentPage := 1
	var result []*fbcrawl.FacebookImage

	collector.OnHTML("a[href*=\"/media/set/\"]", func(element *colly.HTMLElement) {
		currentPage++
		logger.Info("Will fetch page", currentPage)
		err = collector.Visit("http://mbasic.facebook.com" + element.Attr("href"))
	})

	collector.OnHTML("a[href*=\"/photo.php\"]", func(element *colly.HTMLElement) {
		result = append(result, &fbcrawl.FacebookImage{
			Id: getImageIdFromHref(element.Attr("href")),
		})
		//f.detailCollector.Visit(url)
	})

	err = collector.Visit(fmt.Sprintf("https://mbasic.facebook.com/media/set/?set=pcb.%d", postId))
	if err != nil {
		logger.Error("crawl by colly err:", err)
	}
	return err, &fbcrawl.FacebookImageList{Images: result}
}

func (f *Fbcolly) FetchImageUrl(imageId int64) (error, *fbcrawl.FacebookImage) {
	collector := f.collector.Clone()
	err := setupSharedCollector(collector)
	result := fbcrawl.FacebookImage{Id: imageId}

	collector.OnHTML("a", func(element *colly.HTMLElement) {
		result.Url = element.Attr("href")
	})

	err = collector.Visit(fmt.Sprintf("https://mbasic.facebook.com/photo/view_full_size/?fbid=%d", imageId))
	if err != nil {
		logger.Error("crawl by colly err:", err)
	}
	return err, &result
}

func (f *Fbcolly) FetchPost(groupId int64, postId int64) (error, *fbcrawl.FacebookPost) {
	collector := f.collector.Clone()
	err := setupSharedCollector(collector)
	post := &fbcrawl.FacebookPost{Comments: []*fbcrawl.FacebookComment{}}
	commentPaging := 0
	collector.OnHTML("#m_story_permalink_view", func(element *colly.HTMLElement) {
		dataElement := element.DOM.Find("div[data-ft]")
		if dataElement.Length() > 0 {
			var result FbDataFt
			jsonData, isExist := dataElement.Attr("data-ft")
			if isExist {
				logger.Info(jsonData)
				err = json.Unmarshal([]byte(jsonData), &result)
				if err != nil {
					logger.Error(err)
					return
				}
				logger.Info("Post ", result)
				post.Id = result.TopLevelPostId
				post.Group = &fbcrawl.FacebookGroup{Id: result.PageId, Name: dataElement.Find("h3 strong:last-child a").Text()}
				post.User = &fbcrawl.FacebookUser{
					Id:   result.ContentOwnerIdNew,
					Name: dataElement.Find("h3 strong:first-child a").Text(),
				}
				post.CreatedAt = result.PageInsights[strconv.FormatInt(result.PageId, 10)].PublishTime
				//Content

				//NO BACKGROUND TEXT ONLY
				post.Content = strings.Join(dataElement.Find("p").Map(func(i int, selection *goquery.Selection) string {
					return selection.Text()
				}), "\n")

				if len(post.Content) == 0 {
					// TEXT WITH BACKGROUND
					post.Content = dataElement.Find("div[style*=\"background-image:url\"]").Text()
				}

				post.ContentLink = getUrlFromRedirectHref(dataElement.Find("a[href*=\"https://lm.facebook.com/l.php\"]").AttrOr("href", ""))
				post.ReactionCount = getReactionFromText(element.DOM.Find("div[id*=\"sentence_\"]").Text())
				post.ContentImages = (funk.Map(result.PhotoAttachmentsList, func(id string) *fbcrawl.FacebookImage {
					i, _ := strconv.ParseInt(id, 10, 64)
					return &fbcrawl.FacebookImage{
						Id: i,
					}
				})).([]*fbcrawl.FacebookImage)

				if result.PhotoId > 0 {
					post.ContentImage = &fbcrawl.FacebookImage{Id: result.PhotoId}
				}

				logger.Info("content", strings.Join(dataElement.Find("p").Map(func(i int, selection *goquery.Selection) string {
					return selection.Text()
				}), "\n"))
			}

			//Comment
			element.DOM.Find("h3 + div + div + div").Parent().Parent().Each(func(i int, selection *goquery.Selection) {
				//author
				commentId, _ := strconv.ParseInt(selection.AttrOr("id", ""), 10, 64)
				logger.Info("comment", commentId)
				createdAtWhenResult, _ := f.w.Parse(selection.Find("abbr").Text(), time.Now())
				post.Comments = append(post.Comments, &fbcrawl.FacebookComment{
					Id:   commentId,
					Post: &fbcrawl.FacebookPost{Id: post.Id},
					User: &fbcrawl.FacebookUser{
						Id:   getUserIdFromCommentHref(selection.Find("a[href*=\"#comment_form_\"]").AttrOr("href", "")),
						Name: selection.Find("h3 > a").Text(),
					},
					Content:   selection.Find("h3 + div").Text(),
					CreatedAt: createdAtWhenResult.Time.Unix(),
				})
			})

		}
	})

	collector.OnHTML("div[id*=\"see_prev_\"] > a", func(element *colly.HTMLElement) {
		if commentPaging < 3 {
			logger.Info("Comment paging", commentPaging)
			err = collector.Visit("http://mbasic.facebook.com" + element.Attr("href"))
			commentPaging = commentPaging + 1
		}
	})

	err = collector.Visit(fmt.Sprintf("http://mbasic.facebook.com/groups/%d?view=permalink&id=%d&_rdr", groupId, postId))
	return err, post
}

func (f *Fbcolly) LoginWithCookies(cookies string) error {
	collector := f.collector
	return collector.SetCookies("https://mbasic.facebook.com/", storage.UnstringifyCookies(cookies))
}

//func getUsernameFromHref(href string) string {
//	return regexp.MustCompile("/([\\d\\w.]+).*").FindStringSubmatch(href)[1]
//}

func getUserIdFromCommentHref(href string) int64 {
	id, _ := strconv.ParseInt(regexp.MustCompile("#comment_form_(\\d+)").FindStringSubmatch(href)[1], 10, 64)
	return id
}

func getUrlFromRedirectHref(href string) string {
	u, _ := url.Parse(href)
	return u.Query().Get("u")
}

func getImageIdFromHref(href string) int64 {
	u, _ := url.Parse(href)
	i, _ := strconv.ParseInt(u.Query().Get("fbid"), 10, 64)
	return i
}

func getReactionFromText(text string) int64 {
	logger.Error("reaction", text)
	if len(text) > 0 {
		id, _ := strconv.ParseInt(regexp.MustCompile("(\\d+)").FindStringSubmatch(text)[1], 10, 64)
		return id
	}
	return 0
}
