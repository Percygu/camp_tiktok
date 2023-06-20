package service

import (
	"commentsvr/constant"
	"commentsvr/repository"
	"github.com/Percygu/litetiktok_proto/pb"
)

type CommentService struct{}

// PublishComment 发布评论
//func (c *CommentService) PublishComment(ctx context.Context, req *pb.PublishCommentReq) (*pb.PublishCommentRsp, error) {
//	var err error
//
//	tokenUids, _ := ctx.Get("UserId")
//
//	tokenUid := tokenUids.(int64)
//
//	video_id := ctx.Query("video_id")
//	comment_text := ctx.Query("comment_text")
//	actionType := ctx.Query("action_type")
//	comment_id := ctx.Query("comment_id")
//	commentId := int64(0)
//	if actionType == "2" {
//		commentId, err = strconv.ParseInt(req.CommentID, 10, 64)
//		if err != nil {
//			zap.L().Error("commentId error", zap.Error(err))
//			response.Fail(ctx, err.Error(), nil)
//			return
//		}
//	}
//	videoId, err := strconv.ParseInt(req.VideoId, 10, 64)
//	if err != nil {
//		logger.Errorf("videoId error : %s", err)
//		response.Fail(ctx, err.Error(), nil)
//		return
//	}
//
//	commentResponse, err := CommentAction(commentId, req.VideoId, tokenUid, comment_text, actionType)
//	if err != nil {
//		logger.Errorf("commentsvr error : %s", err)
//		response.Fail(ctx, err.Error(), nil)
//		return
//	}
//	response.Success(ctx, "success", commentResponse)
//}

// GetCommentList 获取评论列表
//func (c *CommentService) GetCommentList(ctx context.Context, req *proto.GetUserInfoRequest) (*pb.GetUserInfosp, error) {
//	var err error
//	video_id := ctx.Query("video_id")
//	/* token := ctx.Query("token")
//	_, err = util.VerifyToken(token)
//	if err != nil {
//		log.Errorf("token error : %s", err)
//		response.Fail(ctx, err.Error(), nil)
//		return
//	} */
//	videoId, err := strconv.ParseInt(video_id, 10, 64)
//	if err != nil {
//		logger.Errorf("videoId error : %s", err)
//		response.Fail(ctx, err.Error(), nil)
//		return
//	}
//	listResponse, err := service.CommentList(videoId)
//	if err != nil {
//		logger.Infof("list error : %s", err)
//		response.Fail(ctx, err.Error(), nil)
//		return
//	}
//	response.Success(ctx, "success", listResponse)
//}

func CommentAction(req *pb.PublishCommentReq) (*pb.PublishCommentRsp, error) {
	// 增加评论
	if req.ActionType == 1 {
		// commentInfo, err := repository.CommentAdd(userId, videoId, comment_text)
		_, err := repository.CommentAdd(req.UserId, req.VideoId, req.CommentText)
		if err != nil {
			return nil, err
		}
		return &pb.PublishCommentRsp{
			CommonRsp: &pb.CommonResponse{
				Code: constant.SuccessCode,
				Msg:  constant.SuccessMsg,
			},
		}, nil
	} else { // 删除评论
		err := repository.CommentDelete(req.VideoId, req.CommentId)
		if err != nil {
			return nil, err
		}
		return &pb.PublishCommentRsp{
			CommonRsp: &pb.CommonResponse{
				Code: constant.SuccessCode,
				Msg:  constant.SuccessMsg,
			},
		}, nil
	}
}
