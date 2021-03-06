package acfundanmu

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fastjson"
)

// WatchingUser 就是观看直播的用户的信息，目前没有Medal
type WatchingUser struct {
	UserInfo          `json:"userInfo"`
	AnonymousUser     bool   `json:"anonymousUser"`     // 是否匿名用户
	DisplaySendAmount string `json:"displaySendAmount"` // 赠送的全部礼物的价值，单位是AC币，注意不一定是纯数字的字符串
	CustomData        string `json:"customData"`        // 用户的一些额外信息，格式为json
}

// BillboardUser 就是礼物贡献榜上的用户的信息，没有AnonymousUser、Medal和ManagerType
type BillboardUser WatchingUser

// Summary 就是直播的总结信息
type Summary struct {
	Duration     int64 `json:"duration"`     // 直播时长，单位为毫秒
	LikeCount    int   `json:"likeCount"`    // 点赞总数
	WatchCount   int   `json:"watchCount"`   // 观看过直播的人数总数
	GiftCount    int   `json:"giftCount"`    // 直播收到的付费礼物数量
	DiamondCount int   `json:"diamondCount"` // 直播收到的付费礼物对应的钻石数量，100钻石=1AC币
	BananaCount  int   `json:"bananaCount"`  // 直播收到的香蕉数量
}

// Medal 就是登陆帐号的守护徽章信息
type Medal struct {
	MedalInfo          `json:"medalInfo"`
	UperName           string `json:"uperName"`           // UP主的名字
	UperAvatar         string `json:"uperAvatar"`         // UP主的头像
	WearMedal          bool   `json:"wearMedal"`          // 是否正在佩戴该守护徽章
	FriendshipDegree   int    `json:"friendshipDegree"`   // 目前守护徽章的亲密度
	JoinClubTime       int64  `json:"joinClubTime"`       // 加入守护团的时间，是以毫秒为单位的Unix时间
	CurrentDegreeLimit int    `json:"currentDegreeLimit"` // 守护徽章目前等级的亲密度的上限
}

// MedalDegree 就是守护徽章的亲密度信息
type MedalDegree struct {
	UperID               int64 `json:"uperID"`               // UP主的uid
	GiftDegree           int   `json:"giftDegree"`           // 本日送直播礼物增加的亲密度
	GiftDegreeLimit      int   `json:"giftDegreeLimit"`      // 本日送直播礼物增加的亲密度上限
	PeachDegree          int   `json:"peachDegree"`          // 本日投桃增加的亲密度
	PeachDegreeLimit     int   `json:"peachDegreeLimit"`     // 本日投桃增加的亲密度上限
	LiveWatchDegree      int   `json:"liveWatchDegree"`      // 本日看直播时长增加的亲密度
	LiveWatchDegreeLimit int   `json:"liveWatchDegreeLimit"` // 本日看直播时长增加的亲密度上限
	BananaDegree         int   `json:"bananaDegree"`         // 本日投蕉增加的亲密度
	BananaDegreeLimit    int   `json:"bananaDegreeLimit"`    // 本日投蕉增加的亲密度上限
}

// MedalDetail 就是登陆用户的守护徽章的详细信息
type MedalDetail struct {
	Medal       Medal       `json:"medal"`       // 守护徽章信息
	MedalDegree MedalDegree `json:"medalDegree"` // 守护徽章亲密度信息
	UserRank    int         `json:"userRank"`    // 登陆用户的主播守护徽章亲密度的排名
	ClubName    string      `json:"clubName"`    // 守护徽章名字
}

// MedalList 就是登陆用户拥有的守护徽章列表
type MedalList struct {
	MedalList   []Medal     `json:"medalList"`   // 用户拥有的守护徽章列表
	MedalDetail MedalDetail `json:"medalDetail"` // 用户的指定主播的守护徽章详细信息
}

// LuckyUser 就是抢到红包的用户，没有Medal和ManagerType
type LuckyUser struct {
	UserInfo   `json:"userInfo"`
	GrabAmount int `json:"grabAmount"` // 抢红包抢到的AC币
}

// Playback 就是直播回放的相关信息
type Playback struct {
	Duration  int64  `json:"duration"`  // 录播视频时长，单位是毫秒
	URL       string `json:"url"`       // 录播源链接，目前分为阿里云和腾讯云两种，目前阿里云的下载速度比较快
	BackupURL string `json:"backupURL"` // 备份录播源链接
	M3U8Slice string `json:"m3u8Slice"` // m3u8
	Width     int    `json:"width"`     // 录播视频宽度
	Height    int    `json:"height"`    // 录播视频高度
}

// Manager 就是房管的用户信息，目前没有Medal和ManagerType
type Manager struct {
	UserInfo   `json:"userInfo"`
	CustomData string `json:"customData"` // 用户的一些额外信息，格式为json
	Online     bool   `json:"online"`     // 是否直播间在线？
}

// UserProfile 就是用户信息
type UserProfile struct {
	UserID          int64  `json:"userID"`          // 用户uid
	Nickname        string `json:"nickname"`        // 用户名字
	Avatar          string `json:"avatar"`          // 用户头像
	AvatarFrame     string `json:"avatarFrame"`     // 用户头像挂件
	FollowingCount  int    `json:"followingCount"`  // 用户关注数量
	FansCount       int    `json:"fansCount"`       // 用户粉丝数量
	ContributeCount int    `json:"contributeCount"` // 用户投稿数量
	Signature       string `json:"signature"`       // 用户签名
	VerifiedText    string `json:"verifiedText"`    // 用户认证信息
	IsJoinUpCollege bool   `json:"isJoinUpCollege"` // 用户是否加入阿普学院
	IsFollowing     bool   `json:"isFollowing"`     // 是否关注了用户
}

// UserLiveInfo 就是用户直播信息
type UserLiveInfo struct {
	Profile               UserProfile `json:"profile"`               // 用户信息
	LiveType              LiveType    `json:"liveType"`              // 直播分类
	LiveID                string      `json:"liveID"`                // 直播ID
	StreamName            string      `json:"streamName"`            // 直播源名字（ID）
	Title                 string      `json:"title"`                 // 直播间标题
	LiveStartTime         int64       `json:"liveStartTime"`         // 直播开始的时间，是以毫秒为单位的Unix时间
	Portrait              bool        `json:"portrait"`              // 是否手机直播
	Panoramic             bool        `json:"panoramic"`             // 是否全景直播
	LiveCover             string      `json:"liveCover"`             // 直播间封面
	OnlineCount           int         `json:"onlineCount"`           // 直播间在线人数
	LikeCount             int         `json:"likeCount"`             // 直播间点赞总数
	HasFansClub           bool        `json:"hasFansClub"`           // 主播是否有守护团
	DisableDanmakuShow    bool        `json:"disableDanmakuShow"`    // 是否禁止显示弹幕？
	PaidShowUserBuyStatus bool        `json:"paidShowUserBuyStatus"` // 用户是否购买了付费直播？
}

// UserMedalInfo 就是用户的守护徽章信息
type UserMedalInfo struct {
	UserProfile      `json:"profile"`
	FriendshipDegree int `json:"friendshipDegree"` // 用户守护徽章的亲密度
	Level            int `json:"level"`            // 用户守护徽章的等级
}

// MedalRankList 就是主播守护徽章等级列表
type MedalRankList struct {
	RankList             []UserMedalInfo `json:"rankList"`             // 主播守护徽章等级列表
	ClubName             string          `json:"clubName"`             // 守护徽章名字
	MedalCount           int             `json:"medalCount"`           // 拥有主播守护徽章的用户的数量
	HasMedal             bool            `json:"hasMedal"`             // 登陆用户是否有主播的守护徽章
	UserFriendshipDegree int             `json:"userFriendshipDegree"` // 登陆用户的主播守护徽章的亲密度
	UserRank             int             `json:"userRank"`             // 登陆用户的主播守护徽章的排名
}

// 获取直播间排名前50的在线观众信息列表
func (t *token) getWatchingList(liveID string) (watchingList []WatchingUser, e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("watchingList() error: %w", err)
		}
	}()

	form := t.defaultForm(liveID)
	defer fasthttp.ReleaseArgs(form)
	body, err := t.fetchKuaiShouAPI(watchingListURL, form, false)
	checkErr(err)

	p := t.watchParser.Get()
	defer t.watchParser.Put(p)
	v, err := p.ParseBytes(body)
	checkErr(err)
	if v.GetInt("result") != 1 {
		panic(fmt.Errorf("获取在线观众列表失败，响应为 %s", string(body)))
	}

	watchArray := v.GetArray("data", "list")
	watchingList = make([]WatchingUser, len(watchArray))
	for i, watch := range watchArray {
		o := watch.GetObject()
		o.Visit(func(k []byte, v *fastjson.Value) {
			switch string(k) {
			case "userId":
				watchingList[i].UserID = v.GetInt64()
			case "nickname":
				watchingList[i].Nickname = string(v.GetStringBytes())
			case "avatar":
				watchingList[i].Avatar = string(v.GetStringBytes("0", "url"))
			case "anonymousUser":
				watchingList[i].AnonymousUser = v.GetBool()
			case "displaySendAmount":
				watchingList[i].DisplaySendAmount = string(v.GetStringBytes())
			case "customWatchingListData":
				watchingList[i].CustomData = string(v.GetStringBytes())
			case "managerType":
				watchingList[i].ManagerType = ManagerType(v.GetInt())
			default:
				log.Printf("在线观众列表里出现未处理的key和value：%s %s", string(k), string(v.MarshalTo([]byte{})))
			}
		})
	}

	return watchingList, nil
}

// 获取主播最近七日内的礼物贡献榜前50名观众的详细信息
func (t *token) getBillboard(uid int64) (billboard []BillboardUser, e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("getBillboard() error: %w", err)
		}
	}()

	form := fasthttp.AcquireArgs()
	defer fasthttp.ReleaseArgs(form)
	form.Set("authorId", strconv.FormatInt(uid, 10))
	body, err := t.fetchKuaiShouAPI(billboardURL, form, false)
	checkErr(err)

	p := t.watchParser.Get()
	defer t.watchParser.Put(p)
	v, err := p.ParseBytes(body)
	checkErr(err)
	if v.GetInt("result") != 1 {
		panic(fmt.Errorf("获取主播最近七日内的礼物贡献榜失败，响应为 %s", string(body)))
	}

	billboardArray := v.GetArray("data", "list")
	billboard = make([]BillboardUser, len(billboardArray))
	for i, user := range billboardArray {
		o := user.GetObject()
		o.Visit(func(k []byte, v *fastjson.Value) {
			switch string(k) {
			case "userId":
				billboard[i].UserID = v.GetInt64()
			case "nickname":
				billboard[i].Nickname = string(v.GetStringBytes())
			case "avatar":
				billboard[i].Avatar = string(v.GetStringBytes("0", "url"))
			case "displaySendAmount":
				billboard[i].DisplaySendAmount = string(v.GetStringBytes())
			case "customData":
				billboard[i].CustomData = string(v.GetStringBytes())
			default:
				log.Printf("礼物贡献榜里出现未处理的key和value：%s %s", string(k), string(v.MarshalTo([]byte{})))
			}
		})
	}

	return billboard, nil
}

// 获取直播总结信息
func (t *token) getSummary(liveID string) (summary *Summary, e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("getSummary() error: %w", err)
		}
	}()

	form := t.defaultForm(liveID)
	defer fasthttp.ReleaseArgs(form)
	body, err := t.fetchKuaiShouAPI(endSummaryURL, form, false)
	checkErr(err)

	var p fastjson.Parser
	v, err := p.ParseBytes(body)
	checkErr(err)
	if v.GetInt("result") != 1 {
		panic(fmt.Errorf("获取直播总结信息失败，响应为 %s", string(body)))
	}

	summary = new(Summary)
	o := v.GetObject("data")
	o.Visit(func(k []byte, v *fastjson.Value) {
		switch string(k) {
		case "liveDurationMs":
			summary.Duration = v.GetInt64()
		case "likeCount":
			summary.LikeCount, err = strconv.Atoi(string(v.GetStringBytes()))
			checkErr(err)
		case "watchCount":
			summary.WatchCount, err = strconv.Atoi(string(v.GetStringBytes()))
			checkErr(err)
		case "payWalletTypeToReceiveCount":
			summary.GiftCount = v.GetInt("1")
		case "payWalletTypeToReceiveCurrency":
			summary.DiamondCount = v.GetInt("1")
			summary.BananaCount = v.GetInt("2")
		default:
			log.Printf("直播总结信息里出现未处理的key和value：%s %s", string(k), string(v.MarshalTo([]byte{})))
		}
	})

	return summary, nil
}

// 获取抢到红包的用户列表
func (t *token) getLuckList(liveID, redpackID string) (luckyList []LuckyUser, e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("getLuckList() error: %w", err)
		}
	}()

	if len(t.Cookies) == 0 {
		panic(fmt.Errorf("获取抢到红包的用户列表需要登陆AcFun帐号"))
	}

	form := t.defaultForm(liveID)
	defer fasthttp.ReleaseArgs(form)
	form.Set("redpackBizUnit", "ztLiveAcfunRedpackGift")
	form.Set("redpackId", redpackID)
	body, err := t.fetchKuaiShouAPI(redpackLuckListURL, form, false)
	checkErr(err)

	p := t.watchParser.Get()
	defer t.watchParser.Put(p)
	v, err := p.ParseBytes(body)
	checkErr(err)
	if v.GetInt("result") != 1 {
		panic(fmt.Errorf("获取抢到红包的用户列表失败，响应为 %s", string(body)))
	}

	luckyArray := v.GetArray("data", "luckyList")
	luckyList = make([]LuckyUser, len(luckyArray))
	for i, user := range luckyArray {
		o := user.GetObject()
		o.Visit(func(k []byte, v *fastjson.Value) {
			switch string(k) {
			case "simpleUserInfo":
				o := v.GetObject()
				o.Visit(func(k []byte, v *fastjson.Value) {
					switch string(k) {
					case "userId":
						luckyList[i].UserID = v.GetInt64()
					case "nickname":
						luckyList[i].Nickname = string(v.GetStringBytes())
					case "headPic":
						luckyList[i].Avatar = string(v.GetStringBytes("0", "url"))
					default:
						log.Printf("抢到红包的用户列表里出现未处理的key和value：%s %s", string(k), string(v.MarshalTo([]byte{})))
					}
				})
			case "grabAmount":
				luckyList[i].GrabAmount = v.GetInt()
			default:
				log.Printf("抢到红包的用户列表里出现未处理的key和value：%s %s", string(k), string(v.MarshalTo([]byte{})))
			}
		})
	}

	return luckyList, nil
}

// 获取直播回放的相关信息
func (t *token) getPlayback(liveID string) (playback *Playback, e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("getPlayback() error: %w", err)
		}
	}()

	form := fasthttp.AcquireArgs()
	defer fasthttp.ReleaseArgs(form)
	form.Set("liveId", liveID)
	body, err := t.fetchKuaiShouAPI(playbackURL, form, true)
	checkErr(err)

	var p fastjson.Parser
	v, err := p.ParseBytes(body)
	checkErr(err)
	if v.GetInt("result") != 1 {
		panic(fmt.Errorf("获取直播回放的相关信息失败，响应为 %s", string(body)))
	}
	adaptiveManifest := v.GetStringBytes("data", "adaptiveManifest")
	v, err = p.ParseBytes(adaptiveManifest)
	checkErr(err)
	if len(v.GetArray("adaptationSet")) > 1 {
		log.Println("adaptationSet列表长度大于1，请报告issue")
	}
	v = v.Get("adaptationSet", "0")
	duration := v.GetInt64("duration")
	if len(v.GetArray("representation")) > 1 {
		log.Println("representation列表长度大于1，请报告issue")
	}
	v = v.Get("representation", "0")
	playback = &Playback{
		Duration:  duration,
		URL:       string(v.GetStringBytes("url")),
		BackupURL: string(v.GetStringBytes("backupUrl", "0")),
		M3U8Slice: string(v.GetStringBytes("m3u8Slice")),
		Width:     v.GetInt("width"),
		Height:    v.GetInt("height"),
	}
	if len(v.GetArray("backupUrl")) > 1 {
		log.Println("backupUrl列表长度大于1，请报告issue")
	}

	return playback, nil
}

// 获取直播源信息，和getLiveToken()重复了
func (t *token) getPlayURL() (e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("getPlayURL() error: %w", err)
		}
	}()

	body, err := t.fetchKuaiShouAPI(getPlayURL, nil, false)
	checkErr(err)

	var p fastjson.Parser
	v, err := p.ParseBytes(body)
	checkErr(err)
	if v.GetInt("result") != 1 {
		panic(fmt.Errorf("获取直播源信息失败，响应为 %s", string(body)))
	}

	//videoPlayRes := string(v.GetStringBytes("data", "videoPlayRes"))

	return nil
}

// 获取全部礼物的数据
func (t *token) getAllGift() (gifts map[int64]GiftDetail, e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("getAllGift() error: %w", err)
		}
	}()

	form := fasthttp.AcquireArgs()
	defer fasthttp.ReleaseArgs(form)
	form.Set("visitorId", strconv.FormatInt(t.UserID, 10))
	body, err := t.fetchKuaiShouAPI(allGiftURL, form, false)
	checkErr(err)

	var p fastjson.Parser
	v, err := p.ParseBytes(body)
	checkErr(err)
	if v.GetInt("result") != 1 {
		panic(fmt.Errorf("获取全部礼物的数据失败，响应为 %s", string(body)))
	}

	gifts = updateGiftList(v)

	return gifts, nil
}

// 获取钱包里AC币和拥有的香蕉的数量
func (t *token) getWalletBalance() (accoins int, bananas int, e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("getWalletBalance() error: %w", err)
		}
	}()

	if len(t.Cookies) == 0 {
		panic(fmt.Errorf("获取钱包里AC币和拥有的香蕉的数量需要登陆AcFun帐号"))
	}

	form := fasthttp.AcquireArgs()
	defer fasthttp.ReleaseArgs(form)
	form.Set("visitorId", strconv.FormatInt(t.UserID, 10))
	body, err := t.fetchKuaiShouAPI(walletBalanceURL, form, false)
	checkErr(err)

	var p fastjson.Parser
	v, err := p.ParseBytes(body)
	checkErr(err)
	if v.GetInt("result") != 1 {
		panic(fmt.Errorf("获取拥有的香蕉和钱包里AC币的数量失败，响应为 %s", string(body)))
	}

	o := v.GetObject("data", "payWalletTypeToBalance")
	o.Visit(func(k []byte, v *fastjson.Value) {
		switch string(k) {
		case "1":
			accoins = v.GetInt()
		case "2":
			bananas = v.GetInt()
		default:
			log.Printf("用户钱包里出现未处理的key和value：%s %s", string(k), string(v.MarshalTo([]byte{})))
		}
	})

	return accoins, bananas, nil
}

// 获取主播踢人的历史记录
func (t *token) getKickHistory() (e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("getAuthorKickHistory() error: %w", err)
		}
	}()

	if len(t.Cookies) == 0 {
		panic(fmt.Errorf("获取主播踢人的历史记录需要登陆主播的AcFun帐号"))
	}

	body, err := t.fetchKuaiShouAPI(kickHistoryURL, nil, false)
	checkErr(err)

	var p fastjson.Parser
	v, err := p.ParseBytes(body)
	checkErr(err)
	if v.GetInt("result") != 1 {
		panic(fmt.Errorf("获取主播踢人的历史记录失败，响应为 %s", string(body)))
	} else {
		log.Printf("获取主播踢人的历史记录的响应为 %s", string(body))
	}

	return nil
}

// 获取主播的房管列表
func (t *token) getManagerList() (managerList []Manager, e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("getAuthorManagerList() error: %w", err)
		}
	}()

	if len(t.Cookies) == 0 {
		panic(fmt.Errorf("获取主播的房管列表需要登陆主播的AcFun帐号"))
	}

	form := fasthttp.AcquireArgs()
	defer fasthttp.ReleaseArgs(form)
	form.Set("visitorId", strconv.FormatInt(t.UserID, 10))
	body, err := t.fetchKuaiShouAPI(managerListURL, form, false)
	checkErr(err)

	var p fastjson.Parser
	v, err := p.ParseBytes(body)
	checkErr(err)
	if v.GetInt("result") != 1 {
		panic(fmt.Errorf("获取主播的房管列表失败，响应为 %s", string(body)))
	}

	list := v.GetArray("data", "list")
	managerList = make([]Manager, len(list))
	for i, m := range list {
		o := m.GetObject()
		o.Visit(func(k []byte, v *fastjson.Value) {
			switch string(k) {
			case "userId":
				managerList[i].UserID = v.GetInt64()
			case "nickname":
				managerList[i].Nickname = string(v.GetStringBytes())
			case "avatar":
				managerList[i].Avatar = string(v.GetStringBytes("0", "url"))
			case "customData":
				managerList[i].CustomData = string(v.GetStringBytes())
			case "online":
				managerList[i].Online = v.GetBool()
			default:
				log.Printf("主播的房管列表里出现未处理的key和value：%s %s", string(k), string(v.MarshalTo([]byte{})))
			}
		})
	}

	return managerList, nil
}

// 从json里获取登陆帐号的守护徽章信息，uid是登陆帐号的uid
func getMedalJSON(v *fastjson.Value, uid int64) *Medal {
	medal := new(Medal)
	o := v.GetObject()
	o.Visit(func(k []byte, v *fastjson.Value) {
		switch string(k) {
		case "uperId":
			medal.UperID = v.GetInt64()
		case "clubName":
			medal.ClubName = string(v.GetStringBytes())
		case "level":
			medal.Level = v.GetInt()
		case "uperName":
			medal.UperName = string(v.GetStringBytes())
		case "uperHeadUrl":
			medal.UperAvatar = string(v.GetStringBytes())
		case "wearMedal":
			medal.WearMedal = v.GetBool()
		case "friendshipDegree":
			medal.FriendshipDegree = v.GetInt()
		case "joinClubTime":
			medal.JoinClubTime = v.GetInt64()
		case "currentDegreeLimit":
			medal.CurrentDegreeLimit = v.GetInt()
		default:
			log.Printf("守护徽章信息里出现未处理的key和value：%s %s", string(k), string(v.MarshalTo([]byte{})))
		}
	})
	medal.UserID = uid

	return medal
}

// 从json里获取登陆帐号的守护徽章亲密度信息
func getMedalDegreeJSON(v *fastjson.Value) *MedalDegree {
	medal := new(MedalDegree)
	o := v.GetObject()
	o.Visit(func(k []byte, v *fastjson.Value) {
		switch string(k) {
		case "uperId":
			medal.UperID = v.GetInt64()
		case "giftDegree":
			medal.GiftDegree = v.GetInt()
		case "giftDegreeLimit":
			medal.GiftDegreeLimit = v.GetInt()
		case "peachDegree":
			medal.PeachDegree = v.GetInt()
		case "peachDegreeLimit":
			medal.PeachDegreeLimit = v.GetInt()
		case "liveWatchDegree":
			medal.LiveWatchDegree = v.GetInt()
		case "liveWatchDegreeLimit":
			medal.LiveWatchDegreeLimit = v.GetInt()
		case "bananaDegree":
			medal.BananaDegree = v.GetInt()
		case "bananaDegreeLimit":
			medal.BananaDegreeLimit = v.GetInt()
		default:
			log.Printf("守护徽章亲密度信息里出现未处理的key和value：%s %s", string(k), string(v.MarshalTo([]byte{})))
		}
	})

	return medal
}

// 获取登陆帐号拥有的指定主播的守护徽章详细信息
func (t *token) getMedalDetail(uid int64) (medal *MedalDetail, e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("getMedalDetail() error: %w", err)
		}
	}()

	if len(t.Cookies) == 0 {
		panic(fmt.Errorf("获取登陆帐号拥有的指定主播的守护徽章详细信息需要登陆AcFun帐号"))
	}

	client := &httpClient{
		url:     fmt.Sprintf(medalDetailURL, uid),
		method:  "GET",
		cookies: t.Cookies,
	}
	body, err := client.request()
	checkErr(err)

	var p fastjson.Parser
	v, err := p.ParseBytes(body)
	checkErr(err)
	if !v.Exists("result") || v.GetInt("result") != 0 {
		panic(fmt.Errorf("获取登陆帐号拥有的指定主播的守护徽章详细信息失败，响应为 %s", string(body)))
	}

	medal = new(MedalDetail)
	o := v.GetObject()
	o.Visit(func(k []byte, v *fastjson.Value) {
		switch string(k) {
		case "result":
		case "host-name":
		case "medal":
			medal.Medal = *getMedalJSON(v, t.UserID)
		case "medalDegreeLimit":
			medal.MedalDegree = *getMedalDegreeJSON(v)
		case "rankIndex":
			medal.UserRank, err = strconv.Atoi(string(v.GetStringBytes()))
			checkErr(err)
		default:
			log.Printf("守护徽章详细信息里出现未处理的key和value：%s %s", string(k), string(v.MarshalTo([]byte{})))
		}
	})
	medal.ClubName = medal.Medal.ClubName

	return medal, nil
}

// 获取登陆帐号拥有的守护徽章列表
func (t *token) getMedalList(uid int64) (medalList *MedalList, e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("getMedalInfo() error: %w", err)
		}
	}()

	if len(t.Cookies) == 0 {
		panic(fmt.Errorf("获取登陆帐号拥有的守护徽章列表需要登陆AcFun帐号"))
	}

	client := &httpClient{
		url:     fmt.Sprintf(medalListURL, uid),
		method:  "GET",
		cookies: t.Cookies,
	}
	body, err := client.request()
	checkErr(err)

	var p fastjson.Parser
	v, err := p.ParseBytes(body)
	checkErr(err)
	if !v.Exists("result") || v.GetInt("result") != 0 {
		panic(fmt.Errorf("获取登陆帐号拥有的守护徽章列表失败，响应为 %s", string(body)))
	}

	medalList = new(MedalList)
	o := v.GetObject()
	o.Visit(func(k []byte, v *fastjson.Value) {
		switch string(k) {
		case "result":
		case "host-name":
		case "status":
		case "liveGiftConfig": // 忽略
		case "medalList":
			list := v.GetArray()
			medalList.MedalList = make([]Medal, 0, len(list))
			for _, l := range list {
				medalList.MedalList = append(medalList.MedalList, *getMedalJSON(l, t.UserID))
			}
		case "medal":
			medalList.MedalDetail.Medal = *getMedalJSON(v, t.UserID)
		case "medalDegreeLimit":
			medalList.MedalDetail.MedalDegree = *getMedalDegreeJSON(v)
		case "rankIndex":
			medalList.MedalDetail.UserRank, err = strconv.Atoi(string(v.GetStringBytes()))
			checkErr(err)
		case "clubName":
			medalList.MedalDetail.ClubName = string(v.GetStringBytes())
		default:
			log.Printf("登陆帐号拥有的守护徽章列表里出现未处理的key和value：%s %s", string(k), string(v.MarshalTo([]byte{})))
		}
	})

	return medalList, nil
}

// 获取指定用户正在佩戴的守护徽章信息
func getUserMedal(uid int64) (medal *Medal, e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("getUserMedal() error: %w", err)
		}
	}()

	client := &httpClient{
		url:    fmt.Sprintf(userMedalURL, uid),
		method: "GET",
	}
	body, err := client.request()
	checkErr(err)

	var p fastjson.Parser
	v, err := p.ParseBytes(body)
	checkErr(err)
	if !v.Exists("result") || v.GetInt("result") != 0 {
		panic(fmt.Errorf("获取指定用户正在佩戴的守护徽章信息失败，响应为 %s", string(body)))
	}

	medal = new(Medal)
	o := v.GetObject("wearMedalInfo")
	o.Visit(func(k []byte, v *fastjson.Value) {
		switch string(k) {
		case "uperId":
			medal.UperID = v.GetInt64()
		case "level":
			medal.Level = v.GetInt()
		case "clubName":
			medal.ClubName = string(v.GetStringBytes())
		case "uperName":
			medal.UperName = string(v.GetStringBytes())
		case "uperHeadUrl":
			medal.UperAvatar = string(v.GetStringBytes())
		default:
			log.Printf("指定用户正在佩戴的守护徽章信息里出现未处理的key和value：%s %s", string(k), string(v.MarshalTo([]byte{})))
		}
	})
	medal.UserID = uid
	medal.WearMedal = true

	return medal, nil
}

// 从json里获取用户信息
func getUserProfileJSON(v *fastjson.Value) *UserProfile {
	profile := new(UserProfile)
	o := v.GetObject()
	o.Visit(func(k []byte, v *fastjson.Value) {
		switch string(k) {
		case "name":
			profile.Nickname = string(v.GetStringBytes())
		case "headUrl":
			profile.Avatar = string(v.GetStringBytes())
		case "avatarFrameMobileImg":
			profile.AvatarFrame = string(v.GetStringBytes())
		case "followingCountValue":
			profile.FollowingCount = v.GetInt()
		case "fanCountValue":
			profile.FansCount = v.GetInt()
		case "contributeCountValue":
			profile.ContributeCount = v.GetInt()
		case "signature":
			profile.Signature = string(v.GetStringBytes())
		case "verifiedText":
			profile.VerifiedText = string(v.GetStringBytes())
		case "isJoinUpCollege":
			profile.IsJoinUpCollege = v.GetBool()
		case "isFollowing":
			profile.IsFollowing = v.GetBool()
		}
	})

	return profile
}

// 从json里获取用户直播信息
func getUserLiveInfoJSON(v *fastjson.Value) *UserLiveInfo {
	info := new(UserLiveInfo)
	o := v.GetObject()
	o.Visit(func(k []byte, v *fastjson.Value) {
		switch string(k) {
		case "user":
			info.Profile = *getUserProfileJSON(v)
		case "authorId":
		case "type":
			info.LiveType = *getLiveType(v)
		case "liveId":
			info.LiveID = string(v.GetStringBytes())
		case "streamName":
			info.StreamName = string(v.GetStringBytes())
		case "title":
			info.Title = string(v.GetStringBytes())
		case "createTime":
			info.LiveStartTime = v.GetInt64()
		case "portrait":
			info.Portrait = v.GetBool()
		case "panoramic":
			info.Panoramic = v.GetBool()
		case "coverUrls":
			info.LiveCover = string(v.GetStringBytes("0"))
		case "onlineCount":
			info.OnlineCount = v.GetInt()
		case "likeCount":
			info.LikeCount = v.GetInt()
		case "hasFansClub":
			info.HasFansClub = v.GetBool()
		case "disableDanmakuShow":
			info.DisableDanmakuShow = v.GetBool()
		case "paidShowUserBuyStatus":
			info.PaidShowUserBuyStatus = v.GetBool()
		}
	})
	info.Profile.UserID = v.GetInt64("authorId")

	return info
}

// 获取指定用户的直播信息
func getUserLiveInfo(uid int64, cookies Cookies) (info *UserLiveInfo, e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("getUserLiveInfo() error: %w", err)
		}
	}()

	client := &httpClient{
		url:     fmt.Sprintf(liveInfoURL, uid),
		method:  "GET",
		cookies: cookies,
	}
	body, err := client.request()
	checkErr(err)

	var p fastjson.Parser
	v, err := p.ParseBytes(body)
	checkErr(err)
	if !v.Exists("result") || v.GetInt("result") != 0 {
		panic(fmt.Errorf("获取指定用户的直播信息失败，响应为 %s", string(body)))
	}

	return getUserLiveInfoJSON(v), nil
}

// 获取指定主播的守护榜
func getMedalRankList(uid int64, cookies Cookies) (medalRankList *MedalRankList, e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("getMedalRankList() error: %w", err)
		}
	}()

	client := &httpClient{
		url:     fmt.Sprintf(medalRankURL, uid),
		method:  "GET",
		cookies: cookies,
	}
	body, err := client.request()
	checkErr(err)

	var p fastjson.Parser
	v, err := p.ParseBytes(body)
	checkErr(err)
	if !v.Exists("result") || v.GetInt("result") != 0 {
		panic(fmt.Errorf("获取指定主播的守护榜失败，响应为 %s", string(body)))
	}

	medalRankList = new(MedalRankList)
	o := v.GetObject()
	o.Visit(func(k []byte, v *fastjson.Value) {
		switch string(k) {
		case "result":
		case "host-name":
		case "friendshipDegreeRank":
			list := v.GetArray()
			medalRankList.RankList = make([]UserMedalInfo, len(list))
			for i, l := range list {
				o := l.GetObject()
				o.Visit(func(k []byte, v *fastjson.Value) {
					switch string(k) {
					case "userId":
					case "userInfo":
						medalRankList.RankList[i].UserProfile = *getUserProfileJSON(v)
					case "friendshipDegree":
						medalRankList.RankList[i].FriendshipDegree = v.GetInt()
					case "medalLevel":
						medalRankList.RankList[i].Level = v.GetInt()
					default:
						log.Printf("用户的守护徽章信息里出现未处理的key和value：%s %s", string(k), string(v.MarshalTo([]byte{})))
					}
				})
				medalRankList.RankList[i].UserID = l.GetInt64("userId")
			}
		case "clubName":
			medalRankList.ClubName = string(v.GetStringBytes())
		case "fansTotalCount":
			medalRankList.MedalCount = v.GetInt()
		case "isInFansClub":
			medalRankList.HasMedal = v.GetBool()
		case "curUserFriendshipDegree":
			medalRankList.UserFriendshipDegree = v.GetInt()
		case "curUserRankIndex":
			medalRankList.UserRank, err = strconv.Atoi(string(v.GetStringBytes()))
			checkErr(err)
		default:
			log.Printf("指定主播的守护榜里出现未处理的key和value：%s %s", string(k), string(v.MarshalTo([]byte{})))
		}
	})

	return medalRankList, nil
}

// 获取正在直播的直播间列表
func getLiveList(count int, page int, cookies Cookies) (liveList []UserLiveInfo, e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("getLiveList() error: %w", err)
		}
	}()

	form := fasthttp.AcquireArgs()
	defer fasthttp.ReleaseArgs(form)
	form.Set("count", strconv.Itoa(count))
	form.Set("pcursor", strconv.Itoa(page))
	client := &httpClient{
		url:         liveListURL,
		body:        form.QueryString(),
		method:      "POST",
		cookies:     cookies,
		contentType: formContentType,
	}
	body, err := client.request()
	checkErr(err)

	var p fastjson.Parser
	v, err := p.ParseBytes(body)
	checkErr(err)
	if !v.Exists("result") || v.GetInt("result") != 0 {
		panic(fmt.Errorf("获取正在直播的直播间列表失败，响应为：%s", string(body)))
	}

	list := v.GetArray("liveList")
	liveList = make([]UserLiveInfo, 0, len(list))
	for _, l := range list {
		liveList = append(liveList, *getUserLiveInfoJSON(l))
	}

	return liveList, nil
}

// 获取全部正在直播的直播间列表
func getAllLiveList(cookies Cookies) (liveList []UserLiveInfo, e error) {
	return getLiveList(1000000, 0, cookies)
}

// Distinguish 返回阿里云和腾讯云链接，目前阿里云的下载速度比较快
func (pb *Playback) Distinguish() (aliURL, txURL string) {
	switch {
	case strings.Contains(pb.URL, "alivod"):
		aliURL = pb.URL
	case strings.Contains(pb.URL, "txvod"):
		txURL = pb.URL
	default:
		log.Printf("未能识别的录播链接：%s", pb.URL)
	}

	switch {
	case strings.Contains(pb.BackupURL, "alivod"):
		aliURL = pb.BackupURL
	case strings.Contains(pb.BackupURL, "txvod"):
		txURL = pb.BackupURL
	default:
		log.Printf("未能识别的录播链接：%s", pb.BackupURL)
	}

	return aliURL, txURL
}

// GetWatchingList 返回直播间排名前50的在线观众信息列表，需要设置主播uid
func (ac *AcFunLive) GetWatchingList() ([]WatchingUser, error) {
	return ac.t.getWatchingList(ac.t.liveID)
}

// GetWatchingListWithLiveID 返回直播间排名前50的在线观众信息列表
func (ac *AcFunLive) GetWatchingListWithLiveID(liveID string) ([]WatchingUser, error) {
	return ac.t.getWatchingList(liveID)
}

// GetBillboard 返回指定uid的主播最近七日内的礼物贡献榜前50名观众的详细信息
func (ac *AcFunLive) GetBillboard(uid int64) ([]BillboardUser, error) {
	return ac.t.getBillboard(uid)
}

// GetSummary 返回直播总结信息，需要设置主播uid
func (ac *AcFunLive) GetSummary() (*Summary, error) {
	return ac.t.getSummary(ac.t.liveID)
}

// GetSummaryWithLiveID 返回直播总结信息
func (ac *AcFunLive) GetSummaryWithLiveID(liveID string) (*Summary, error) {
	return ac.t.getSummary(liveID)
}

// GetLuckList 返回抢到红包的用户列表，需要登陆AcFun帐号
func (ac *AcFunLive) GetLuckList(liveID, redpackID string) ([]LuckyUser, error) {
	return ac.t.getLuckList(liveID, redpackID)
}

// GetPlayback 返回直播回放的相关信息，目前部分直播没有回放
func (ac *AcFunLive) GetPlayback(liveID string) (*Playback, error) {
	return ac.t.getPlayback(liveID)
}

// GetGiftList 返回指定主播直播间的礼物数据，需要设置主播uid
func (ac *AcFunLive) GetGiftList() map[int64]GiftDetail {
	gifts := make(map[int64]GiftDetail)
	for k, v := range ac.t.gifts {
		gifts[k] = v
	}
	return gifts
}

// GetAllGiftList 返回全部礼物的数据
func (ac *AcFunLive) GetAllGiftList() (map[int64]GiftDetail, error) {
	return ac.t.getAllGift()
}

// GetWalletBalance 返回钱包里AC币和拥有的香蕉的数量，需要登陆AcFun帐号
func (ac *AcFunLive) GetWalletBalance() (accoins int, bananas int, e error) {
	return ac.t.getWalletBalance()
}

// GetKickHistory 返回主播踢人的历史记录，需要登陆主播的AcFun帐号，未测试
func (ac *AcFunLive) GetKickHistory() (e error) {
	return ac.t.getKickHistory()
}

// GetManagerList 返回主播的房管列表，需要登陆主播的AcFun帐号
func (ac *AcFunLive) GetManagerList() ([]Manager, error) {
	return ac.t.getManagerList()
}

// GetMedalDetail 返回登陆帐号拥有的指定主播的守护徽章详细信息，需要登陆AcFun帐号
func (ac *AcFunLive) GetMedalDetail(uid int64) (*MedalDetail, error) {
	return ac.t.getMedalDetail(uid)
}

// GetMedalList 返回登陆用户拥有的守护徽章列表，uid为想要获取守护徽章详细信息的主播的uid，可以为0，可用于获取指定主播的守护徽章名字，需要登陆AcFun帐号
func (ac *AcFunLive) GetMedalList(uid int64) (*MedalList, error) {
	return ac.t.getMedalList(uid)
}

// GetUserLiveInfo 返回uid指定用户的直播信息
func (ac *AcFunLive) GetUserLiveInfo(uid int64) (*UserLiveInfo, error) {
	return getUserLiveInfo(uid, ac.t.Cookies)
}

// GetMedalRankList 返回uid指定主播的守护榜（守护徽章亲密度排名前50名的用户），可用于获取指定主播的守护徽章名字
func (ac *AcFunLive) GetMedalRankList(uid int64) (medalRankList *MedalRankList, e error) {
	return getMedalRankList(uid, ac.t.Cookies)
}

// GetLiveList 返回正在直播的直播间列表，count为每页的直播间数量，page为第几页（从0开始数起）
func (ac *AcFunLive) GetLiveList(count int, page int) ([]UserLiveInfo, error) {
	return getLiveList(count, page, ac.t.Cookies)
}

// GetAllLiveList 返回全部正在直播的直播间列表
func (ac *AcFunLive) GetAllLiveList() ([]UserLiveInfo, error) {
	return getAllLiveList(ac.t.Cookies)
}

// GetUserMedal 返回uid指定用户正在佩戴的守护徽章信息，没有FriendshipDegree、JoinClubTime和CurrentDegreeLimit
func GetUserMedal(uid int64) (medal *Medal, e error) {
	return getUserMedal(uid)
}

// GetUserLiveInfo 返回uid指定用户的直播信息
func GetUserLiveInfo(uid int64) (*UserLiveInfo, error) {
	return getUserLiveInfo(uid, nil)
}

// GetMedalRankList 返回uid指定主播的守护榜（守护徽章亲密度排名前50名的用户），可用于获取指定主播的守护徽章名字
func GetMedalRankList(uid int64) (medalRankList *MedalRankList, e error) {
	return getMedalRankList(uid, nil)
}

// GetLiveList 返回正在直播的直播间列表，count为每页的直播间数量，page为第几页（从0开始数起）
func GetLiveList(count int, page int) ([]UserLiveInfo, error) {
	return getLiveList(count, page, nil)
}

// GetAllLiveList 返回全部正在直播的直播间列表
func GetAllLiveList() ([]UserLiveInfo, error) {
	return getAllLiveList(nil)
}
