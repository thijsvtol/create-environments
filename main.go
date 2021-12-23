package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-github/v41/github"
	"golang.org/x/oauth2"
)

type env struct {
	repo                  string
	repoOwner             string
	token                 string
	environments          []string
	waitTime              int
	requiredReviewers     []string
	protectedBranchesOnly bool
	customBranches        bool
}

type service struct {
	ctx    context.Context
	client *github.Client
	env    *env
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
	protectedBranchesOnly, err := strconv.ParseBool(os.Getenv("INPUT_PROTECTED_BRANCHES_ONLY"))
	if err != nil {
		log.Fatalln("protected_branches_only is not a boolean")
	}

	e := &env{
		repoOwner:             repo[0],
		repo:                  repo[1],
		token:                 os.Getenv("INPUT_TOKEN"),
		environments:          environments,
		waitTime:              waitTime,
		requiredReviewers:     requiredReviewers,
		protectedBranchesOnly: protectedBranchesOnly,
		customBranches:        false,
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
	var createdRepoEnvironments []*github.Environment
	for _, env := range s.env.environments {
		opt := &github.CreateUpdateEnvironment{
			WaitTimer: &s.env.waitTime,
			Reviewers: s.getUsers(),
			DeploymentBranchPolicy: &github.BranchPolicy{
				ProtectedBranches:    &s.env.protectedBranchesOnly,
				CustomBranchPolicies: &s.env.customBranches,
			},
		}

		environments, _, err := s.client.Repositories.CreateUpdateEnvironment(s.ctx, s.env.repoOwner, s.env.repo, env, opt)
		if err != nil {
			fmt.Sprintf("Options: %v", opt)
			log.Fatalln(err)
			return nil, err
		}

		createdRepoEnvironments = append(createdRepoEnvironments, environments)
	}

	for _, env := range createdRepoEnvironments {
		log.Printf("Created environment [%v] %v", *env.Name, *env.URL)
	}

	return createdRepoEnvironments, nil
}

func (s *service) getUsers() []*github.EnvReviewers {
	var retrievedUsers []*github.EnvReviewers
	for _, user := range s.env.requiredReviewers {
		if strings.Contains(user, "/") {
			orgTeam := strings.Split(user, "/")
			team, _, err := s.client.Teams.GetTeamBySlug(s.ctx, orgTeam[0], orgTeam[1])
			if err != nil {
				log.Fatalln(err)
				return nil
			}

			t := &github.EnvReviewers{
				Type: team.Organization.Type,
				ID:   team.ID,
			}
			retrievedUsers = append(retrievedUsers, t)
		} else if user != "" {
			user, _, err := s.client.Users.Get(s.ctx, user)
			if err != nil {
				log.Fatalln(err)
				return nil
			}

			u := &github.EnvReviewers{
				Type: user.Type,
				ID:   user.ID,
			}
			retrievedUsers = append(retrievedUsers, u)
		}
	}

	return retrievedUsers
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
