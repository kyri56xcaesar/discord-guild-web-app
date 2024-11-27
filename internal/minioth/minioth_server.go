package minioth

import (
	"bufio"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type MService struct {
	Engine  *gin.Engine
	Config  *MConfig
	Minioth *Minioth
}

type MConfig struct {
	ConfPath string
	Ip       string
	Port     string
}

const (
	DEFAULT_conf_name string = "minioth.conf"
	DEFAULT_conf_path string = "configs/"
)

func newConfig() MConfig {
	return MConfig{
		ConfPath: DEFAULT_conf_path + DEFAULT_conf_name,
		Ip:       "localhost",
		Port:     "8081",
	}
}

func NewMSerivce(m *Minioth) MService {
	cfg := newConfig()
	cfg.loadConfig("")

	srv := MService{
		Minioth: m,
		Engine:  gin.Default(),
		Config:  &cfg,
	}

	return srv
}

func (srv *MService) ServeHTTP() {
	minioth := srv.Minioth

	apiV1 := srv.Engine.Group("/v1")
	{
		apiV1.GET("/user", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"content": minioth.Select("users"),
			})
		})
		apiV1.POST("/user", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "activo",
			})
		})

		apiV1.GET("/group", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"content": minioth.Select("groups"),
			})
		})

		apiV1.POST("/group", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "activo",
			})
		})

	}
	server := &http.Server{
		Addr:              srv.Config.Ip + ":" + srv.Config.Port,
		Handler:           srv.Engine,
		ReadHeaderTimeout: time.Second * 5,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	<-ctx.Done()

	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}

func (cfg *MConfig) loadConfig(path string) error {
	var confFile *os.File
	var err error

	curpath, err := os.Getwd()
	cfg.ConfPath = curpath + path
	log.Printf("Current path: %v", curpath)
	confFile, err = os.Open(path)
	if err != nil {
		log.Printf("Given Path: %v, err: %v", path, err)
		confFile, err = os.Open(DEFAULT_conf_path + DEFAULT_conf_name)
		if err != nil {
			log.Print(err)
			return err
		}
		log.Print("Default config opened.")
	}
	defer confFile.Close()

	scanner := bufio.NewScanner(confFile)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "PORT":
			if value != "" {
				cfg.Port = value
			}
		case "IP":
			if value != "" {
				cfg.Ip = value
			}

		default:
		}
	}

	return nil
}
