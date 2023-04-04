package reddit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type CommentResponse []struct {
	Kind string `json:"kind"`
	Data struct {
		After     any    `json:"after"`
		Dist      int    `json:"dist"`
		Modhash   any    `json:"modhash"`
		GeoFilter string `json:"geo_filter"`
		Children  []struct {
			Kind string `json:"kind"`
			Data struct {
				ApprovedAtUtc     any    `json:"approved_at_utc"`
				Body              string `json: "body"`
				Subreddit         string `json:"subreddit"`
				Selftext          string `json:"selftext"`
				UserReports       []any  `json:"user_reports"`
				Saved             bool   `json:"saved"`
				ModReasonTitle    any    `json:"mod_reason_title"`
				Gilded            int    `json:"gilded"`
				Clicked           bool   `json:"clicked"`
				Title             string `json:"title"`
				LinkFlairRichtext []struct {
					E string `json:"e"`
					T string `json:"t"`
				} `json:"link_flair_richtext"`
				SubredditNamePrefixed      string  `json:"subreddit_name_prefixed"`
				Hidden                     bool    `json:"hidden"`
				Pwls                       int     `json:"pwls"`
				LinkFlairCSSClass          string  `json:"link_flair_css_class"`
				Downs                      int     `json:"downs"`
				TopAwardedType             any     `json:"top_awarded_type"`
				ParentWhitelistStatus      string  `json:"parent_whitelist_status"`
				HideScore                  bool    `json:"hide_score"`
				Name                       string  `json:"name"`
				Quarantine                 bool    `json:"quarantine"`
				LinkFlairTextColor         string  `json:"link_flair_text_color"`
				UpvoteRatio                float64 `json:"upvote_ratio"`
				AuthorFlairBackgroundColor any     `json:"author_flair_background_color"`
				SubredditType              string  `json:"subreddit_type"`
				Ups                        int     `json:"ups"`
				TotalAwardsReceived        int     `json:"total_awards_received"`
				MediaEmbed                 struct {
				} `json:"media_embed"`
				AuthorFlairTemplateID any    `json:"author_flair_template_id"`
				IsOriginalContent     bool   `json:"is_original_content"`
				AuthorFullname        string `json:"author_fullname"`
				SecureMedia           any    `json:"secure_media"`
				IsRedditMediaDomain   bool   `json:"is_reddit_media_domain"`
				IsMeta                bool   `json:"is_meta"`
				Category              any    `json:"category"`
				SecureMediaEmbed      struct {
				} `json:"secure_media_embed"`
				LinkFlairText       string `json:"link_flair_text"`
				CanModPost          bool   `json:"can_mod_post"`
				Score               int    `json:"score"`
				ApprovedBy          any    `json:"approved_by"`
				IsCreatedFromAdsUI  bool   `json:"is_created_from_ads_ui"`
				AuthorPremium       bool   `json:"author_premium"`
				Thumbnail           string `json:"thumbnail"`
				Edited              any    `json:"edited"`
				AuthorFlairCSSClass any    `json:"author_flair_css_class"`
				AuthorFlairRichtext []any  `json:"author_flair_richtext"`
				Gildings            struct {
				} `json:"gildings"`
				ContentCategories   any     `json:"content_categories"`
				IsSelf              bool    `json:"is_self"`
				ModNote             any     `json:"mod_note"`
				Created             float64 `json:"created"`
				LinkFlairType       string  `json:"link_flair_type"`
				Wls                 int     `json:"wls"`
				RemovedByCategory   any     `json:"removed_by_category"`
				BannedBy            any     `json:"banned_by"`
				AuthorFlairType     string  `json:"author_flair_type"`
				Domain              string  `json:"domain"`
				AllowLiveComments   bool    `json:"allow_live_comments"`
				SelftextHTML        any     `json:"selftext_html"`
				Likes               any     `json:"likes"`
				SuggestedSort       any     `json:"suggested_sort"`
				BannedAtUtc         any     `json:"banned_at_utc"`
				URLOverriddenByDest string  `json:"url_overridden_by_dest"`
				ViewCount           any     `json:"view_count"`
				Archived            bool    `json:"archived"`
				NoFollow            bool    `json:"no_follow"`
				IsCrosspostable     bool    `json:"is_crosspostable"`
				Pinned              bool    `json:"pinned"`
				Over18              bool    `json:"over_18"`
				AllAwardings        []struct {
					GiverCoinReward          any    `json:"giver_coin_reward"`
					SubredditID              any    `json:"subreddit_id"`
					IsNew                    bool   `json:"is_new"`
					DaysOfDripExtension      any    `json:"days_of_drip_extension"`
					CoinPrice                int    `json:"coin_price"`
					ID                       string `json:"id"`
					PennyDonate              any    `json:"penny_donate"`
					CoinReward               int    `json:"coin_reward"`
					IconURL                  string `json:"icon_url"`
					DaysOfPremium            any    `json:"days_of_premium"`
					IconHeight               int    `json:"icon_height"`
					TiersByRequiredAwardings any    `json:"tiers_by_required_awardings"`
					ResizedIcons             []struct {
						URL    string `json:"url"`
						Width  int    `json:"width"`
						Height int    `json:"height"`
					} `json:"resized_icons"`
					IconWidth                        int    `json:"icon_width"`
					StaticIconWidth                  int    `json:"static_icon_width"`
					StartDate                        any    `json:"start_date"`
					IsEnabled                        bool   `json:"is_enabled"`
					AwardingsRequiredToGrantBenefits any    `json:"awardings_required_to_grant_benefits"`
					Description                      string `json:"description"`
					EndDate                          any    `json:"end_date"`
					StickyDurationSeconds            any    `json:"sticky_duration_seconds"`
					SubredditCoinReward              int    `json:"subreddit_coin_reward"`
					Count                            int    `json:"count"`
					StaticIconHeight                 int    `json:"static_icon_height"`
					Name                             string `json:"name"`
					ResizedStaticIcons               []struct {
						URL    string `json:"url"`
						Width  int    `json:"width"`
						Height int    `json:"height"`
					} `json:"resized_static_icons"`
					IconFormat    any    `json:"icon_format"`
					AwardSubType  string `json:"award_sub_type"`
					PennyPrice    any    `json:"penny_price"`
					AwardType     string `json:"award_type"`
					StaticIconURL string `json:"static_icon_url"`
				} `json:"all_awardings"`
				Awarders                 []any   `json:"awarders"`
				MediaOnly                bool    `json:"media_only"`
				CanGild                  bool    `json:"can_gild"`
				Spoiler                  bool    `json:"spoiler"`
				Locked                   bool    `json:"locked"`
				AuthorFlairText          any     `json:"author_flair_text"`
				TreatmentTags            []any   `json:"treatment_tags"`
				Visited                  bool    `json:"visited"`
				RemovedBy                any     `json:"removed_by"`
				NumReports               any     `json:"num_reports"`
				Distinguished            any     `json:"distinguished"`
				SubredditID              string  `json:"subreddit_id"`
				AuthorIsBlocked          bool    `json:"author_is_blocked"`
				ModReasonBy              any     `json:"mod_reason_by"`
				RemovalReason            any     `json:"removal_reason"`
				LinkFlairBackgroundColor string  `json:"link_flair_background_color"`
				ID                       string  `json:"id"`
				IsRobotIndexable         bool    `json:"is_robot_indexable"`
				NumDuplicates            int     `json:"num_duplicates"`
				ReportReasons            any     `json:"report_reasons"`
				Author                   string  `json:"author"`
				DiscussionType           any     `json:"discussion_type"`
				NumComments              int     `json:"num_comments"`
				SendReplies              bool    `json:"send_replies"`
				Media                    any     `json:"media"`
				ContestMode              bool    `json:"contest_mode"`
				AuthorPatreonFlair       bool    `json:"author_patreon_flair"`
				AuthorFlairTextColor     any     `json:"author_flair_text_color"`
				Permalink                string  `json:"permalink"`
				WhitelistStatus          string  `json:"whitelist_status"`
				Stickied                 bool    `json:"stickied"`
				URL                      string  `json:"url"`
				SubredditSubscribers     int     `json:"subreddit_subscribers"`
				CreatedUtc               float64 `json:"created_utc"`
				NumCrossposts            int     `json:"num_crossposts"`
				ModReports               []any   `json:"mod_reports"`
				IsVideo                  bool    `json:"is_video"`
			} `json:"data"`
		} `json:"children"`
		Before any `json:"before"`
	} `json:"data"`
}

func LoadComments(subReddit string, postId string, sortingMethod string, depth int, bearerToken string) []string {

	r, err := http.NewRequest("GET", "https://oauth.reddit.com/r/"+subReddit+"/comments/"+postId+"?sort="+sortingMethod+"&depth="+strconv.Itoa(depth), nil)

	if err != nil {
		log.Println("Error while creating new request")
		panic(err)
	}

	r.Header.Add("Authorization", "bearer "+bearerToken)
	r.Header.Add("User-Agent", "posts_summarizer/0.0.1")

	client := &http.Client{}

	res, err := client.Do(r)

	if err != nil {
		log.Println("Error while getting response from Reddit API")
		panic(err)
	}

	defer res.Body.Close()

	responseData, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Println("Error in reading raw comments from response")
		log.Fatal(err)
	}
	fmt.Println("Response Status:", res.Status)

	var resJson CommentResponse

	err = json.Unmarshal(responseData, &resJson)

	if err != nil {
		log.Println("Error in unmarshalling")
		log.Fatal(err)
	}

	var comments []string
	for _, v := range resJson[1].Data.Children {
		comments = append(comments, v.Data.Body)
	}

	return comments
}
