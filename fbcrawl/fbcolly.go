package fbcolly

import "C"
import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/gocolly/colly/v2/storage"
	"github.com/google/logger"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
	"github.com/pkg/errors"
	"github.com/xlzd/gotp"
	"net"
	"net/http"
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

type FbDataFt struct {
	ContentOwnerIdNew json.Number `json:"content_owner_id_new"`
	PhotoId           int64       `json:"photo_id,string"`
	PageId            int64       `json:"page_id,string"`
	TopLevelPostId    int64       `json:"top_level_post_id,string"`
}

func setupSharedCollector(collector *colly.Collector, onError func(error)) {

	extensions.Referer(collector)
	collector.AllowURLRevisit = true
	var lastUrl string

	collector.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})

	collector.OnRequest(func(request *colly.Request) {
		lastUrl = request.URL.RawPath
		logger.Info("OnRequest ", request.URL)

		request.Headers.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		request.Headers.Set("accept-language", "vi,en-US;q=0.9,en;q=0.8,zh-CN;q=0.7,zh-TW;q=0.6,zh;q=0.5")
		request.Headers.Set("cache-control", "no-cache")
		request.Headers.Set("dnt", "1")
		request.Headers.Set("pragma", "no-cache")
		request.Headers.Set("sec-fetch-dest", "document")
		request.Headers.Set("sec-fetch-mode", "navigate")
		request.Headers.Set("sec-fetch-site", "none")
		request.Headers.Set("sec-fetch-user", "?1")
		request.Headers.Set("upgrade-insecure-requests", "1")
		request.Headers.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.109 Safari/537.36")

		request.ResponseCharacterEncoding = "utf-8"
	})
	collector.OnResponse(func(response *colly.Response) {
		name := strconv.FormatInt(time.Now().Unix(), 10) + ".html"
		logger.Info("OnResponse ./" + name)
		_ = response.Save("./" + name)
		//logger.Info(string(response.Body))
	})

	collector.OnHTML("a[href*=\"177066345680802\"]", func(element *colly.HTMLElement) {
		logger.Error("RateLimit reached ")
		onError(errors.New("RateLimit reached"))
	})

	// Set error handler
	collector.OnError(func(r *colly.Response, err error) {
		logger.Error("Request URL:", r.Request.URL, " failed with response:", r.StatusCode, r.Headers, "\nError:", err)
		onError(err)

	})
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

func FacebookRule() rules.Rule {

	return &rules.F{
		RegExp: regexp.MustCompile("(20\\d{2})"),
		Applier: func(m *rules.Match, c *rules.Context, o *rules.Options, ref time.Time) (bool, error) {

			year, err := strconv.Atoi(m.Captures[0])
			if err != nil {
				return false, errors.Wrap(err, "year rule")
			}
			c.Year = &year
			return true, nil
		},
	}
}

func New() *Fbcolly {
	f := Fbcolly{}
	f.collector = colly.NewCollector()
	f.w = when.New(nil)
	f.w.Add(en.All...)
	f.w.Add(common.All...)
	f.w.Add(FacebookRule())
	return &f
}

func (f *Fbcolly) Login(email string, password string, totpSecret string) (*pb.LoginResponse, error) {

	collector := f.collector.Clone()
	var err error
	setupSharedCollector(collector, func(errInner error) {
		err = errInner
	})

	logger.Info("Login using email", email)
	loggedIn := false
	firstLogin := true
	collector.OnHTML("form[action*=\"/basic-lite/cookie/consent\"]", func(element *colly.HTMLElement) {
		logger.Info("OnHTML consent form")
		loginURL, err, reqMap := getForm(element, err)
		if err != nil {
			logger.Error(err)
			return
		}
		err = collector.Post(loginURL, reqMap)
		if err != nil {
			logger.Error("post err:", err)
		}
	})
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
		_ = collector.Visit("https://mbasic.facebook.com" + element.Attr("href"))
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
		err = nil
	})

	_ = collector.Visit("https://mbasic.facebook.com/")

	var cookies string
	if loggedIn {
		cookies = storage.StringifyCookies(collector.Cookies("https://mbasic.facebook.com/"))
	} else {
		cookies = ""
		err = errors.New(fmt.Sprintf("CAN'T LOGIN"))

	}
	return &pb.LoginResponse{Cookies: cookies}, err
}

func (f *Fbcolly) FetchGroupFeed(groupId int64, nextCursor string) (*pb.FacebookPostList, error) {
	collector := f.collector.Clone()
	var err error
	setupSharedCollector(collector, func(errInner error) {
		err = errInner
	})
	result := pb.FacebookPostList{Posts: []*pb.FacebookPost{}}

	collector.OnHTML("#m_group_stories_container > :last-child a", func(element *colly.HTMLElement) {
		result.NextCursor = "https://mbasic.facebook.com" + element.Attr("href")
	})
	collector.OnHTML("article", func(element *colly.HTMLElement) {
		dataElement := element
		post := &pb.FacebookPost{}
		var fbDataFt FbDataFt
		jsonData := dataElement.Attr("data-ft")

		logger.Info(jsonData)
		jsonErr := json.Unmarshal([]byte(jsonData), &fbDataFt)
		if jsonErr != nil {
			logger.Error(jsonErr)
			return
		}
		logger.Info("Post ", fbDataFt)
		post.Id = fbDataFt.TopLevelPostId
		post.Group = &pb.FacebookGroup{Id: fbDataFt.PageId, Name: dataElement.DOM.Find("h3 strong:nth-child(2) a").Text()}
		userId, _ := fbDataFt.ContentOwnerIdNew.Int64()
		if userId > 0 {
			userA := dataElement.DOM.Find("h3 strong:nth-child(1) a")
			post.User = &pb.FacebookUser{
				Id:       userId,
				Username: getUsernameFromHref(userA.AttrOr("href", "")),
				Name:     userA.Text(),
			}
			if len(post.User.Username) == 0 {
				logger.Error("Invalid username")
				return
			}

			createdAtWhenResult, err := f.w.Parse(element.DOM.Find("abbr").Text(), time.Now())
			if err != nil || createdAtWhenResult == nil {
				logger.Error(err)
				return
			}

			post.CreatedAt = createdAtWhenResult.Time.Unix()

			//Content

			//NO BACKGROUND TEXT ONLY
			post.Content = strings.Join(dataElement.DOM.Find("p").Map(func(i int, selection *goquery.Selection) string {
				return selection.Text()
			}), "\n")

			if len(post.Content) == 0 {
				// TEXT WITH BACKGROUND
				post.Content = dataElement.DOM.Find("div[style*=\"background-image:url\"]").Text()
			}

			if len(post.Content) == 0 {
				post.Content = dataElement.DOM.Find("[data-ft=\"{\\\"tn\\\":\\\"*s\\\"}\"] span :first-child").Text()
			}

			if len(post.Content) == 0 {
				post.Content = dataElement.DOM.Find("div[data-ft]").Text()
			}

			post.ContentLink = getUrlFromRedirectHref(dataElement.DOM.Find("a[href*=\"https://lm.facebook.com/l.php\"]").AttrOr("href", ""))
			post.ReactionCount = getNumberFromText(element.DOM.Find("span[id*=\"like_\"]").Text())
			post.CommentCount = getNumberFromText(element.DOM.Find("span[id*=\"like_\"] ~ a").Text())

			post.ContentImages = []*pb.FacebookImage{}
			element.DOM.Find("a[href*=\"/photo.php\"]").Each(func(i int, selection *goquery.Selection) {
				post.ContentImages = append(post.ContentImages, &pb.FacebookImage{
					Id: getImageIdFromHref(selection.AttrOr("href", "")),
				})
			})

			if len(post.ContentImages) > 0 {
				post.ContentImage = post.ContentImages[0]
			}
			result.Posts = append(result.Posts, post)
		}
	})
	if len(nextCursor) > 0 {
		_ = collector.Visit(nextCursor)
	} else {
		_ = collector.Visit(fmt.Sprintf("https://mbasic.facebook.com/groups/%d", groupId))
	}
	return &result, err
}

func (f *Fbcolly) FetchUserInfo(userIdOrUsername string) (*pb.FacebookUser, error) {
	collector := f.collector.Clone()
	var err error
	setupSharedCollector(collector, func(errInner error) {
		err = errInner
	})

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

	_ = collector.Visit(fmt.Sprintf("https://mbasic.facebook.com/%s", userIdOrUsername))
	return result, err
}

func (f *Fbcolly) FetchGroupInfo(groupIdOrUsername string) (*pb.FacebookGroup, error) {
	collector := f.collector.Clone()
	var err error
	setupSharedCollector(collector, func(errInner error) {
		err = errInner
	})
	result := &pb.FacebookGroup{}

	collector.OnHTML("a[href=\"#groupMenuBottom\"] h1", func(element *colly.HTMLElement) {
		result.Name = element.Text
	})
	collector.OnHTML("a[href*=\"view=member\"]", func(element *colly.HTMLElement) {
		result.Id = getNumberFromText(element.Attr("href"))
		result.MemberCount, _ = strconv.ParseInt(element.DOM.Closest("tr").Find("td:last-child").Text(), 10, 64)
	})

	_ = collector.Visit(fmt.Sprintf("https://mbasic.facebook.com/groups/%s?view=info", groupIdOrUsername))
	return result, err
}

func (f *Fbcolly) FetchContentImages(postId int64, nextCursor string) (*pb.FacebookImageList, error) {
	collector := f.collector.Clone()
	var err error
	setupSharedCollector(collector, func(errInner error) {
		err = errInner
	})
	result := pb.FacebookImageList{Images: []*pb.FacebookImage{}}

	collector.OnHTML("a[href*=\"/media/set/\"]", func(element *colly.HTMLElement) {
		result.NextCursor = "https://mbasic.facebook.com" + element.Attr("href")
	})

	collector.OnHTML("a[href*=\"/photo.php\"]", func(element *colly.HTMLElement) {
		result.Images = append(result.Images, &pb.FacebookImage{
			Id: getImageIdFromHref(element.Attr("href")),
		})
		//f.detailCollector.Visit(url)
	})
	if len(nextCursor) > 0 {
		_ = collector.Visit(nextCursor)
	} else {
		_ = collector.Visit(fmt.Sprintf("https://mbasic.facebook.com/media/set/?set=pcb.%d", postId))
	}
	return &result, err
}

func (f *Fbcolly) FetchImageUrl(imageId int64) (*pb.FacebookImage, error) {
	collector := f.collector.Clone()
	var err error
	setupSharedCollector(collector, func(errInner error) {
		err = errInner
	})
	result := pb.FacebookImage{Id: imageId}

	collector.OnHTML("a[href*=\"fbcdn\"]", func(element *colly.HTMLElement) {
		result.Url = element.Attr("href")
	})

	_ = collector.Visit(fmt.Sprintf("https://mbasic.facebook.com/photo/view_full_size/?fbid=%d", imageId))
	return &result, err
}

func (f *Fbcolly) FetchPost(groupId int64, postId int64, commentNextCursor string) (*pb.FacebookPost, error) {
	collector := f.collector.Clone()
	var err error
	setupSharedCollector(collector, func(errInner error) {
		err = errInner
	})
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
					Id:       userId,
					Username: getUsernameFromHref(dataElement.Find("h3 strong:first-child a").AttrOr("href", "")),
					Name:     dataElement.Find("h3 strong:first-child a").Text(),
				}

				createdAtWhenResult, err := f.w.Parse(element.DOM.Find("div[data-ft] abbr").Text(), time.Now())
				if err != nil || createdAtWhenResult == nil {
					logger.Error(err)
					return
				}

				post.CreatedAt = createdAtWhenResult.Time.Unix()
				//Content

				//NO BACKGROUND TEXT ONLY
				post.Content = strings.Join(dataElement.Find("p").Map(func(i int, selection *goquery.Selection) string {
					return selection.Text()
				}), "\n")

				if len(post.Content) == 0 {
					// TEXT WITH BACKGROUND
					post.Content = dataElement.Find("div[style*=\"background-image:url\"]").Text()
				}

				if len(post.Content) == 0 {
					post.Content = dataElement.Find("[data-ft=\"{\\\"tn\\\":\\\"*s\\\"}\"] span :first-child").Text()
				}

				if len(post.Content) == 0 {
					post.Content = dataElement.Find("div[data-ft]").Text()
				}

				post.ContentLink = getUrlFromRedirectHref(dataElement.Find("a[href*=\"https://lm.facebook.com/l.php\"]").AttrOr("href", ""))
				post.ReactionCount = getNumberFromText(element.DOM.Find("div[id*=\"sentence_\"]").Text())
				post.ContentImages = []*pb.FacebookImage{}
				element.DOM.Find("a[href*=\"/photo.php\"]").Each(func(i int, selection *goquery.Selection) {
					post.ContentImages = append(post.ContentImages, &pb.FacebookImage{
						Id: getImageIdFromHref(selection.AttrOr("href", "")),
					})
				})

				if len(post.ContentImages) > 0 {
					post.ContentImage = post.ContentImages[0]
				}
			}

			//Comment
			element.DOM.Find("h3 + div + div + div").Parent().Parent().Each(func(i int, selection *goquery.Selection) {
				//author
				commentId, _ := strconv.ParseInt(selection.AttrOr("id", ""), 10, 64)
				if commentId > 0 {
					createdAtWhenResult, err := f.w.Parse(selection.Find("abbr").Text(), time.Now())
					if err != nil || createdAtWhenResult == nil {
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
						comment := pb.FacebookComment{
							Id:   commentId,
							Post: &pb.FacebookPost{Id: post.Id},
							User: &pb.FacebookUser{
								Username: getUsernameFromHref(selection.Find("h3 > a").AttrOr("href", "")),
								Name:     selection.Find("h3 > a").Text(),
							},
							Content:   selection.Find("h3 + div").Text(),
							CreatedAt: createdAtWhenResult.Time.Unix(),
						}
						post.Comments.Comments = append(post.Comments.Comments, &comment)
					}
				}
			})

		}
	})

	collector.OnHTML("div[id*=\"see_prev_\"] > a", func(element *colly.HTMLElement) {
		post.Comments.NextCursor = element.Attr("href")
	})
	if len(commentNextCursor) > 0 {
		_ = collector.Visit(commentNextCursor)
	} else {
		_ = collector.Visit(fmt.Sprintf("https://mbasic.facebook.com/groups/%d?view=permalink&id=%d&_rdr", groupId, postId))
	}

	return post, err
}

func (f *Fbcolly) LoginWithCookies(cookies string) error {
	collector := f.collector
	return collector.SetCookies("https://mbasic.facebook.com/", storage.UnstringifyCookies(cookies))
}

func (f *Fbcolly) FetchMyGroups() (*pb.FacebookGroupList, error) {
	collector := f.collector.Clone()
	var err error
	setupSharedCollector(collector, func(errInner error) {
		err = errInner
	})
	result := &pb.FacebookGroupList{Groups: []*pb.FacebookGroup{}}

	collector.OnHTML("li table a", func(element *colly.HTMLElement) {
		result.Groups = append(result.Groups, &pb.FacebookGroup{
			Id:   getNumberFromText(element.Attr("href")),
			Name: element.Text,
		})
	})

	_ = collector.Visit("https://mbasic.facebook.com/groups/?seemore")
	return result, err
}

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

func getUsernameFromHref(href string) string {
	parsed, _ := url.Parse(href)
	if strings.HasPrefix(parsed.Path, "/profile.php") {
		return parsed.Query().Get("id")
	} else {
		if len(parsed.Path) > 0 {
			return strings.Split(parsed.Path[1:], "/")[0]
		}

	}
	return ""
}
