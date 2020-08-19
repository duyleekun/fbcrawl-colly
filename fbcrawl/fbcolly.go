package fbcolly

import "C"
import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/gocolly/colly/v2/storage"
	"github.com/google/logger"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
	"github.com/thoas/go-funk"
	"github.com/xlzd/gotp"
	"net/url"
	"qnetwork.net/fbcrawl/fbcrawl/pb"
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
	ContentOwnerIdNew    json.Number              `json:"content_owner_id_new"`
	PhotoAttachmentsList []string                 `json:"photo_attachments_list"`
	PhotoId              int64                    `json:"photo_id,string"`
	PageId               int64                    `json:"page_id,string"`
	TopLevelPostId       int64                    `json:"top_level_post_id,string"`
	PageInsights         map[string]FbDataInsight `json:"page_insights"`
}

func setupSharedCollector(collector *colly.Collector) error {
	var err error
	extensions.Referer(collector)
	collector.AllowURLRevisit = true
	var lastUrl string
	collector.OnRequest(func(request *colly.Request) {
		lastUrl = request.URL.RawPath
		logger.Info("OnRequest ", request.URL)
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
	})
	collector.OnResponse(func(response *colly.Response) {
		logger.Info("OnResponse ./last.html")
		_ = response.Save("./last.html")
		//logger.Info(string(resp.Body))
	})

	collector.OnHTML("a[href*=\"177066345680802\"", func(element *colly.HTMLElement) {
		logger.Error("RateLimit reached ", lastUrl)
	})
	collector.OnError(func(resp *colly.Response, errHttp error) {
		err = errHttp
		logger.Error("OnError", err)
	})
	return err
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

func (f *Fbcolly) Login(email string, password string, totpSecret string) (string, error) {
	collector := f.collector.Clone()
	err := setupSharedCollector(collector)

	logger.Info("Login using email", email)
	loggedIn := false
	firstLogin := true
	collector.OnHTML("#login_form", func(element *colly.HTMLElement) {
		if firstLogin {
			firstLogin = false
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
			if len(totpSecret) > 0 {
				code := gotp.NewDefaultTOTP(totpSecret).Now()
				reqMap["approvals_code"] = code
				shouldSubmit = true
			}

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

func (f *Fbcolly) FetchGroupFeed(groupId int64, nextCursor string) (error, *pb.FacebookPostList) {
	collector := f.collector.Clone()
	err := setupSharedCollector(collector)
	result := pb.FacebookPostList{Posts: []*pb.FacebookPost{}}

	collector.OnHTML("#m_group_stories_container > :last-child a", func(element *colly.HTMLElement) {
		result.NextCursor = "http://mbasic.facebook.com" + element.Attr("href")
	})
	collector.OnHTML("#m_group_stories_container div[role=\"article\"]", func(element *colly.HTMLElement) {
		dataElement := element
		post := &pb.FacebookPost{}
		var fbDataFt FbDataFt
		jsonData := dataElement.Attr("data-ft")

		logger.Info(jsonData)
		err = json.Unmarshal([]byte(jsonData), &fbDataFt)
		if err != nil {
			logger.Error(err)
			return
		}
		logger.Info("Post ", fbDataFt)
		post.Id = fbDataFt.TopLevelPostId
		post.Group = &pb.FacebookGroup{Id: fbDataFt.PageId, Name: dataElement.DOM.Find("h3 strong:nth-child(2) a").Text()}
		userId, _ := fbDataFt.ContentOwnerIdNew.Int64()
		post.User = &pb.FacebookUser{
			Id:   userId,
			Name: dataElement.DOM.Find("h3 strong:nth-child(1) a").Text(),
		}
		post.CreatedAt = fbDataFt.PageInsights[strconv.FormatInt(fbDataFt.PageId, 10)].PublishTime
		//Content

		//NO BACKGROUND TEXT ONLY
		post.Content = strings.Join(dataElement.DOM.Find("p").Map(func(i int, selection *goquery.Selection) string {
			return selection.Text()
		}), "\n")

		if len(post.Content) == 0 {
			// TEXT WITH BACKGROUND
			post.Content = dataElement.DOM.Find("div[style*=\"background-image:url\"]").Text()
		}

		post.ContentLink = getUrlFromRedirectHref(dataElement.DOM.Find("a[href*=\"https://lm.facebook.com/l.php\"]").AttrOr("href", ""))
		post.ReactionCount = getNumberFromText(element.DOM.Find("span[id*=\"like_\"]").Text())
		post.CommentCount = getNumberFromText(element.DOM.Find("span[id*=\"like_\"] ~ a").Text())
		post.ContentImages = (funk.Map(fbDataFt.PhotoAttachmentsList, func(id string) *pb.FacebookImage {
			i, _ := strconv.ParseInt(id, 10, 64)
			return &pb.FacebookImage{
				Id: i,
			}
		})).([]*pb.FacebookImage)

		if fbDataFt.PhotoId > 0 {
			post.ContentImage = &pb.FacebookImage{Id: fbDataFt.PhotoId}
		}
		result.Posts = append(result.Posts, post)
	})
	if len(nextCursor) > 0 {
		err = collector.Visit(nextCursor)
	} else {
		err = collector.Visit(fmt.Sprintf("https://mbasic.facebook.com/groups/%d", groupId))
	}

	if err != nil {
		logger.Error("crawl by colly err:", err)
	}
	return err, &result
}

func (f *Fbcolly) FetchUserInfo(userIdOrUsername string) (error, *pb.FacebookUser) {
	collector := f.collector.Clone()
	err := setupSharedCollector(collector)

	result := &pb.FacebookUser{}

	collector.OnHTML("a[href*=\"lst=\"]", func(element *colly.HTMLElement) {
		parsed, _ := url.Parse(element.Attr("href"))
		result.Username = strings.Split(parsed.Path[1:], "/")[0]
		result.Id = getUserIdFromCommentHref(parsed.Query().Get("lst"))
	})

	collector.OnHTML("a[href*=\"/friends\"]", func(element *colly.HTMLElement) {
		result.FriendCount = getNumberFromText(element.Text)
	})

	collector.OnHTML("#objects_container", func(element *colly.HTMLElement) {
		result.Name = element.DOM.Find("strong").First().Text()
	})

	err = collector.Visit(fmt.Sprintf("https://mbasic.facebook.com/%s", userIdOrUsername))
	if err != nil {
		logger.Error("crawl by colly err:", err)
	}
	return err, result
}

func (f *Fbcolly) FetchGroupInfo(groupIdOrUsername string) (error, *pb.FacebookGroup) {
	collector := f.collector.Clone()
	err := setupSharedCollector(collector)
	result := &pb.FacebookGroup{}

	collector.OnHTML("a[href=\"#groupMenuBottom\"] h1", func(element *colly.HTMLElement) {
		result.Name = element.Text
	})
	collector.OnHTML("a[href*=\"view=member\"]", func(element *colly.HTMLElement) {
		result.Id = getNumberFromText(element.Attr("href"))
		result.MemberCount, _ = strconv.ParseInt(element.DOM.Closest("tr").Find("td:last-child").Text(), 10, 64)
	})

	err = collector.Visit(fmt.Sprintf("https://mbasic.facebook.com/groups/%s?view=info", groupIdOrUsername))
	if err != nil {
		logger.Error("crawl by colly err:", err)
	}
	return err, result
}

func (f *Fbcolly) FetchContentImages(postId int64, nextCursor string) (error, *pb.FacebookImageList) {
	collector := f.collector.Clone()
	err := setupSharedCollector(collector)
	result := pb.FacebookImageList{Images: []*pb.FacebookImage{}}

	collector.OnHTML("a[href*=\"/media/set/\"]", func(element *colly.HTMLElement) {
		result.NextCursor = "http://mbasic.facebook.com" + element.Attr("href")
	})

	collector.OnHTML("a[href*=\"/photo.php\"]", func(element *colly.HTMLElement) {
		result.Images = append(result.Images, &pb.FacebookImage{
			Id: getImageIdFromHref(element.Attr("href")),
		})
		//f.detailCollector.Visit(url)
	})
	if len(nextCursor) > 0 {
		err = collector.Visit(nextCursor)
	} else {
		err = collector.Visit(fmt.Sprintf("https://mbasic.facebook.com/media/set/?set=pcb.%d", postId))
	}

	if err != nil {
		logger.Error("crawl by colly err:", err)
	}
	return err, &result
}

func (f *Fbcolly) FetchImageUrl(imageId int64) (error, *pb.FacebookImage) {
	collector := f.collector.Clone()
	err := setupSharedCollector(collector)
	result := pb.FacebookImage{Id: imageId}

	collector.OnHTML("a[href*=\"fbcdn\"]", func(element *colly.HTMLElement) {
		result.Url = element.Attr("href")
	})

	err = collector.Visit(fmt.Sprintf("https://mbasic.facebook.com/photo/view_full_size/?fbid=%d", imageId))
	if err != nil {
		logger.Error("crawl by colly err:", err)
	}
	return err, &result
}

func (f *Fbcolly) FetchPost(groupId int64, postId int64, commentNextCursor string) (error, *pb.FacebookPost) {
	collector := f.collector.Clone()
	err := setupSharedCollector(collector)
	post := &pb.FacebookPost{Comments: &pb.CommentList{Comments: []*pb.FacebookComment{}}}
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
				post.Group = &pb.FacebookGroup{Id: result.PageId, Name: dataElement.Find("h3 strong:last-child a").Text()}
				userId, _ := result.ContentOwnerIdNew.Int64()
				post.User = &pb.FacebookUser{
					Id:   userId,
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
				post.ReactionCount = getNumberFromText(element.DOM.Find("div[id*=\"sentence_\"]").Text())
				post.ContentImages = (funk.Map(result.PhotoAttachmentsList, func(id string) *pb.FacebookImage {
					i, _ := strconv.ParseInt(id, 10, 64)
					return &pb.FacebookImage{
						Id: i,
					}
				})).([]*pb.FacebookImage)

				if result.PhotoId > 0 {
					post.ContentImage = &pb.FacebookImage{Id: result.PhotoId}
				}
			}

			//Comment
			element.DOM.Find("h3 + div + div + div").Parent().Parent().Each(func(i int, selection *goquery.Selection) {
				//author
				commentId, _ := strconv.ParseInt(selection.AttrOr("id", ""), 10, 64)
				if commentId > 0 {
					createdAtWhenResult, err := f.w.Parse(selection.Find("abbr").Text(), time.Now())
					if err != nil {
						logger.Error(err)
						return
					}
					parsed, err := url.Parse(selection.Find("h3 > a").AttrOr("href", ""))
					if err != nil {
						logger.Error(err)
						return
					}
					if len(parsed.Path) == 0 {
						logger.Error("Empty path for commentId ", commentId)
						return
					}
					if len(parsed.Path) > 1 {
						post.Comments.Comments = append(post.Comments.Comments, &pb.FacebookComment{
							Id:   commentId,
							Post: &pb.FacebookPost{Id: post.Id},
							User: &pb.FacebookUser{
								Username: parsed.Path[1:],
								Name:     selection.Find("h3 > a").Text(),
							},
							Content:   selection.Find("h3 + div").Text(),
							CreatedAt: createdAtWhenResult.Time.Unix(),
						})
					}
				}
			})

		}
	})

	collector.OnHTML("div[id*=\"see_prev_\"] > a", func(element *colly.HTMLElement) {
		post.Comments.NextCursor = "http://mbasic.facebook.com" + element.Attr("href")
	})
	if len(commentNextCursor) > 0 {
		err = collector.Visit(commentNextCursor)
	} else {
		err = collector.Visit(fmt.Sprintf("http://mbasic.facebook.com/groups/%d?view=permalink&id=%d&_rdr", groupId, postId))
	}

	return err, post
}

func (f *Fbcolly) LoginWithCookies(cookies string) error {
	collector := f.collector
	return collector.SetCookies("https://mbasic.facebook.com/", storage.UnstringifyCookies(cookies))
}

func (f *Fbcolly) FetchMyGroups() (error, *pb.FacebookGroupList) {
	collector := f.collector.Clone()
	err := setupSharedCollector(collector)
	result := &pb.FacebookGroupList{Groups: []*pb.FacebookGroup{}}

	collector.OnHTML("li table a", func(element *colly.HTMLElement) {
		result.Groups = append(result.Groups, &pb.FacebookGroup{
			Id:   getNumberFromText(element.Attr("href")),
			Name: element.Text,
		})
	})

	err = collector.Visit("https://mbasic.facebook.com/groups/?seemore")
	if err != nil {
		logger.Error("crawl by colly err:", err)
	}
	return err, result
}

//func getUsernameFromHref(href string) string {
//	return regexp.MustCompile("/([\\d\\w.]+).*").FindStringSubmatch(href)[1]
//}

func getUserIdFromCommentHref(href string) int64 {
	match := regexp.MustCompile("\\d+:(\\d+)\\d+").FindStringSubmatch(href)
	if len(match) > 0 {
		id, _ := strconv.ParseInt(match[1], 10, 64)
		return id
	}
	return 0
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

func getNumberFromText(text string) int64 {
	logger.Info("getNumberFromText ", text)
	if len(text) > 0 {
		match := regexp.MustCompile("(\\d+)\\s?([km]?)").FindStringSubmatch(text)
		if len(match) > 0 {
			count, _ := strconv.ParseInt(match[1], 10, 64)
			switch match[2] {
			case "k":
				count *= 1000
			case "m":
				count *= 1000000
			}
			return count
		}
	}
	return 0
}
