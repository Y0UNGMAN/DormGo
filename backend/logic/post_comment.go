package logic

import (
	"errors"

	"github.com/Y0UNGMAN/DormGo/backend/model"
	"gorm.io/gorm"
)

type CommentLogic struct{}

var Comment = new(CommentLogic)

type PublishCommentParams struct {
	PostID   uint
	UserID   uint
	Content  string
	ParentID uint
}

type GetPostCommentsResp struct {
	List  []*model.DgComment `json:"list"`
	Total int64              `json:"total"`
}

func (l *CommentLogic) PublishComment(p *PublishCommentParams) error {
	newComment := model.DgComment{
		PostID:      p.PostID,
		CommenterID: p.UserID,
		Content:     p.Content,
		ParentID:    p.ParentID,
		RootID:      0,
		ReplyToUID:  0,
	}

	return model.DB.Transaction(func(tx *gorm.DB) error {
		if p.ParentID > 0 {
			parentComment, err := model.GetCommentById(tx, p.ParentID)
			if err != nil {
				return errors.New("被回复的评论不存在")
			}
			newComment.ReplyToUID = parentComment.CommenterID
			if parentComment.RootID == 0 {
				newComment.RootID = parentComment.ID
			} else {
				newComment.RootID = parentComment.RootID
			}

			err = model.IncrementReplyCount(tx, newComment.RootID)
			if err != nil {
				return err
			}
		}

		err := model.CreateComment(tx, &newComment)
		if err != nil {
			return err
		}
		err = model.AddCommentCount(tx, newComment.PostID)
		return nil

	})

}

func (l *CommentLogic) GetPostComments(postID uint, page int, pageSize int) (*GetPostCommentsResp, error) {
	// 1. 计算分页偏移量
	offset := (page - 1) * pageSize

	//2. 调用 Model：获取顶级评论列表
	roots, total, err := model.GetRootCommentList(model.DB, postID, offset, pageSize)
	if err != nil {
		return nil, err
	}

	for i := range roots {
		if roots[i].ReplyCount > 0 {
			previews, secerr := model.GetPreviewReplies(model.DB, roots[i].ID, 3)
			if secerr == nil {
				roots[i].SubComments = previews
			}
		}
	}

	return &GetPostCommentsResp{
		List:  roots,
		Total: total,
	}, nil

}
