package beego

import (
	"encoding/json"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/astaxie/beego/session"
)

// MIME (Multipurpose Internet Mail Extensions) 是描述消息内容类型的因特网标准
// 可见 http://www.w3school.com.cn/media/media_mimeref.asp
// 这里使用标准库的 mime.AddExtensionType()将扩展名和mimetype关联
// mime的定义见 mime.go
func registerMime() error {
	for k, v := range mimemaps {
		mime.AddExtensionType(k, v)
	}
	return nil
}
// 设置不同错误的默认回调方法
func registerDefaultErrorHandler() error {
	m := map[string]func(http.ResponseWriter, *http.Request){
		"401": unauthorized,
		"402": paymentRequired,
		"403": forbidden,
		"404": notFound,
		"405": methodNotAllowed,
		"500": internalServerError,
		"501": notImplemented,
		"502": badGateway,
		"503": serviceUnavailable,
		"504": gatewayTimeout,
	}
	for e, h := range m {
		if _, ok := ErrorMaps[e]; !ok {
			ErrorHandler(e, h)
		}
	}
	return nil
}
//判断是否需要初始化GlobalSessions(Session管理器)
func registerSession() error {
	// BConfig(位于 config.go内的全局变量)
	if BConfig.WebConfig.Session.SessionOn {
		var err error
		sessionConfig := AppConfig.String("sessionConfig")
		if sessionConfig == "" {
			// 启用默认的 Session配置
			conf := map[string]interface{}{
				"cookieName":      BConfig.WebConfig.Session.SessionName,
				"gclifetime":      BConfig.WebConfig.Session.SessionGCMaxLifetime,
				"providerConfig":  filepath.ToSlash(BConfig.WebConfig.Session.SessionProviderConfig),
				"secure":          BConfig.Listen.EnableHTTPS,
				"enableSetCookie": BConfig.WebConfig.Session.SessionAutoSetCookie,
				"domain":          BConfig.WebConfig.Session.SessionDomain,
				"cookieLifeTime":  BConfig.WebConfig.Session.SessionCookieLifeTime,
			}
			confBytes, err := json.Marshal(conf)
			if err != nil {
				return err
			}
			sessionConfig = string(confBytes)
		}
		if GlobalSessions, err = session.NewManager(BConfig.WebConfig.Session.SessionProvider, sessionConfig); err != nil {
			return err
		}
		//开启一个goroutine来处理session的回收,定义于 session.session.go:227
		go GlobalSessions.GC()
	}
	return nil
}
// 构建模板
func registerTemplate() error {
	if err := BuildTemplate(BConfig.WebConfig.ViewsPath); err != nil {
		if BConfig.RunMode == DEV {
			Warn(err)
		}
		return err
	}
	return nil
}
//判断是否需要加入文档的路由
func registerDocs() error {
	if BConfig.WebConfig.EnableDocs {
		Get("/docs", serverDocs)
		Get("/docs/*", serverDocs)
	}
	return nil
}
// 判断是否需要启动进程内监控服务器
func registerAdmin() error {
	if BConfig.Listen.EnableAdmin {
		go beeAdminApp.Run()
	}
	return nil
}
