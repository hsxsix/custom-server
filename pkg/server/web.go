/**
 * @File: web
 * @Author: hsien
 * @Description:
 * @Date: 9/17/21 2:28 PM
 */

package server

import (
	"context"
	"custom_server/pkg/config"
	"custom_server/pkg/database"
	"custom_server/pkg/log"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var once sync.Once

type GinApp struct {
	name    string
	db      map[string]*gorm.DB
	engine  *gin.Engine
	cronSrv *CronSrv
	config  *config.Config
}

func NewWebServer(cfgPath string, debug bool) *GinApp {
	app := &GinApp{
		db:     make(map[string]*gorm.DB),
		config: config.DefaultConfig(),
	}

	app.loadConfig(cfgPath)
	if !app.config.Server.Debug || !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	logLevel := app.config.Log.LogLevel
	if debug {
		logLevel = "DEBUG"
	}
	log.WithOption(log.SetLevel(
		logLevel),
		log.ColorPrint(true),
		log.FileName(app.config.Log.LogFile))

	app.engine = gin.New()
	return app
}

func (g *GinApp) loadConfig(path string) {
	cfg, err := config.LoadFromFile(path)

	if err != nil {
		log.Warn("load config failed", log.NameError("error", err))
		return
	}
	g.config = cfg

	if g.config != nil && len(g.config.DataBase) != 0 {
		for _, v := range g.config.DataBase {
			v.LogLevel = g.config.Log.LogLevel
		}
	}
}

func (g *GinApp) WithName(name string) *GinApp {
	g.name = name
	return g
}

func (g *GinApp) WithDB(name string) *GinApp {
	if g.config == nil || len(g.config.DataBase) == 0 {
		log.Warn("database initialized failed, config is empty, please check the configuration")
		return g
	}

	cfg, ok := g.config.DataBase[name]
	if !ok {
		log.Warnf("database initialized failed, not found database config with name %s", name)
		return g
	}

	db, err := database.InitDB(cfg)
	if err != nil {
		log.Warnf("database initialized failed, type:%s, error: %v", cfg.Type, err)
		return g
	}
	g.db[name] = db
	log.Infof("success initialized database with name: %s, type: %s", name, cfg.Type)

	return g
}

func (g *GinApp) WithCron() *GinApp {
	once.Do(func() {
		g.cronSrv = NewCron()
		log.Info("success initialized cron")
	})
	return g
}

func (g *GinApp) AddCronJob(job CronJob) {
	if g.cronSrv == nil {
		log.Warn("cron not initialize")
		return
	}

	if err := g.cronSrv.AddJob(job); err != nil {
		log.Warn("add task error", log.String("name", job.Name()), log.NameError("error", err))
		return
	}

	log.Debug("add cron job", log.String("name", job.Name()))
}

func (g *GinApp) AddCronFunc(name, spec string, cmd func()) {
	if g.cronSrv == nil {
		log.Warn("cron not initialize")
		return
	}

	if err := g.cronSrv.AddFunc(name, spec, cmd); err != nil {
		log.Warn("add task error", log.String("name", name), log.NameError("error", err))
		return
	}

	log.Debug("add cron job", log.String("name", name))
}

func (g *GinApp) CronEntries() []Entry {
	if g.cronSrv == nil {
		log.Warn("cron not initialize")
		return nil
	}
	return g.cronSrv.Entries()
}

type routerFunc func(app *GinApp)

func (g *GinApp) RegisterRouter(hfs ...routerFunc) *GinApp {
	for _, hf := range hfs {
		hf(g)
	}
	return g
}

func (g *GinApp) DB(name string) *gorm.DB {
	if len(g.db) == 0 {
		log.Panic("not found database instance, gorm db maybe not initialize")
	}
	db, ok := g.db[name]
	if !ok {
		log.Panicf("not found database instance with name: %s", name)
	}
	return db
}

func (g *GinApp) Engine() *gin.Engine {
	return g.engine
}

func (g *GinApp) Config() *config.Config {
	return g.config
}

func (g *GinApp) DebugMode() bool {
	return g.config.Server.Debug
}

func (g *GinApp) Cron() *CronSrv {
	if g.cronSrv == nil {
		log.Panic("cron not initialize")
	}
	return g.cronSrv
}

func (g *GinApp) Run() {
	if g.cronSrv != nil {
		g.cronSrv.Start()
	}

	srv := &http.Server{
		Addr:    g.config.Server.Addr,
		Handler: g.engine,
	}

	go func() {
		log.Infof("%s start listening %s...", g.name, g.config.Server.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen %s error: %v", g.config.Server.Addr, err)
		}
	}()

	// wait os signal to shut down server
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	if g.cronSrv != nil {
		g.cronSrv.Stop()
	}

	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown:", log.NameError("error", err))
	}

	log.Info("server exiting!")
}
