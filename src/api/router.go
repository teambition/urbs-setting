package api

import (
	"github.com/teambition/gear"
	tracing "github.com/teambition/gear-tracing"

	"github.com/teambition/urbs-setting/src/bll"
	"github.com/teambition/urbs-setting/src/middleware"
	"github.com/teambition/urbs-setting/src/util"
)

func init() {
	util.DigProvide(newAPIs)
	util.DigProvide(newRouters)
}

// APIs ..
type APIs struct {
	Healthz *Healthz
	User    *User
	Group   *Group
	Product *Product
	Module  *Module
	Setting *Setting
	Label   *Label
}

func newAPIs(blls *bll.Blls) *APIs {
	return &APIs{
		Healthz: &Healthz{blls: blls},
		User:    &User{blls: blls},
		Group:   &Group{blls: blls},
		Product: &Product{blls: blls},
		Module:  &Module{blls: blls},
		Setting: &Setting{blls: blls},
		Label:   &Label{blls: blls},
	}
}

func newRouters(apis *APIs) []*gear.Router {

	router := gear.NewRouter()
	// health check
	router.Get("/healthz", apis.Healthz.Get)
	// 读取指定用户的灰度标签，包括继承自群组的标签，返回轻量级 labels，无身份验证，用于网关
	router.Get("/users/:uid/labels:cache", tracing.New(), apis.User.ListLablesInCache)

	routerV1 := gear.NewRouter(gear.RouterOptions{
		Root: "/v1",
	})
	routerV1.Use(middleware.Auth)
	routerV1.Use(tracing.New())

	// ***** user ******
	// 读取指定用户的灰度标签，支持条件筛选
	routerV1.Get("/users/:uid/labels", apis.User.ListLables)
	// 读取指定用户的功能配置项，支持条件筛选
	routerV1.Get("/users/:uid/settings", apis.User.ListSettings)
	// 查询指定用户是否存在
	routerV1.Get("/users/:uid+:exists", apis.User.CheckExists)
	// 批量添加用户
	routerV1.Post("/users:batch", apis.User.BatchAdd)
	// 删除指定用户的指定灰度标签
	routerV1.Delete("/users/:uid/labels/:hid", apis.User.RemoveLable)
	// 更新指定用户的指定配置项
	routerV1.Put("/users/:uid/settings/:hid", apis.User.UpdateSetting)
	// 删除指定用户的指定配置项
	routerV1.Delete("/users/:uid/settings/:hid", apis.User.RemoveSetting)

	// ***** group ******
	// 读取指定群组的灰度标签，支持条件筛选
	routerV1.Get("/groups/:uid/labels", apis.Group.ListLables)
	// 读取指定群组的功能配置项，支持条件筛选
	routerV1.Get("/groups/:uid/settings", apis.Group.ListSettings)
	// 读取群组列表，支持条件筛选
	routerV1.Get("/groups", apis.Group.List)
	// 查询指定群组是否存在
	routerV1.Get("/groups/:uid+:exists", apis.Group.CheckExists)
	// 批量添加群组
	routerV1.Post("/groups:batch", apis.Group.BatchAdd)
	// 更新指定群组
	routerV1.Put("/groups/:uid", apis.Group.Update)
	// 删除指定群组
	routerV1.Delete("/groups/:uid", apis.Group.Delete)
	// 读取群组成员列表，支持条件筛选
	routerV1.Get("/groups/:uid/members", apis.Group.ListMembers)
	// 指定群组批量添加成员
	routerV1.Post("/groups/:uid/members:batch", apis.Group.BatchAddMembers)
	// 指定群组根据条件清理成员
	routerV1.Delete("/groups/:uid/members", apis.Group.RemoveMembers)
	// 删除指定群组的指定灰度标签
	routerV1.Delete("/groups/:uid/labels/:hid", apis.Group.RemoveLable)
	// 更新指定用户的指定配置项
	routerV1.Put("/groups/:uid/settings/:hid", apis.Group.UpdateSetting)
	// 删除指定群组的指定配置项
	routerV1.Delete("/groups/:uid/settings/:hid", apis.Group.RemoveSetting)

	// ***** product ******
	// 读取产品列表，支持条件筛选
	routerV1.Get("/products", apis.Product.List)
	// 创建产品
	routerV1.Post("/products", apis.Product.Create)
	// 更新指定产品
	routerV1.Put("/products/:product", apis.Product.Update)
	// 下线指定产品功能模块
	routerV1.Put("/products/:product+:offline", apis.Product.Offline)
	// 重新上线指定产品功能模块
	routerV1.Put("/products/:product+:online", apis.Product.Online)
	// 删除指定产品
	routerV1.Delete("/products/:product", apis.Product.Delete)

	// ***** module ******
	// 读取指定产品的功能模块
	routerV1.Get("/products/:product/modules", apis.Module.List)
	// 指定产品创建功能模块
	routerV1.Post("/products/:product/modules", apis.Module.Create)
	// 更新指定产品功能模块
	routerV1.Put("/products/:product/modules/:module", apis.Module.Update)
	// 下线指定产品功能模块
	routerV1.Put("/products/:product/modules/:module+:offline", apis.Module.Offline)
	// 重新上线指定产品功能模块
	routerV1.Put("/products/:product/modules/:module+:online", apis.Module.Online)

	// ***** setting ******
	// 读取指定产品功能模块的配置项
	routerV1.Get("/products/:product/modules/:module/settings", apis.Setting.List)
	// 创建指定产品功能模块配置项
	routerV1.Post("/products/:product/modules/:module/settings", apis.Setting.Create)
	// 更新指定产品功能模块配置项
	routerV1.Put("/products/:product/modules/:module/settings/:setting", apis.Setting.Update)
	// 下线指定产品功能模块配置项
	routerV1.Put("/products/:product/modules/:module/settings/:setting+:offline", apis.Setting.Offline)
	// 重新上线指定产品功能模块配置项
	routerV1.Put("/products/:product/modules/:module/settings/:setting+:online", apis.Setting.Online)
	// 批量为用户或群组设置产品功能模块配置项
	routerV1.Post("/products/:product/modules/:module/settings/:setting+:assign", apis.Setting.Assign)

	// ***** label ******
	// 读取指定产品灰度标签
	routerV1.Get("/products/:product/labels", apis.Label.List)
	// 创建指定产品灰度标签
	routerV1.Post("/products/:product/labels", apis.Label.Create)
	// 更新指定产品灰度标签
	routerV1.Put("/products/:product/labels/:label", apis.Label.Update)
	// 更新指定产品灰度标签
	routerV1.Delete("/products/:product/labels/:label", apis.Label.Delete)
	// 下线指定产品灰度标签
	routerV1.Put("/products/:product/labels/:label+:offline", apis.Label.Offline)
	// 重新上线指定产品灰度标签
	routerV1.Put("/products/:product/labels/:label+:online", apis.Label.Online)
	// 批量为用户或群组设置产品灰度标签
	routerV1.Post("/products/:product/labels/:label+:assign", apis.Label.Assign)

	return []*gear.Router{routerV1, router}
}
