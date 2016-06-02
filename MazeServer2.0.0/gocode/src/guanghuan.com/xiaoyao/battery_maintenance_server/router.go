package main

const (
    MaintenanceItem = "/maintenance/:content"
    //MaintenanceItem = "/maintenance"
)

const (
    MaintenanceSubItem_Prop  = "prop"
    MaintenanceSubItem_CDkey = "CDkeyExchange"

    SDKCallBackSubItem_CallBack = "/sdk/callback"
)

const (
    MaintenanceSubject_PROP  = "maintenanceprop"
    MaintenanceSubject_CDkey = "maintenancecdkey"

    SDKCallBackSubject_CallBack = "sdkcallback"
)

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

//默认的路由表管理器
var (
    DefHttpPostTable        = CreateHttpPostToNatsRouteTable(0)
    DefHttpPostNoTokenTable = CreateHttpPostToNatsRouteTable(0)
)

func (rt HttpPostToNatsRouteTable) AddHttpPostRoute(tag string, uri string, subj string) {
    m := map[string]*HttpPostToNatsRoute(rt)
    m[tag] = NewHttpPostToNatsRoute(tag, uri, subj)
}

//增加消息路由信息
func AddHttpPostRoute(tag string, uri string, subj string) {
    DefHttpPostTable.AddHttpPostRoute(tag, uri, subj)
}

//初始化消息路由表
func InitRouteTable() {
    DefHttpPostTable.AddHttpPostRoute(MaintenanceSubItem_Prop, MaintenanceSubItem_Prop, MaintenanceSubject_PROP)
    DefHttpPostTable.AddHttpPostRoute(MaintenanceSubItem_CDkey, MaintenanceSubItem_CDkey, MaintenanceSubject_CDkey)

    DefHttpPostTable.AddHttpPostRoute(SDKCallBackSubItem_CallBack, SDKCallBackSubItem_CallBack, SDKCallBackSubject_CallBack)
}
