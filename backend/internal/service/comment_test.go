package service

import (
	"errors"
	"strings"
	"testing"

	"personal-blog/backend/internal/model"
)

func TestPrepareCommentInputValidatesLengthsAndEmail(t *testing.T) {
	tests := []struct {
		name  string
		input CommentInput
	}{
		{name: "nickname too short", input: CommentInput{Nickname: "a", Content: "valid content"}},
		{name: "nickname too long", input: CommentInput{Nickname: strings.Repeat("访", 51), Content: "valid content"}},
		{name: "content too short", input: CommentInput{Nickname: "valid", Content: "a"}},
		{name: "content too long", input: CommentInput{Nickname: "valid", Content: strings.Repeat("评", 1001)}},
		{name: "invalid email", input: CommentInput{Nickname: "valid", Email: "not-an-email", Content: "valid content"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, _, err := prepareCommentInput(test.input); !errors.Is(err, ErrInvalidComment) {
				t.Fatalf("expected invalid comment error, got %v", err)
			}
		})
	}
}

func TestCommentStatusesAreRestricted(t *testing.T) {
	for _, status := range []string{
		model.CommentStatusPending,
		model.CommentStatusApproved,
		model.CommentStatusRejected,
	} {
		if !IsValidCommentStatus(status) {
			t.Fatalf("expected %q to be valid", status)
		}
	}
	if IsValidCommentStatus("published") {
		t.Fatal("unexpected non-comment status to be valid")
	}
}
