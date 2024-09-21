package store

import (
	"log"
	"strconv"
	"strings"
)

// Compare function to detect differences between two GHState structs
func CompareStates(oldState, newState GHState) []string {
	var changes []string

	for _, repo := range newState.RepoStates {
		if oldRepo := findRepoByName(oldState.RepoStates, repo.Name); oldRepo != nil {
			for _, pr := range repo.PRs {
				if oldPR := findPRByNumber(oldRepo.PRs, pr.Number); oldPR != nil {
					if pr.Body != oldPR.Body {
						changes = append(changes, "PR Body Change: "+repo.Name+" PR#"+strconv.Itoa(int(pr.Number))+"\n")
					}
					for _, review := range pr.Reviews {
						if oldReview := findReviewByID(oldPR.Reviews, review.ID); oldReview != nil {
							if review.Body != oldReview.Body {
								changes = append(changes, "Review Body Change: "+repo.Name+" PR#"+strconv.Itoa(int(pr.Number))+" Reviewer: "+review.Login+"\n")
							}
							for _, comment := range review.Comments {
								if oldComment := findCommentByID(oldReview.Comments, comment.ID); oldComment != nil {
									if comment.Body != oldComment.Body {
										changes = append(changes, "Review Comment Body Change: "+repo.Name+" PR#"+strconv.Itoa(int(pr.Number))+" Reviewer: "+review.Login+" CommentID: "+strconv.Itoa(int(comment.ID))+" --- " + comment.Body + "\n")
									}
								} else {
									changes = append(changes, "Review Comment Added: "+repo.Name+" PR#"+strconv.Itoa(int(pr.Number))+" Reviewer: "+review.Login+" CommentID: "+strconv.Itoa(int(comment.ID))+" ---  "+comment.Body+"\n")
								}
							}
						} else {
							changes = append(changes, "Review Added: "+repo.Name+" PR#"+strconv.Itoa(int(pr.Number))+" Reviewer: "+review.Login+"\n")
						}
					}
				} else {
					changes = append(changes, "PR State Change: "+repo.Name+" PR#"+strconv.Itoa(int(pr.Number))+"\n")
				}
			}
		} else {
			changes = append(changes, "Repo State Change: "+repo.Name+"\n")
		}
	}

	log.Println("my diff: ", strings.Join(changes, ""))

	return changes
}

func findRepoByName(prs []RepoState, name string) *RepoState {
	for _, repo := range prs {
		if repo.Name == name {
			return &repo
		}
	}
	return nil
}

func findPRByNumber(prs []PR, number int) *PR {
	for _, pr := range prs {
		if pr.Number == number {
			return &pr
		}
	}
	return nil
}

func findReviewByID(reviews []PRReview, id int64) *PRReview {
	for _, review := range reviews {
		if review.ID == id {
			return &review
		}
	}
	return nil
}

func findCommentByID(comments []PRReviewComment, id int64) *PRReviewComment {
	for _, comment := range comments {
		if comment.ID == id {
			return &comment
		}
	}
	return nil
}

