package gorm

import (
	"kama_chat_server/internal/dao"
	"kama_chat_server/internal/dto/request"
	"kama_chat_server/internal/model"
	"kama_chat_server/pkg/zlog"
	"errors"
	"gorm.io/gorm"
	"time"
)

type groupInfoService struct {
}

var GroupInfoService = new(groupInfoService)

// SaveGroup 保存群聊，创建群聊也用这接口
func (g *groupInfoService) SaveGroup(groupReq request.GroupRequest) error {
	var group model.GroupInfo
	res := dao.GormDB.First(&group, "uuid = ?", groupReq.Uuid)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			// 创建群聊
			group.Uuid = groupReq.Uuid
			group.Name = groupReq.Name
			group.Notice = groupReq.Notice
			group.AddMode = groupReq.AddMode
			group.Avatar = groupReq.Avatar
			group.MemberCnt = 1
			group.Members = append(group.Members, groupReq.OwnerId)
			group.OwnerId = groupReq.OwnerId
			group.CreatedAt = time.Now()
			group.UpdatedAt = time.Now()
			if res := dao.GormDB.Create(&group); res.Error != nil {
				zlog.Error(res.Error.Error())
				return res.Error
			}
			return nil
		} else {
			zlog.Error(res.Error.Error())
			return res.Error
		}
	}
	// 更新群聊
	group.Uuid = groupReq.Uuid
	group.Name = groupReq.Name
	group.Notice = groupReq.Notice
	group.AddMode = groupReq.AddMode
	group.Avatar = groupReq.Avatar
	group.MemberCnt = 1
	group.Members = append(group.Members, groupReq.OwnerId)
	group.OwnerId = groupReq.OwnerId
	group.CreatedAt = time.Now()
	group.UpdatedAt = time.Now()
	return nil
}

// GetAllMembers 获取所有成员信息
func (g *groupInfoService) GetAllMembers(groupId string) ([]string, error) {
	var group model.GroupInfo
	res := dao.GormDB.Preload("Members").First(&group, "uuid = ?", groupId)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			zlog.Info("group not existed")
			return nil, nil
		} else {
			zlog.Error(res.Error.Error())
			return nil, res.Error
		}
	} else {
		return group.Members, nil
	}
}

// LoadMyGroup 获取我创建的群聊
func (g *groupInfoService) LoadMyGroup(ownerId string) []model.GroupInfo {
	var groupList []model.GroupInfo
	if res := dao.GormDB.Order("created_at DESC").Where("owner_id = ?", ownerId).Find(&groupList); res.Error != nil {
		zlog.Error(res.Error.Error())
		return nil
	}
	return groupList
}

// GetGroupInfo 获取聊天详情
func (g *groupInfoService) GetGroupInfo(userId string, groupId string) (model.GroupInfo, error) {
	var user model.UserInfo
	if res := dao.GormDB.First(&user, "uuid = ?", userId); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			zlog.Info(res.Error.Error())
			return model.GroupInfo{}, res.Error
		} else {
			zlog.Error(res.Error.Error())
			return model.GroupInfo{}, res.Error
		}
	}
	// 找到了，看是否被删除（禁用）
	if user.DeletedAt.Valid { // 已经删除了
		zlog.Info("user already be deleted, no access to this group")
		return model.GroupInfo{}, nil
	} else {
		var group model.GroupInfo
		if res := dao.GormDB.First(&group, "uuid = ?", groupId); res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				zlog.Info(res.Error.Error())
				return model.GroupInfo{}, res.Error
			} else {
				zlog.Error(res.Error.Error())
				return model.GroupInfo{}, res.Error
			}
		}
		return group, nil
	}
}

//func (g *groupInfoService) checkUserAndGroupValid(userId string, groupId string)

// GetGroupInfo4Chat 获取聊天会话群聊详情
func (g *groupInfoService) GetGroupInfo4Chat() error {

}

// LeaveGroup 退群
func (g *groupInfoService) LeaveGroup(userId string, groupId string) error {
	if userId == 
}
