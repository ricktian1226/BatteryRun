package main

type RouteType int32

const (
    RouteTypeUnknown RouteType = 0
    RouteHttpPost    RouteType = 1
    RouteHttpGet     RouteType = 2
    RouteNats        RouteType = 3
)

type RouteEntry struct {
    Type RouteType
    Uri  string
}

type RoutePath struct {
    Tag      string
    Inbound  *RouteEntry
    Outbound *RouteEntry
}

func NewRoutePath(tag string, in_type RouteType, in_uri string, out_type RouteType, out_uri string) (r *RoutePath) {
    r = &RoutePath{
        Tag: tag,
        Inbound: &RouteEntry{
            Type: in_type,
            Uri:  in_uri,
        },
        Outbound: &RouteEntry{
            Type: out_type,
            Uri:  out_uri,
        },
    }
    return
}

type HttpPostToNatsRoute struct {
    RoutePath
}

func (r *HttpPostToNatsRoute) GetHttpUri() string {
    return r.Inbound.Uri
}
func (r *HttpPostToNatsRoute) GetNatsSubject() string {
    return r.Outbound.Uri
}

func NewHttpPostToNatsRoute(tag string, in_uri string, out_subj string) (r *HttpPostToNatsRoute) {
    r = &HttpPostToNatsRoute{
        RoutePath: *NewRoutePath(tag, RouteHttpPost, in_uri, RouteNats, out_subj),
    }
    return
}

type HttpGetToNatsRoute struct {
    RoutePath
}

func NewHttpGetToNatsRoute(tag string, in_uri string, out_subj string) (r *HttpGetToNatsRoute) {
    r = &HttpGetToNatsRoute{
        RoutePath: *NewRoutePath(tag, RouteHttpGet, in_uri, RouteNats, out_subj),
    }
    return
}

func (r *HttpGetToNatsRoute) GetHttpUri() string {
    return r.Inbound.Uri
}

func (r *HttpGetToNatsRoute) GetNatsSubject() string {
    return r.Outbound.Uri
}

type HttpPostToNatsRouteTable map[string]*HttpPostToNatsRoute

func CreateHttpPostToNatsRouteTable(min_cap int) (rt HttpPostToNatsRouteTable) {
    if min_cap <= 0 {
        min_cap = 10
    }
    rt = make(HttpPostToNatsRouteTable, min_cap)
    return
}

func (rt HttpPostToNatsRouteTable) GetRoutePath(tag string) (r *HttpPostToNatsRoute) {
    m := map[string]*HttpPostToNatsRoute(rt)
    return m[tag]
}

var (
    DefHttpPostTable        HttpPostToNatsRouteTable = CreateHttpPostToNatsRouteTable(0)
    DefHttpPostNoTokenTable HttpPostToNatsRouteTable = CreateHttpPostToNatsRouteTable(0)
)

func (rt HttpPostToNatsRouteTable) AddHttpPostRoute(tag string, uri string, subj string) {
    m := map[string]*HttpPostToNatsRoute(rt)
    m[tag] = NewHttpPostToNatsRoute(tag, uri, subj)
}

func AddHttpPostRoute(tag string, uri string, subj string) {
    DefHttpPostTable.AddHttpPostRoute(tag, uri, subj)
}

func InitRouteTable() {
    DefHttpPostTable.AddHttpPostRoute("/v1/login", "/v1/login/:token", "login")
    DefHttpPostTable.AddHttpPostRoute("/v1/user", "/v1/user/:token", "userdata")
    DefHttpPostTable.AddHttpPostRoute("/v1/friend", "/v1/friend/:token", "frienddata")
    DefHttpPostTable.AddHttpPostRoute("/v1/newgame", "/v1/newgame/:token", "newgame")
    DefHttpPostTable.AddHttpPostRoute("/v1/gameresult", "/v1/gameresult/:token", "gameresult")
    DefHttpPostTable.AddHttpPostRoute("/v2/gameresult", "/v2/gameresult/:token", "gameresult2")

    //DefHttpPostTable.AddHttpPostRoute("/v1/stamina", "/v1/stamina/:token", "stamina")
    DefHttpPostTable.AddHttpPostRoute("/v1/gift/query", "/v1/gift/query/:token", "gift_query")
    DefHttpPostTable.AddHttpPostRoute("/v1/gift/op", "/v1/gift/op/:token", "gift_op")
    DefHttpPostTable.AddHttpPostRoute("/v1/goods/query", "/v1/goods/query/:token", "goods_query")
    DefHttpPostTable.AddHttpPostRoute("/v1/goods/buy", "/v1/goods/buy/:token", "goods_buy")
    DefHttpPostTable.AddHttpPostRoute("/v2/announcement", "/v2/announcement/:token", "announcement")

    DefHttpPostNoTokenTable.AddHttpPostRoute("/iap_verify/order_request", "/iap_verify/order_request", "iap_verify_order")
    DefHttpPostNoTokenTable.AddHttpPostRoute("/iap_verify/order_verify", "/iap_verify/order_verify", "iap_verify_validate")
    DefHttpPostNoTokenTable.AddHttpPostRoute("/v1/device/device_id", "/v1/device/device_id", "device_id")
    DefHttpPostNoTokenTable.AddHttpPostRoute("/v2/lotto/lotto_op", "/v2/lotto/lotto_op", "lotto_op")

    DefHttpPostNoTokenTable.AddHttpPostRoute("/v2/systemmail/systemmail_op", "/v2/systemmail/systemmail_op", "systemmail_op")
    DefHttpPostNoTokenTable.AddHttpPostRoute("/v2/friendmail/friendmail_op", "/v2/friendmail/friendmail_op", "friendmail_op")
    DefHttpPostNoTokenTable.AddHttpPostRoute("/v2/roleinfolist/roleinfolist_op", "/v2/roleinfolist/roleinfolist_op", "roleinfolist_op")
    DefHttpPostNoTokenTable.AddHttpPostRoute("/v2/jigsaw/jigsaw_op", "/v2/jigsaw/jigsaw_op", "jigsaw_op")
    DefHttpPostNoTokenTable.AddHttpPostRoute("/v2/rune/rune_op", "/v2/rune/rune_op", "rune_op")
    DefHttpPostNoTokenTable.AddHttpPostRoute("/v2/beforegameprop/beforegameprop_op", "/v2/beforegameprop/beforegameprop_op", "beforegameprop_op")
    DefHttpPostNoTokenTable.AddHttpPostRoute("/v2/signin/query", "/v2/signin/query", "signin_query")
    DefHttpPostNoTokenTable.AddHttpPostRoute("/v2/signin/sign", "/v2/signin/sign", "signin_sign")
    DefHttpPostNoTokenTable.AddHttpPostRoute("/v2/usermission/query", "/v2/usermission/query", "usermission_query")
    DefHttpPostNoTokenTable.AddHttpPostRoute("/v2/usermission/confirm", "/v2/usermission/confirm", "usermission_confirm")
    DefHttpPostNoTokenTable.AddHttpPostRoute("/v2/checkpoint/query_range", "/v2/checkpoint/query_range", "checkpoint_query_range")
    DefHttpPostNoTokenTable.AddHttpPostRoute("/v2/checkpoint/query_detail", "/v2/checkpoint/query_detail", "checkpoint_query_detail") //查询记忆点排行版（好友排行榜或者全局排行榜）
    // DefHttpPostNoTokenTable.AddHttpPostRoute("/v2/checkpoint/commit", "/v2/checkpoint/commit", "checkpoint_commit")
    DefHttpPostNoTokenTable.AddHttpPostRoute("/v2/checkpoint/unlock", "/v2/checkpoint/unlock", "checkpoint_unlock")
    DefHttpPostNoTokenTable.AddHttpPostRoute("/v2/wallet/query", "/v2/wallet/query", "wallet_query")
    DefHttpPostNoTokenTable.AddHttpPostRoute("/v2/res/res_op", "/v2/res/res_op", "res_op")
    DefHttpPostNoTokenTable.AddHttpPostRoute("/v2/prop/prop_res_query", "/v2/prop/prop_res_query", "prop_res_query")
    DefHttpPostNoTokenTable.AddHttpPostRoute("/v2/memcache/memcache_op", "/v2/memcache/memcache_op", "memcache_op")
    DefHttpPostNoTokenTable.AddHttpPostRoute("/v2/advertisement/advertisement", "/v2/advertisement/advertisement", "advertisement")
    DefHttpPostNoTokenTable.AddHttpPostRoute("/v2/bind/bind", "/v2/bind/bind", "bind")

    DefHttpPostNoTokenTable.AddHttpPostRoute("/v2/signin/query_new", "/v2/signin/query_new", "new_signin_query")
    DefHttpPostNoTokenTable.AddHttpPostRoute("/v2/maintenance/cdkey", "/v2/maintenance/cdkey", "cdkey")
    // 分享消息路由
    DefHttpPostNoTokenTable.AddHttpPostRoute("/v2/share/query", "/v2/share/query", "share_query")
    DefHttpPostNoTokenTable.AddHttpPostRoute("/v2/share/op", "/v2/share/op", "share_op")

    DefHttpPostNoTokenTable.AddHttpPostRoute("/v2/sdk/order_op", "/v2/sdk/order_op", "order_op")          // sdk 操作
    DefHttpPostNoTokenTable.AddHttpPostRoute("/v2/sdk/order_query", "/v2/sdk/order_query", "order_query") // 订单查询

    DefHttpPostNoTokenTable.AddHttpPostRoute("/v2/creatname", "/v2/creatname", "creatname") // 玩家起名
    DefHttpPostNoTokenTable.AddHttpPostRoute("/v2/ranklist", "/v2/ranklist", "ranklist")    //
}
