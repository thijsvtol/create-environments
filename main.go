package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/google/go-github/v41/github"
	"golang.org/x/oauth2"
)

type env struct {
	repo              string
	repoOwner         string
	token             string
	environments      []string
	waitTime          int
	requiredReviewers []string
}

type service struct {
	ctx    context.Context
	client *github.Client
	env    *env
	wg     sync.WaitGroup
}

func environment() *env {
	repo := strings.Split(os.Getenv("INPUT_REPO"), "/")
	environments := strings.Split(os.Getenv("INPUT_ENVIRONMENTS"), ",")
	if len(environments) == 0 {
		log.Fatalln("The environments variable is required and could have multiple values separated by comma")
	}
	waitTime, err := strconv.Atoi(os.Getenv("INPUT_WAIT_TIME"))
	if err != nil {
		log.Fatalln("wait_time is not a number")
	}
	requiredReviewers := strings.Split(os.Getenv("INPUT_REQUIRED_REVIEWERS"), ",")
	if len(environments) == 0 {
		log.Fatalln("The required_reviewers variable is required and could have multiple values separated by comma")
	}

	e := &env{
		repoOwner:         repo[0],
		repo:              repo[1],
		token:             os.Getenv("INPUT_TOKEN"),
		environments:      environments,
		waitTime:          waitTime,
		requiredReviewers: requiredReviewers,
	}
	return e
}

func (e *env) debugPrint() {
	log.Printf("Repo: %v", e.repo)
	log.Printf("Repo Owner: %v", e.repoOwner)
	log.Printf("Token: %v", e.token)
	log.Printf("Environments: %v", e.environments)
	log.Printf("Wait time: %v", e.waitTime)
	log.Printf("Required reviewers: %v", e.requiredReviewers)
}

func (s *service) createUpdateEnvironments() ([]*github.Environment, error) {
	var repoEnvironments []*github.Environment

	for _, env := range s.env.environments {
		opt := &github.CreateUpdateEnvironment{
			WaitTimer: &s.env.waitTime,
		}

		environments, _, err := s.client.Repositories.CreateUpdateEnvironment(s.ctx, s.env.repoOwner, s.env.repo, env, opt)
		if err != nil {
			log.Fatalln(err)
			return nil, err
		}

		repoEnvironments = append(repoEnvironments, environments)
	}
	log.Printf("Created/updated the following environments: ", repoEnvironments)

	return repoEnvironments, nil
}

func main() {
	log.SetOutput(os.Stdout)

	env := environment()
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: env.token},
	)
	tc := oauth2.NewClient(ctx, ts)

	svc := &service{
		ctx:    ctx,
		client: github.NewClient(tc),
		env:    env,
	}
	env.debugPrint()

	svc.createUpdateEnvironments()
}
