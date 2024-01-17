package main

import (
	"context"
	"log"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
	"github.com/ishua/a3bot5/fsnotes/internal/clients/mygit"
	"github.com/ishua/a3bot5/fsnotes/internal/domain"
	"github.com/ishua/a3bot5/fsnotes/internal/repo"
	"github.com/redis/go-redis/v9"
)

type MyConfig struct {
	Redis           string `default:"redis:6379" env:"REDIS" usage:"connect str to redis"`
	SubChannel      string `default:"fsnotes" usage:"channel for subscribe jobs"`
	TbotChannel     string `default:"tbot" usage:"channel for jobs for tbot"`
	RepoPath        string `default:"data/fsnotes" usage:" path to repository fsnotes"`
	RepoUrl         string `required:"true"`
	RepoAccessToken string `env:"REPOACCESSTOKEN" required:"true"`
	RepoDiaryPath   string `required:"true"`
}

type Pubsub struct {
	*redis.Client
}

var (
	cfg    MyConfig
	rdb    *Pubsub
	myGit  *mygit.Repo
	myRepo *repo.GitFile
)

// init config
func init() {
	loader := aconfig.LoaderFor(&cfg, aconfig.Config{
		Files: []string{"conf/fsnote_config.yaml"},
		FileDecoders: map[string]aconfig.FileDecoder{
			".yaml": aconfigyaml.New(),
		},
	})
	if err := loader.Load(); err != nil {
		panic(err)
	}
}

// // init redis
func init() {
	rdb = &Pubsub{redis.NewClient(&redis.Options{
		Addr:     cfg.Redis,
		Password: "", // no password set
		DB:       0,  // use default DB
	})}
	log.Println("Redis: " + cfg.Redis)
}
func (r Pubsub) Pub(ctx context.Context, value string) {
	err := r.Publish(ctx, cfg.TbotChannel, value).Err()
	if err != nil {
		log.Println("publish error: " + err.Error())
	}
}

// init git
func init() {
	var err error
	myGit, err = mygit.NewClient(cfg.RepoPath, cfg.RepoUrl, cfg.RepoAccessToken)
	if err != nil {
		panic(err)
	}
}

func main() {

	ctx := context.Background()
	myRepo = repo.NewGitFile(cfg.RepoPath, cfg.RepoDiaryPath, myGit)
	d := domain.NewDiary(myRepo)
	m := domain.NewModel(d, rdb)

	pubsub := rdb.Subscribe(ctx, cfg.SubChannel)

	log.Println("Subscribe to channel: " + cfg.SubChannel)
	// Close the subscription when we are done.
	defer pubsub.Close()

	ch := pubsub.Channel()

	for msg := range ch {
		go m.DoJob(ctx, msg.Payload)
	}
}
