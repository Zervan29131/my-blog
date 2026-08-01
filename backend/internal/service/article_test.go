package service

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"personal-blog/backend/internal/model"
)

func TestChangingPublishedArticleToDraftPreservesPublishedAt(t *testing.T) {
	database, mock := newServiceTestDatabase(t)
	publishedAt := time.Now().UTC().Add(-24 * time.Hour)
	createdAt := publishedAt.Add(-time.Hour)
	mock.ExpectQuery(`SELECT \* FROM "articles" WHERE "articles"\."id" = \$1 ORDER BY "articles"\."id" LIMIT \$2`).
		WithArgs(uint64(1), 1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "title", "slug", "summary", "content", "status", "published_at", "created_at", "updated_at",
		}).AddRow(
			uint64(1), "Published", "published-post", "Summary", "Content",
			model.ArticleStatusPublished, publishedAt, createdAt, publishedAt,
		))
	mock.ExpectQuery(`SELECT count\(\*\) FROM "articles" WHERE slug = \$1 AND id <> \$2`).
		WithArgs("published-post", uint64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "articles" SET .* WHERE id = \$[0-9]+`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	articleService := NewArticleService(database)
	article, err := articleService.Update(context.Background(), 1, ArticleInput{
		Title:   "Draft again",
		Slug:    "published-post",
		Summary: "Summary",
		Content: "Updated content",
		Status:  model.ArticleStatusDraft,
	})
	if err != nil {
		t.Fatalf("change article to draft: %v", err)
	}
	if article.PublishedAt == nil || !article.PublishedAt.Equal(publishedAt) {
		t.Fatalf("expected first publication time %s, got %v", publishedAt, article.PublishedAt)
	}
	if article.Status != model.ArticleStatusDraft {
		t.Fatalf("expected draft status, got %s", article.Status)
	}
}

func TestArticlePaginationCalculatesTotalPages(t *testing.T) {
	page := newArticlePage(nil, 2, 10, 21)
	if page.TotalPages != 3 {
		t.Fatalf("expected 3 total pages, got %d", page.TotalPages)
	}
}
