package service

import "errors"

var (
	ErrSiteContentNotFound = errors.New("site content resource not found")
	ErrInvalidSiteContent  = errors.New("invalid site content")
	ErrSiteContentLimit    = errors.New("site content limit exceeded")
	ErrSiteContentConflict = errors.New("site content conflict")
	ErrArticleNotPublished = errors.New("article is not published")
)

type ResourceOrder struct {
	ID        uint64
	SortOrder int
}

func validateResourceOrder(items []ResourceOrder, maximum int) error {
	if len(items) == 0 || len(items) > maximum {
		return ErrInvalidSiteContent
	}
	seen := make(map[uint64]struct{}, len(items))
	for _, item := range items {
		if item.ID == 0 {
			return ErrInvalidSiteContent
		}
		if _, exists := seen[item.ID]; exists {
			return ErrInvalidSiteContent
		}
		seen[item.ID] = struct{}{}
	}
	return nil
}
