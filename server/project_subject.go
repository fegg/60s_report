package calendar

import (
	"context"
	"github.com/go-kratos/kratos/pkg/log"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"med-common/app/service/med-calendar/model"
	"med-common/app/service/prescription-service/pkg/dbutil"
	"med-common/library/code"
	"time"
)

func (d *Dao) FetchOneProjectSubject(ctx context.Context, filter map[string]interface{}) (res *model.ProjectSubject, err error) {
	res = &model.ProjectSubject{}

	whereSql, params, err := dbutil.GetWhereConditionParams(filter)
	db := d.orm.Model(&model.ProjectSubject{}).Where(whereSql, params...)
	err = db.Order("id desc").First(res).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		log.Errorc(ctx, "[dao][FetchOneProjectSubject] 查询失败，err：%+v,", err)
	}
	return
}

func (d *Dao) GetProjectSubjectList(ctx context.Context, filter map[string]interface{}) (res []*model.ProjectSubject, err error) {
	res = make([]*model.ProjectSubject, 0)

	whereSql, params, err := dbutil.GetWhereConditionParams(filter)
	db := d.orm.Model(&model.ProjectSubject{}).Where(whereSql, params...)
	err = db.Order("id desc").Find(&res).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		log.Errorc(ctx, "[dao][GetProjectSubjectList] 查询失败，err：%+v,", err)
	}

	return
}

func (d *Dao) FetchProjectSubject(ctx context.Context, filter map[string]interface{}) (res map[int64]*model.ProjectSubject, err error) {
	res = make(map[int64]*model.ProjectSubject)

	var data []*model.ProjectSubject
	whereSql, params, err := dbutil.GetWhereConditionParams(filter)
	db := d.orm.Model(&model.ProjectSubject{}).Where(whereSql, params...)
	err = db.Order("id desc").Find(&data).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		log.Errorc(ctx, "[dao][FetchProjectSubject] 查询失败，err：%+v,", err)
	}
	for _, v := range data {
		res[v.Id] = v
	}
	return
}

func (d *Dao) FetchProjectSubjectGroupByProjectSubjectId(ctx context.Context, projectSubjectId int64, filter map[string]interface{}) (res map[int64]*model.ProjectSubjectGroup, err error) {
	res = make(map[int64]*model.ProjectSubjectGroup)

	filter["project_subject_id"] = projectSubjectId
	var data []*model.ProjectSubjectGroup
	whereSql, params, err := dbutil.GetWhereConditionParams(filter)
	db := d.orm.Model(&model.ProjectSubjectGroup{}).Where(whereSql, params...)
	err = db.Order("id desc").Find(&data).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		log.Errorc(ctx, "[dao][FetchProjectSubjectGroupByProjectSubjectId] 查询失败，err：%+v,", err)
	}
	for _, v := range data {
		res[v.Id] = v
	}
	return
}

func (d *Dao) FetchProjectSubjectCheckGroupByProjectSubjectId(ctx context.Context, projectSubjectId int64, filter map[string]interface{}) (res map[string]*model.ProjectSubjectCheckGroup, err error) {
	res = make(map[string]*model.ProjectSubjectCheckGroup)

	var checkGroupData []*model.ProjectSubjectCheckGroup

	filter["project_subject_id"] = projectSubjectId
	whereSql, params, err := dbutil.GetWhereConditionParams(filter)
	db := d.orm.Model(&model.ProjectSubjectCheckGroup{}).Where(whereSql, params...)
	err = db.Order("id desc").Find(&checkGroupData).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		log.Errorc(ctx, "[dao][FetchProjectSubjectCheckGroupByProjectSubjectId] 查询失败，err：%+v,", err)
		return
	}

	for _, v := range checkGroupData {
		res[v.UserId] = v
	}

	return
}

func (d *Dao) FetchProjectSubjectCheckGroup(ctx context.Context, filter map[string]interface{}) (res map[int64][]*model.ProjectSubjectCheckGroup, err error) {
	res = make(map[int64][]*model.ProjectSubjectCheckGroup)

	var checkGroupData []*model.ProjectSubjectCheckGroup

	whereSql, params, err := dbutil.GetWhereConditionParams(filter)
	db := d.orm.Model(&model.ProjectSubjectCheckGroup{}).Where(whereSql, params...)
	err = db.Order("id desc").Find(&checkGroupData).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		log.Errorc(ctx, "[dao][FetchProjectSubjectCheckGroupByProjectSubjectId] 查询失败，err：%+v,", err)
		return
	}

	for _, v := range checkGroupData {
		if _, ok := res[v.ProjectSubjectId]; !ok {
			res[v.ProjectSubjectId] = make([]*model.ProjectSubjectCheckGroup, 0)
		}
		res[v.ProjectSubjectId] = append(res[v.ProjectSubjectId], v)
	}

	return
}

func (d *Dao) GetProjectSubjectGroupByProjectSubjectId(ctx context.Context, projectSubjectId int64, filter map[string]interface{}) (res []*model.ProjectSubjectGroup, err error) {
	res = make([]*model.ProjectSubjectGroup, 0)

	filter["project_subject_id"] = projectSubjectId
	whereSql, params, err := dbutil.GetWhereConditionParams(filter)
	db := d.orm.Model(&model.ProjectSubjectGroup{}).Where(whereSql, params...)
	err = db.Order("id desc").Find(&res).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		log.Errorc(ctx, "[dao][GetProjectSubjectGroupByProjectSubjectId] 查询失败，err：%+v,", err)
		return
	}

	return
}

func (d *Dao) FetchProjectSubjectGroup(ctx context.Context, filter map[string]interface{}) (res map[int64][]*model.ProjectSubjectGroup, err error) {
	res = make(map[int64][]*model.ProjectSubjectGroup, 0)

	var groupData []*model.ProjectSubjectGroup
	whereSql, params, err := dbutil.GetWhereConditionParams(filter)
	db := d.orm.Model(&model.ProjectSubjectGroup{}).Where(whereSql, params...)
	err = db.Order("id desc").Find(&groupData).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		log.Errorc(ctx, "[dao][FetchProjectSubjectGroup] 查询失败，err：%+v,", err)
		return
	}

	for _, v := range groupData {
		if _, ok := res[v.ProjectSubjectId]; !ok {
			res[v.ProjectSubjectId] = make([]*model.ProjectSubjectGroup, 0)
		}
		res[v.ProjectSubjectId] = append(res[v.ProjectSubjectId], v)
	}

	return
}

func (d *Dao) GetProjectSubjectGroupScore(ctx context.Context, filter map[string]interface{}) (res []*model.ProjectSubjectGroupScore, err error) {
	res = make([]*model.ProjectSubjectGroupScore, 0)

	whereSql, params, err := dbutil.GetWhereConditionParams(filter)
	db := d.orm.Model(&model.ProjectSubjectGroupScore{}).Where(whereSql, params...)
	err = db.Order("id desc").Find(&res).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		log.Errorc(ctx, "[dao][GetProjectSubjectGroupScore] 查询失败，err：%+v,", err)
		return
	}

	return
}

func (d *Dao) FetchProjectSubjectGroupScore(ctx context.Context, filter map[string]interface{}) (res map[int64][]*model.ProjectSubjectGroupScore, err error) {
	res = make(map[int64][]*model.ProjectSubjectGroupScore, 0)

	var scoreData []*model.ProjectSubjectGroupScore
	whereSql, params, err := dbutil.GetWhereConditionParams(filter)
	db := d.orm.Model(&model.ProjectSubjectGroupScore{}).Where(whereSql, params...)
	err = db.Order("id desc").Find(&scoreData).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		log.Errorc(ctx, "[dao][FetchProjectSubjectGroupScore] 查询失败，err：%+v,", err)
		return
	}

	for _, v := range scoreData {
		if _, ok := res[v.ProjectSubjectId]; !ok {
			res[v.ProjectSubjectId] = make([]*model.ProjectSubjectGroupScore, 0)
		}
		res[v.ProjectSubjectId] = append(res[v.ProjectSubjectId], v)
	}

	return
}

func (d *Dao) PostProjectSubject(ctx context.Context, projectSubject *model.ProjectSubject, project *model.Project, checkGroups []*model.ProjectSubjectCheckGroup, parentProject *model.Project) (res int64, err error) {
	new_db := d.orm.Begin()
	if err := new_db.Error; err != nil {
		log.Errorc(ctx, "[dao][PostProjectSubject] 启动事务失败，err：%+v,", err)
		return 0, errors.Wrap(err, "启动事务失败")
	}
	defer func() {
		if err != nil {
			new_db.Rollback()
			log.Errorc(ctx, "[dao][PostProjectSubject] 提交项目开题数据失败，err：%+v,", err)
			return
		}
	}()
	// 创建
	if projectSubject.Id == 0 {
		// 创建项目开题数据
		if err = new_db.Model(&model.ProjectSubject{}).Create(&projectSubject).Error; err != nil {
			return
		}
		// 创建项目数据
		var count int64 = 0
		d.orm.Model(&model.Project{}).Where("state = ?", 1).Where("project_name = ?", project.ProjectName).Count(&count)
		if count > 0 {
			err = code.MedErrProjectNameExists
			return
		}
		// 判断父项目coin是否足够
		if project.ParentId > 0 && parentProject.Id > 0 {
			efCount := new_db.Model(&model.Project{}).Where("id = ?", project.ParentId).
				Where("coin - ? >= 0", project.Coin).
				Updates(map[string]interface{}{
					"coin": gorm.Expr("coin - ?", project.Coin),
				}).RowsAffected
			if efCount == 0 {
				err = code.MedErrMedCoinNotEnough
				return
			}
		} else { // 扣除总办
			efCount := new_db.Model(&model.OfficeWater{}).Where("year = ?", time.Now().Year()).
				Where("coin - ? >= 0", project.Coin).
				Updates(map[string]interface{}{
					"coin": gorm.Expr("coin - ?", project.Coin),
				}).RowsAffected
			if efCount == 0 {
				err = code.MedErrMedAgentCoinNotEnough
				return
			}
		}
		project.ProjectSubjectId = projectSubject.Id
		if err = new_db.Model(&model.Project{}).Create(&project).Error; err != nil {
			return
		}
	} else {
		// 修改
		var dbProjectSubject *model.ProjectSubject
		if err = new_db.Model(&model.ProjectSubject{}).Where("id = ?", projectSubject.Id).First(&dbProjectSubject).Error; err != nil {
			if err = dbutil.IgnoreNoRecordErr(err); err != nil {
				log.Errorc(ctx, "[dao][PostProjectSubject] 查询开题详情失败 params(%d)，err：%+v,", projectSubject.Id, err)
				return
			} else {
				return 0, code.MedErrProjectSubjectNotExists
			}
		}
		if dbProjectSubject.Status == 3 {
			return 0, code.MedErrProjectSubjectHasFinish
		}
		if err = new_db.Model(&model.ProjectSubject{}).Where("id = ?", projectSubject.Id).Updates(projectSubject).Error; err != nil {
			return 0, code.MedErrDBFailed
		}
		// 删掉之前的项目课题考核人
		if err = new_db.Model(&model.ProjectSubjectCheckGroup{}).Where("project_subject_id = ?", projectSubject.Id).Updates(map[string]interface{}{"state": -1}).Error; err != nil {
			return 0, code.MedErrDBFailed
		}
		if len(checkGroups) == 0 && projectSubject.Status == 2 { // 如果开启打分必须设置考核人
			return 0, code.MedErrHasNotSetChecker
		}
		if len(checkGroups) > 0 {
			// 新增新的项目课题考核人
			for _, v := range checkGroups {
				v.ProjectSubjectId = projectSubject.Id
			}
			if err = new_db.Model(&model.ProjectSubjectCheckGroup{}).Where("project_subject_id = ?", projectSubject.Id).Create(&checkGroups).Error; err != nil {
				return 0, code.MedErrDBFailed
			}
		}
	}

	if err = new_db.Commit().Error; err != nil {
		new_db.Rollback()
		log.Errorc(ctx, "[dao][PostProjectSubject] 提交事务失败，err：%+v,", err)
		return 0, err
	}
	res = projectSubject.Id

	return
}

func (d *Dao) BidProjectSubject(ctx context.Context, userId string, projectSubjectId int64, groupName string, remark string) (err error) {
	if projectSubjectId == 0 {
		return
	}
	var dbProjectSubject *model.ProjectSubject
	if err = d.orm.Model(&model.ProjectSubject{}).Where("id = ?", projectSubjectId).First(&dbProjectSubject).Error; err != nil {
		if err = dbutil.IgnoreNoRecordErr(err); err != nil {
			return code.MedErrDBFailed
		} else {
			return code.MedErrProjectSubjectNotExists
		}
	}
	switch dbProjectSubject.Status {
	case 2:
		return code.MedErrCanNotBidProjectSubject
	case 3:
		return code.MedErrProjectSubjectHasFinish
	}
	var dbProjectSubjectGroup *model.ProjectSubjectGroup
	if err = d.orm.Model(&model.ProjectSubjectGroup{}).Where("project_subject_id = ?", projectSubjectId).Where("captain_user_id = ?", userId).Where("state = ?", 1).
		First(&dbProjectSubjectGroup).Error; err != nil {
		if err = dbutil.IgnoreNoRecordErr(err); err != nil {
			return
		}
	}
	if dbProjectSubjectGroup.Id > 0 {
		// 已经投标了
		return code.MedErrHasBidProjectSubject
	}

	projectSubjectGroup := &model.ProjectSubjectGroup{
		ProjectSubjectId: projectSubjectId,
		GroupName:        groupName,
		CaptainUserId:    userId,
		GroupUserIds:     "[]",
		State:            1,
		CreatedAt:        time.Now(),
		Remark:           remark,
	}

	if err = d.orm.Model(&model.ProjectSubjectGroup{}).Create(&projectSubjectGroup).Error; err != nil {
		log.Errorc(ctx, "[dao][BidProjectSubject] 项目竞标数据库错误 params(%+v)，err：%+v,", *projectSubjectGroup, err)
		return
	}
	return
}



func (d *Dao) UpdateProjectSubject(ctx context.Context, userId string, projectSubjectId int64, groupName string, remark string) (err error) {
	if projectSubjectId == 0 {
		return
	}
	var dbProjectSubject *model.ProjectSubject
	if err = d.orm.Model(&model.ProjectSubject{}).Where("id = ?", projectSubjectId).First(&dbProjectSubject).Error; err != nil {
		if err = dbutil.IgnoreNoRecordErr(err); err != nil {
			return code.MedErrDBFailed
		} else {
			return code.MedErrProjectSubjectNotExists
		}
	}

	var dbProjectSubjectGroup *model.ProjectSubjectGroup
	if err = d.orm.Model(&model.ProjectSubjectGroup{}).Where("project_subject_id = ?", projectSubjectId).Where("captain_user_id = ?", userId).Where("state = ?", 1).
		First(&dbProjectSubjectGroup).Error; err != nil {
		if err = dbutil.IgnoreNoRecordErr(err); err == nil {
			return code.MedErrGroupNotBidProjectSubject
		}
	}

	projectSubjectGroupMap := make(map[string]interface{})
	projectSubjectGroupMap["remark"] = remark

	if err = d.orm.Model(&model.ProjectSubjectGroup{}).Where("id = ?", dbProjectSubjectGroup.Id).
		Updates(projectSubjectGroupMap).Error; err != nil {
		log.Errorc(ctx, "[dao][UpdateProjectSubject] 更新竞标组信息失败 params(%s, %s)，err：%+v,", groupName, remark, err)
		return
	}
	return
}

func (d *Dao) ScoreProjectSubjectGroup(ctx context.Context, userId string, projectSubjectId int64, checkGroups []*model.ProjectSubjectGroupScore) (err error) {
	new_db := d.orm.Begin()
	if err := new_db.Error; err != nil {
		log.Errorc(ctx, "[dao][ScoreProjectSubjectGroup] 启动事务失败，err：%+v,", err)
		return errors.Wrap(err, "启动事务失败")
	}
	defer func() {
		if err != nil {
			new_db.Rollback()
			log.Errorc(ctx, "[dao][ScoreProjectSubjectGroup] 项目课题打分失败，err：%+v,", err)
			return
		}
	}()
	var dbProjectSubject *model.ProjectSubject
	if err = d.orm.Model(&model.ProjectSubject{}).Where("id = ?", projectSubjectId).First(&dbProjectSubject).Error; err != nil {
		if err = dbutil.IgnoreNoRecordErr(err); err != nil {
			return code.MedErrDBFailed
		} else {
			return code.MedErrProjectSubjectNotExists
		}
	}
	switch dbProjectSubject.Status {
	case 1:
		return code.MedErrCanNotScoreProjectSubject
	case 3:
		return code.MedErrProjectSubjectHasFinish
	}

	if len(checkGroups) == 0 {
		return
	}
	// 删除老的
	if err = new_db.Model(&model.ProjectSubjectGroupScore{}).Where("project_subject_id = ?", projectSubjectId).Where("check_user_id = ?", userId).Where("state = ?", 1).
		Updates(map[string]interface{}{"state": -1}).Error; err != nil {
		return code.MedErrDBFailed
	}
	if err = new_db.Model(&model.ProjectSubjectGroupScore{}).Create(&checkGroups).Error; err != nil {
		return code.MedErrDBFailed
	}
	// 如果所有人都已经打分
	if d.checkAllHasScore(new_db, projectSubjectId) {
		// 标记 项目课题标记为已完成
		if err = new_db.Model(&model.ProjectSubject{}).Where("id = ?", projectSubjectId).Where("status = ?", 2).Updates(map[string]interface{}{"status": 3}).Error; err != nil {
			return code.MedErrDBFailed
		}
	} else {
		// 标记 项目课题打分中
		//if err = new_db.Model(&model.ProjectSubject{}).Where("id = ?", projectSubjectId).Where("status = ?", 2).Updates(map[string]interface{}{"state": 2}).Error; err != nil {
		//	return code.MedErrDBFailed
		//}
	}

	if err = new_db.Commit().Error; err != nil {
		new_db.Rollback()
		log.Errorc(ctx, "[dao][ScoreProjectSubjectGroup] 提交事务失败，err：%+v,", err)
		return errors.Wrap(err, "提交事务失败")
	}
	return
}

// 判断项目是否全部完成打分
func (d *Dao) checkAllHasScore(new_db *gorm.DB, projectSubjectId int64) bool {
	// 已经打分的考核人id
	hasScoreUserIds := make([]string, 0)
	if err := new_db.Model(&model.ProjectSubjectGroupScore{}).Where("project_subject_id = ?", projectSubjectId).Where("state = ?", 1).Group("check_user_id").Pluck("check_user_id", &hasScoreUserIds).Error; err != nil {
		return false
	}
	// 项目所有考核人
	checkGroupUserIds := make([]string, 0)
	if err := new_db.Model(&model.ProjectSubjectCheckGroup{}).Where("project_subject_id = ?", projectSubjectId).Where("state = ?", 1).Pluck("user_id", &checkGroupUserIds).Error; err != nil {
		return false
	}

	if len(checkGroupUserIds) > 0 && len(hasScoreUserIds) == len(checkGroupUserIds) {
		return true
	}
	return false
}
