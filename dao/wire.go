//go:build wireinject

package dao

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewBase64Captcha,
	NewTokenSessionStorage,
	NewSource,
	NewContact,
	NewContactGroup,
	NewGroupMember,
	NewUsers,
	NewGroup,
	NewGroupApply,
	NewTalkRecordGroup,
	NewTalkRecordFriend,
	NewTalkSession,
	NewTalkRecordGroupDel,
	NewEmoticon,
	NewGroupVote,
	NewFileUpload,
	NewArticleClass,
	NewArticle,
	NewArticleAnnex,
	NewDepartment,
	NewOrganize,
	NewPosition,
	NewRobot,
	NewSequence,
	NewAdmin,
	NewCaptchaStorage,
	//NewClientStorage,
	//NewContactRemark,
	//NewMessageStorage,
	//NewRelation,
	//NewRoomStorage,
	//NewSidStorage,
	//NewSmsStorage,
	//NewVote,
	//NewUnreadStorage,
	//NewGroupApplyStorage,
)
