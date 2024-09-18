package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/slack-go/slack"
)

func getAllUserUIDs(ctx context.Context, client *slack.Client, pageSize int) ([]string, error) {
	var uids []string
	var err error

	pages := 0
	pager := client.GetUsersPaginated(slack.GetUsersOptionLimit(pageSize))
	for {
		// Note reassignment of pager to the value returned by Next()
		pager, err = pager.Next(ctx)
		if failedErr := pager.Failure(err); failedErr != nil {
			var rateLimited *slack.RateLimitedError
			if errors.As(failedErr, &rateLimited) && rateLimited.Retryable() {
				fmt.Println("Rate limited by Slack API; sleeping", rateLimited.RetryAfter)
				select {
				case <-ctx.Done():
					return uids, ctx.Err()
				case <-time.After(rateLimited.RetryAfter):
					continue
				}
			}
			return uids, fmt.Errorf("paginating users: %w", failedErr)
		}
		if pager.Done(err) {
			break
		}

		for _, user := range pager.Users {
			uids = append(uids, user.ID)
		}

		pages++
	}

	fmt.Printf("Pagination complete after %d pages\n", pages)

	return uids, nil
}

func main() {
	client := slack.New("YOUR_TOKEN_HERE")

	uids, err := getAllUserUIDs(context.Background(), client, 1000)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Collected %d UIDs\n", len(uids))
}
