package store

import (
	"fmt"
	"strings"

	"github.com/bebrws/goPR/internal/models"
	"github.com/sirupsen/logrus"
)

// CompareStates compares the old and new GHState and returns a list of concise strings describing the changes.
func CompareStates(oldState, newState models.GHState) []string {
	var changes []string

	for _, newRepo := range newState.RepoStates {
		oldRepo := findRepoByName(oldState.RepoStates, newRepo.Name)
		if oldRepo == nil {
			changes = append(changes, fmt.Sprintf("New repository: %s", newRepo.Name))
			continue
		}

		for _, newPR := range newRepo.PRs {
			oldPR := findPRByNumber(oldRepo.PRs, newPR.Number)
			if oldPR == nil {
				changes = append(changes, fmt.Sprintf("New PR: %s PR#%d", newRepo.Name, newPR.Number))
				continue
			}

			if newPR.Body != oldPR.Body {
				changes = append(changes, fmt.Sprintf("PR Body Change: %s PR#%d", newRepo.Name, newPR.Number))
			}

			for _, newReview := range newPR.Reviews {
				oldReview := findReviewByID(oldPR.Reviews, newReview.ID)
				if oldReview == nil {
					changes = append(changes, fmt.Sprintf("New Review: %s PR#%d Reviewer: %s", newRepo.Name, newPR.Number, newReview.Login))
					continue
				}

				if newReview.Body != oldReview.Body {
					changes = append(changes, fmt.Sprintf("Review Body Change: %s PR#%d Reviewer: %s", newRepo.Name, newPR.Number, newReview.Login))
				}

				for _, newComment := range newReview.Comments {
					oldComment := findCommentByID(oldReview.Comments, newComment.ID)
					if oldComment == nil {
						changes = append(changes, fmt.Sprintf("New Review Comment: %s PR#%d Reviewer: %s CommentID: %d", newRepo.Name, newPR.Number, newReview.Login, newComment.ID))
						continue
					}

					if newComment.Body != oldComment.Body {
						changes = append(changes, fmt.Sprintf("Review Comment Body Change: %s PR#%d Reviewer: %s CommentID: %d", newRepo.Name, newPR.Number, newReview.Login, newComment.ID))
					}
				}
			}
		}
	}

	logrus.Info("Changes detected: ", strings.Join(changes, "\n"))
	return changes
}

func findRepoByName(repos []models.RepoState, name string) *models.RepoState {
	for _, repo := range repos {
		if repo.Name == name {
			return &repo
		}
	}
	return nil
}

func findPRByNumber(prs []models.PR, number int) *models.PR {
	for _, pr := range prs {
		if pr.Number == number {
			return &pr
		}
	}
	return nil
}

func findReviewByID(reviews []models.PRReview, id int64) *models.PRReview {
	for _, review := range reviews {
		if review.ID == id {
			return &review
		}
	}
	return nil
}

func findCommentByID(comments []models.PRReviewComment, id int64) *models.PRReviewComment {
	for _, comment := range comments {
		if comment.ID == id {
			return &comment
		}
	}
	return nil
}
